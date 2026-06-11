// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command availability is a synthetic monitor for the dws command surface.
//
// It hits a real discovery endpoint (production or pre), enumerates every
// server and tool the way the CLI does at runtime, builds the canonical
// catalog, and reports whether each command is actually resolvable. It is
// meant to run unattended on a schedule (see
// .github/workflows/command-availability.yml) and exits non-zero the moment
// any command is "broken" so CI can alert.
//
// "Broken" here covers the failure modes the project has actually hit:
//   - a server whose runtime list-tools call fails (stale/wrong discovery
//     config, gateway down) -> runtime failure
//   - two tools collapsing onto the same CLI command name so one shadows the
//     other (the list-direct class of bug) -> name collision
//   - a tool that resolves to no usable command name at all -> unresolvable
//
// Usage:
//
//	availability [--endpoint URL] [--json] [--timeout 90s]
//
// Endpoint resolution order: --endpoint flag, then DWS_DISCOVERY_URL env,
// then the production default. Set DWS_PAT to inject an Authorization bearer
// header for endpoints that gate list-tools behind auth.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	authpkg "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/cache"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/discovery"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/ir"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/market"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/transport"
)

// defaultDiscoveryURL is the production discovery endpoint. Inlined rather
// than imported from internal/config so the monitor compiles across CLI
// versions whose config constants may have drifted.
const defaultDiscoveryURL = "https://mcp.dingtalk.com"

// Exit codes are part of the contract with CI: they let the workflow route a
// failure to the right alert (commands broken vs. credentials expired) instead
// of lumping every non-zero into one scary "broken" message.
const (
	exitHealthy         = 0
	exitBroken          = 1 // real command issues: runtime failures, collisions, unresolvable
	exitInternalError   = 2 // the monitor itself could not run
	exitAuthUnavailable = 3 // --require-auth set but no usable dws credentials
)

// Report is the machine-readable monitor result emitted with --json.
type Report struct {
	Endpoint        string           `json:"endpoint"`
	Status          string           `json:"status"` // healthy | broken | auth_unavailable | server_list_only
	ServersFound    int              `json:"servers_found"`
	ServersRuntime  int              `json:"servers_runtime_ok"`
	CommandsChecked int              `json:"commands_checked"`
	RuntimeChecked  bool             `json:"runtime_checked"`
	AuthError       string           `json:"auth_error,omitempty"`
	Healthy         bool             `json:"healthy"`
	RuntimeFailures []failureEntry   `json:"runtime_failures"`
	Collisions      []collisionEntry `json:"command_name_collisions"`
	Unresolvable    []string         `json:"unresolvable_tools"`
}

type failureEntry struct {
	ServerKey string `json:"server_key"`
	Error     string `json:"error"`
}

type collisionEntry struct {
	Product  string   `json:"product"`
	Command  string   `json:"command"`
	RPCNames []string `json:"rpc_names"`
}

// knownMergedCommands lists command paths where several RPC tools map onto
// ONE CLI command on purpose: the discovery overlay publishes the union of
// their flags and the surviving leaf routes by flag presence, so the overlap
// is a design decision rather than the list-direct shadowing bug. Keyed by
// the product ID with its instance-hash suffix stripped (chat-f80f4305 ->
// chat), because the suffix is not stable across environments.
//
// Verified entries only. "message send" merges send_message_as_user (group,
// --group) with send_direct_message_as_user (direct, --user /
// --open-dingtalk-id); both paths confirmed working against production.
var knownMergedCommands = map[string]map[string]bool{
	"chat": {"message send": true},
}

// productKey strips the trailing instance-hash suffix from a product ID
// ("chat-f80f4305" -> "chat") so allowlist entries stay environment-agnostic.
// Only a suffix that looks like a hex hash is removed.
func productKey(id string) string {
	if i := strings.LastIndex(id, "-"); i > 0 {
		suffix := id[i+1:]
		if len(suffix) >= 6 && strings.IndexFunc(suffix, func(r rune) bool {
			return !strings.ContainsRune("0123456789abcdef", r)
		}) == -1 {
			return id[:i]
		}
	}
	return id
}

func main() {
	endpointFlag := flag.String("endpoint", "", "discovery endpoint (default: $DWS_DISCOVERY_URL or production)")
	jsonOut := flag.Bool("json", false, "emit the full report as JSON")
	timeout := flag.Duration("timeout", 90*time.Second, "overall timeout for the sweep")
	requireAuth := flag.Bool("require-auth", false, "exit 3 if no dws credentials are available instead of degrading to a server-list-only sweep")
	selftest := flag.String("selftest", "", "emit a synthetic report to exercise the alert path without touching the real surface: 'broken' or 'auth'")
	flag.Parse()

	endpoint := resolveEndpoint(*endpointFlag)

	// Self-test path: synthesize a result to verify the downstream alert
	// pipeline end to end. Findings are clearly marked "selftest" so a real
	// reader can't mistake them for an actual outage.
	switch strings.ToLower(strings.TrimSpace(*selftest)) {
	case "broken":
		emit(Report{
			Endpoint: endpoint, Status: "broken", RuntimeChecked: true, Healthy: false,
			Collisions: []collisionEntry{{Product: "selftest", Command: "synthetic alert", RPCNames: []string{"selftest_a", "selftest_b"}}},
		}, *jsonOut)
		os.Exit(exitBroken)
	case "auth":
		emit(Report{Endpoint: endpoint, Status: "auth_unavailable", Healthy: false}, *jsonOut)
		os.Exit(exitAuthUnavailable)
	}

	httpClient := &http.Client{Timeout: 60 * time.Second}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Runtime list-tools at mcp-gw is auth-gated. Use the transport client's
	// own WithAuth so we send the SAME headers the CLI does — both
	// "Authorization: Bearer" AND "x-user-access-token" — gated by its trusted
	// -domain check (defaults to *.dingtalk.com). A bare Bearer header alone
	// gets a 400. Server-list discovery (market) stays unauthenticated.
	transportClient := transport.NewClient(httpClient)
	hasCreds := false
	token, authErr := resolveToken(ctx)
	if token != "" {
		transportClient = transportClient.WithAuth(token, nil)
		hasCreds = true
	} else if *requireAuth {
		// Credentials were expected (CI) but couldn't be resolved — most often
		// an expired/rotated auth bundle. Signal this distinctly (exit 3) so it
		// alerts as "refresh credentials", not "commands are broken". Surface
		// the underlying reason so the operator isn't left guessing.
		if authErr != "" {
			fmt.Fprintf(os.Stderr, "availability: auth unavailable: %s\n", authErr)
		}
		emit(Report{Endpoint: endpoint, Status: "auth_unavailable", AuthError: authErr, Healthy: false}, *jsonOut)
		os.Exit(exitAuthUnavailable)
	} else {
		fmt.Fprintln(os.Stderr, "availability: no dws credentials (run 'dws auth login' "+
			"or set DWS_PAT); checking the server-list layer only")
	}

	tmp, err := os.MkdirTemp("", "dws-availability-cache-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "availability: cannot create cache dir: %v\n", err)
		os.Exit(exitInternalError)
	}
	defer os.RemoveAll(tmp)

	svc := discovery.NewService(
		market.NewClient(endpoint, httpClient),
		transportClient,
		cache.NewStore(tmp),
	)

	servers, err := svc.DiscoverServers(ctx)
	if err != nil {
		// A total discovery failure is the loudest possible "broken" signal:
		// nobody can run any command. Report and fail hard.
		report := Report{Endpoint: endpoint, Status: "broken", Healthy: false,
			RuntimeFailures: []failureEntry{{ServerKey: "<discovery>", Error: err.Error()}}}
		emit(report, *jsonOut)
		os.Exit(exitBroken)
	}

	// Without credentials, the per-server runtime probe would 400 on every
	// server and report false breakage. Don't cry wolf: validate only that
	// the server list is reachable, and say plainly that runtime was skipped.
	if !hasCreds {
		report := Report{
			Endpoint:       endpoint,
			Status:         "server_list_only",
			ServersFound:   len(servers),
			RuntimeChecked: false,
			Healthy:        len(servers) > 0,
		}
		emit(report, *jsonOut)
		if !report.Healthy {
			os.Exit(exitBroken)
		}
		return
	}

	runtimeServers, failures := svc.DiscoverAllRuntime(ctx, servers)
	catalog := ir.BuildCatalog(runtimeServers)

	report := Report{
		Endpoint:       endpoint,
		Status:         "healthy",
		ServersFound:   len(servers),
		ServersRuntime: len(runtimeServers),
		RuntimeChecked: true,
		Healthy:        true,
	}

	for _, f := range failures {
		report.RuntimeFailures = append(report.RuntimeFailures, failureEntry{
			ServerKey: f.ServerKey,
			Error:     errString(f.Err),
		})
	}

	for _, product := range catalog.Products {
		byCommand := map[string][]string{}
		merged := knownMergedCommands[productKey(product.ID)]
		for _, tool := range product.Tools {
			if tool.Hidden {
				continue
			}
			report.CommandsChecked++
			cmd := effectiveCommand(tool)
			if cmd == "" {
				report.Unresolvable = append(report.Unresolvable,
					fmt.Sprintf("%s.%s", product.ID, tool.RPCName))
				continue
			}
			byCommand[cmd] = append(byCommand[cmd], tool.RPCName)
		}
		for cmd, rpcs := range byCommand {
			if len(rpcs) > 1 && !merged[cmd] {
				sort.Strings(rpcs)
				report.Collisions = append(report.Collisions, collisionEntry{
					Product: product.ID, Command: cmd, RPCNames: rpcs,
				})
			}
		}
	}

	sort.Slice(report.RuntimeFailures, func(i, j int) bool {
		return report.RuntimeFailures[i].ServerKey < report.RuntimeFailures[j].ServerKey
	})
	sort.Slice(report.Collisions, func(i, j int) bool {
		return report.Collisions[i].Product+report.Collisions[i].Command <
			report.Collisions[j].Product+report.Collisions[j].Command
	})
	sort.Strings(report.Unresolvable)

	report.Healthy = len(report.RuntimeFailures) == 0 &&
		len(report.Collisions) == 0 &&
		len(report.Unresolvable) == 0
	if !report.Healthy {
		report.Status = "broken"
	}

	emit(report, *jsonOut)
	if !report.Healthy {
		os.Exit(exitBroken)
	}
}

func resolveEndpoint(flagVal string) string {
	if v := strings.TrimSpace(flagVal); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv("DWS_DISCOVERY_URL")); v != "" {
		return v
	}
	return defaultDiscoveryURL
}

// resolveToken returns a bearer token for the runtime probe. An explicit
// DWS_PAT env wins (handy for headless/CI overrides); otherwise it reads the
// logged-in dws OAuth credentials from the config dir and auto-refreshes via
// the refresh token. Empty return means "not logged in" — the caller degrades
// to a server-list-only sweep rather than failing.
func resolveToken(ctx context.Context) (token string, errReason string) {
	if t := strings.TrimSpace(os.Getenv("DWS_PAT")); t != "" {
		return t, ""
	}
	configDir := strings.TrimSpace(os.Getenv("DWS_CONFIG_DIR"))
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", "cannot resolve home dir: " + err.Error()
		}
		configDir = filepath.Join(home, ".dws")
	}
	silent := slog.New(slog.NewTextHandler(io.Discard, nil))
	provider := authpkg.NewOAuthProvider(configDir, silent)
	t, err := provider.GetAccessToken(ctx)
	if err != nil {
		return "", err.Error()
	}
	return strings.TrimSpace(t), ""
}

// effectiveCommand mirrors how the CLI derives a command's full path: the
// noun-group plus the verb (CLIName), falling back to the RPC name when no
// CLIName override exists. e.g. create_base -> Group "base" + CLIName
// "create" = "base create", which is distinct from create_chart's "chart
// create". Using CLIName alone would wrongly collapse every "create" verb
// into one command and report false collisions. A real collision is two
// tools resolving to the SAME group+verb path (the list-direct class of bug).
func effectiveCommand(tool ir.ToolDescriptor) string {
	cli := strings.TrimSpace(tool.CLIName)
	if cli == "" {
		cli = strings.TrimSpace(tool.RPCName)
	}
	if group := strings.TrimSpace(tool.Group); group != "" {
		return group + " " + cli
	}
	return cli
}

func errString(err error) string {
	if err == nil {
		return "unknown error"
	}
	return err.Error()
}

func emit(r Report, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(r)
		return
	}

	status := "HEALTHY"
	if !r.Healthy {
		status = "BROKEN"
	}
	fmt.Printf("dws command availability: %s\n", status)
	fmt.Printf("  endpoint:         %s\n", r.Endpoint)
	fmt.Printf("  servers found:    %d\n", r.ServersFound)
	if !r.RuntimeChecked {
		fmt.Printf("  runtime:          skipped (no credentials — server-list only)\n")
		return
	}
	fmt.Printf("  servers runtime:  %d\n", r.ServersRuntime)
	fmt.Printf("  commands checked: %d\n", r.CommandsChecked)

	if len(r.RuntimeFailures) > 0 {
		fmt.Printf("\n  runtime failures (%d):\n", len(r.RuntimeFailures))
		for _, f := range r.RuntimeFailures {
			fmt.Printf("    - %s: %s\n", f.ServerKey, f.Error)
		}
	}
	if len(r.Collisions) > 0 {
		fmt.Printf("\n  command-name collisions (%d):\n", len(r.Collisions))
		for _, c := range r.Collisions {
			fmt.Printf("    - %s %q <- %s\n", c.Product, c.Command, strings.Join(c.RPCNames, ", "))
		}
	}
	if len(r.Unresolvable) > 0 {
		fmt.Printf("\n  unresolvable tools (%d):\n", len(r.Unresolvable))
		for _, u := range r.Unresolvable {
			fmt.Printf("    - %s\n", u)
		}
	}
}
