// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/testseam"
)

func isolatedFilesystemFixture() map[string]string {
	return map[string]string{":root": "deny", ":minimal": "read", ":workspace_roots": "write", ":tmpdir": "deny", ":slash_tmp": "deny"}
}

func isolatedSavedFixture(cwd string) codexSavedSettings {
	entries := []any{}
	for path, access := range isolatedFilesystemFixture() {
		entries = append(entries, map[string]any{"path": map[string]any{"type": "special", "value": map[string]string{"kind": strings.TrimPrefix(path, ":")}}, "access": access})
	}
	profile, _ := json.Marshal(map[string]any{"type": "managed", "network": "restricted", "file_system": map[string]any{"type": "restricted", "entries": entries}})
	return codexSavedSettings{Cwd: cwd, Roots: []string{cwd}, Approval: json.RawMessage(`"on-request"`), Reviewer: "user", Profile: profile}
}

func TestCrossPlatformCoverageCodexNativeIsolationForward(t *testing.T) {
	for _, scenario := range []string{"start", "resume", "unknown-profile", "unsupported", "network-enabled", "missing-temp-deny", "extra-write", "mcp-config", "plugin-config", "connector-config", "wrong-start-profile", "wider-saved", "wrong-cwd", "resume-error", "wrong-resumed-thread"} {
		t.Run(scenario, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("CODEX_HOME", home)
			root := t.TempDir()
			cwd, err := codexIsolatedConversationDir(root, "conversation")
			if err != nil {
				t.Fatal(err)
			}
			fs := isolatedFilesystemFixture()
			profileConfig := map[string]any{"filesystem": fs, "network": map[string]bool{"enabled": false}}
			profiles := map[string]any{"strict": profileConfig}
			if scenario == "unknown-profile" {
				profiles = map[string]any{}
			}
			if scenario == "network-enabled" {
				profileConfig["network"] = map[string]bool{"enabled": true}
			}
			if scenario == "missing-temp-deny" {
				delete(fs, ":tmpdir")
			}
			if scenario == "extra-write" {
				fs[root] = "write"
			}
			configuration := map[string]any{"permissions": profiles}
			if scenario == "mcp-config" {
				configuration["mcp_servers"] = map[string]any{"external": map[string]bool{"enabled": true}}
			}
			if scenario == "plugin-config" {
				configuration["plugins"] = map[string]any{"external": map[string]bool{"enabled": true}}
			}
			if scenario == "connector-config" {
				configuration["apps"] = map[string]any{"external": map[string]bool{"enabled": true}}
			}
			config, _ := json.Marshal(map[string]any{"config": configuration})
			writer := &bufferWriteCloser{}
			messages := []codexRPCMessage{{ID: codexIntPtr(1), Result: json.RawMessage(`{}`)}, {ID: codexIntPtr(2), Result: config}}
			if scenario == "unsupported" {
				messages[1].Error = &codexRPCError{Code: -32601, Message: "unsupported"}
			}
			f := &codexAppServerForwarder{bin: "fixture", workDir: root, permissionProfile: "strict", developerInstructions: "employee", sessions: newCodexThreadSessions("")}
			result := map[string]any{"thread": map[string]string{"id": "thread"}, "cwd": cwd, "runtimeWorkspaceRoots": []string{cwd}, "approvalPolicy": "on-request", "approvalsReviewer": "user", "activePermissionProfile": map[string]string{"id": "strict"}}
			if scenario == "wrong-start-profile" {
				result["activePermissionProfile"] = nil
			}
			resume := scenario != "start" && scenario != "wrong-start-profile" && scenario != "unknown-profile" && scenario != "unsupported"
			if resume {
				f.sessions.setThreadID("conversation", "thread")
				saved := isolatedSavedFixture(cwd)
				if scenario == "wider-saved" {
					saved.Profile = json.RawMessage(`{"type":"disabled"}`)
				}
				if scenario == "wrong-cwd" {
					saved.Cwd = root
				}
				meta, _ := json.Marshal(map[string]any{"type": "session_meta", "payload": map[string]string{"id": "thread"}})
				turn, _ := json.Marshal(map[string]any{"type": "turn_context", "payload": saved})
				path := filepath.Join(home, "rollout.jsonl")
				if err := os.WriteFile(path, append(append(meta, '\n'), append(turn, '\n')...), 0600); err != nil {
					t.Fatal(err)
				}
				read, _ := json.Marshal(map[string]any{"thread": map[string]string{"id": "thread", "path": path}})
				messages = append(messages, codexRPCMessage{ID: codexIntPtr(3), Result: read})
			}
			if scenario == "wrong-resumed-thread" {
				result["thread"] = map[string]string{"id": "other"}
			}
			raw, _ := json.Marshal(result)
			id := 3
			if resume {
				id = 4
			}
			response := codexRPCMessage{ID: codexIntPtr(id), Result: raw}
			if scenario == "resume-error" {
				response.Error = &codexRPCError{Code: -32600, Message: "already has an active writer"}
			}
			messages = append(messages, response, codexRPCMessage{Method: "turn/completed", Params: json.RawMessage(`{"threadId":"thread","turn":{"status":"completed","items":[{"type":"agentMessage","text":"done"}]}}`)})
			client := unitCodexClient(writer, messages...)
			testseam.Swap(t, &codexNewAppServerClient, func(_ context.Context, _ string, env []string, actual string) (*codexAppServerClient, error) {
				if actual != cwd {
					t.Fatalf("wrong process cwd %s", actual)
				}
				if !strings.Contains(strings.Join(env, "\n"), "CODEX_HOME=") {
					t.Fatal("dedicated home lost")
				}
				return client, nil
			})
			_, err = f.forward(context.Background(), "conversation", "hello")
			wantOK := scenario == "start" || scenario == "resume"
			if (err == nil) != wantOK {
				t.Fatalf("err=%v", err)
			}
			turns := 0
			for _, request := range decodeCodexRequests(t, writer.String()) {
				if request.Method == "thread/start" || request.Method == "thread/resume" || request.Method == "turn/start" {
					if request.Params["permissions"] != "strict" || request.Params["sandbox"] != nil || request.Params["sandboxPolicy"] != nil || request.Params["config"] != nil || request.Params["cwd"] != cwd || !reflect.DeepEqual(request.Params["runtimeWorkspaceRoots"], []any{cwd}) {
						t.Fatalf("unsafe request: %+v", request)
					}
				}
				if request.Method == "turn/start" {
					turns++
				}
			}
			if !wantOK && turns != 0 {
				t.Fatal("model called after isolation failure")
			}
			if resume && !wantOK && f.sessions.threadID("conversation") != "thread" {
				t.Fatal("failed resume discarded history")
			}
		})
	}
}

func TestCrossPlatformCoverageCodexIsolationOptions(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	for _, scenario := range []string{"valid", "missing-workdir", "memory-off", "legacy-bypass", "builtin-profile"} {
		opts := connectAgentOptions{WorkDir: t.TempDir(), Memory: true, CodexPermissionProfile: "strict"}
		switch scenario {
		case "missing-workdir":
			opts.WorkDir = ""
		case "memory-off":
			opts.Memory = false
		case "legacy-bypass":
			opts.CodexInitialPermissionMode = "bypass"
		case "builtin-profile":
			opts.CodexPermissionProfile = ":danger-full-access"
		}
		err := validateCodexIsolationOptions(opts)
		if (err == nil) != (scenario == "valid") {
			t.Fatalf("%s: %v", scenario, err)
		}
	}
}

func TestCrossPlatformCoverageCodexIsolationLaunchNeverUsesLegacySandbox(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	env, err := codexIsolatedEnvironment(nil)
	if err != nil {
		t.Fatal(err)
	}
	var command *exec.Cmd
	testseam.Swap(t, &codexExecCommandContext, func(ctx context.Context, bin string, args ...string) *exec.Cmd {
		command = exec.CommandContext(ctx, filepath.Join(t.TempDir(), "missing-codex"))
		command.Args = append([]string{bin}, args...)
		return command
	})
	if _, err := newCodexAppServerClient(context.Background(), "codex", env, t.TempDir()); err == nil {
		t.Fatal("missing executable unexpectedly started")
	}
	args := strings.Join(command.Args, " ")
	if strings.Contains(args, "sandbox") || strings.Contains(args, "--remote") || !strings.Contains(args, "features.apps=false") || !strings.Contains(args, "features.plugins=false") {
		t.Fatalf("unsafe isolated launch: %s", args)
	}
	if !reflect.DeepEqual(command.Env, env) {
		t.Fatal("isolated process inherited additional environment")
	}
}

func TestCrossPlatformCoverageCodexIsolationPathsAndEnvironment(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CODEX_HOME", home)
	t.Setenv("PRIVATE_TOKEN", "not-in-child")
	t.Setenv("CODEX_APP_SERVER_SOCKET", "personal-socket")
	env, err := codexIsolatedEnvironment(nil)
	if err != nil || strings.Contains(strings.Join(env, "\n"), "not-in-child") || strings.Contains(strings.Join(env, "\n"), "personal-socket") {
		t.Fatal("unsafe environment", err)
	}
	root := t.TempDir()
	a, err := codexIsolatedConversationDir(root, "../../outside")
	if err != nil {
		t.Fatal(err)
	}
	b, err := codexIsolatedConversationDir(root, "different")
	if err != nil || a == b || filepath.Dir(a) != filepath.Dir(b) {
		t.Fatal("conversations not isolated", err)
	}
	if _, err := codexIsolatedConversationDir(root, ""); err == nil {
		t.Fatal("empty identity accepted")
	}
	if err := os.Remove(a); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), a); err != nil {
		t.Skip("symlink unavailable:", err)
	}
	if _, err := codexIsolatedConversationDir(root, "../../outside"); err == nil {
		t.Fatal("symlink accepted")
	}
	opts := connectAgentOptions{WorkDir: root, Memory: true, CodexPermissionProfile: "strict"}
	path := codexIsolatedStorePath("codex-threads.json", opts)
	if path == "codex-threads.json" || !strings.Contains(path, "isolated-v1") {
		t.Fatal("legacy namespace reused")
	}
	opts.CodexPermissionProfile = "other"
	if path == codexIsolatedStorePath("codex-threads.json", opts) {
		t.Fatal("profile namespace reused")
	}
	t.Setenv("CODEX_HOME", "")
	if _, err := codexIsolatedEnvironment(nil); err == nil {
		t.Fatal("missing dedicated home accepted")
	}
}

func TestCrossPlatformCoverageCodexIsolationRejectsCorruptStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "isolated.json")
	if err := os.WriteFile(path, []byte(`broken`), 0600); err != nil {
		t.Fatal(err)
	}
	f := &codexAppServerForwarder{sessions: newCodexThreadSessions(path)}
	if !errors.Is(f.checkIsolatedStore(), errCodexSavedPermissions) {
		t.Fatal("corrupt history silently reset")
	}
}
