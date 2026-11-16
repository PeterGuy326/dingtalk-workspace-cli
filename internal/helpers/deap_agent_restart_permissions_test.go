// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package helpers

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/testseam"
)

// Exercise the saved adapter and disk mapping through the production factory,
// as the daemon does after restart; no host process or model is started.
func TestCrossPlatformCoverageEmployeeCodexRestartPreservesPermissions(t *testing.T) {
	for _, mode := range []string{"ask", "native", "native-wider"} {
		for _, managed := range []bool{false, true} {
			t.Run(mode+"/managed="+jsonScalar(managed), func(t *testing.T) {
				clearChannelEnv(t)
				cfg := employeeBoundaryConfig(t)
				t.Setenv("DWS_CONFIG_DIR", deapConnectConfigDir())
				home := t.TempDir()
				t.Setenv("CODEX_HOME", home)
				binDir := t.TempDir()
				if err := writeExecStub(binDir, "codex"); err != nil {
					t.Fatal(err)
				}
				t.Setenv("PATH", binDir)
				t.Setenv("DWS_CONNECT_NO_INSTALL", "1")
				cfg.Binding.Channel = "codex"
				cfg.Options = connectAgentOptions{WorkDir: t.TempDir(), Memory: true, CodexInitialPermissionMode: "ask"}
				if mode != "ask" {
					cfg.Options.CodexPermissionProfile = "strict"
				}
				if err := saveDigitalEmployeeBinding(deapConnectConfigDir(), cfg.Binding); err != nil {
					t.Fatal(err)
				}
				adapter := filepath.Join(digitalEmployeeRuntimeDir(cfg.Binding.DWSProfile), "adapter.json")
				if err := writeEmployeeJSON(adapter, cfg); err != nil {
					t.Fatal(err)
				}
				ctx := context.Background()
				key := "conversation"
				if managed {
					ctx = context.WithValue(ctx, agentProgressKey{}, func(agentProgress) {})
					key = codexManagedConversationKey(key)
				}
				var mapping, cwd string
				for phase := 0; phase < 2; phase++ {
					loaded, err := loadDigitalEmployeeConfig(cfg.Binding.DWSProfile)
					if err != nil {
						t.Fatal(err)
					}
					if loaded.Options.CodexInitialPermissionMode != "ask" || loaded.Options.CodexPermissionProfile != cfg.Options.CodexPermissionProfile {
						t.Fatal("saved permission selection changed across restart")
					}
					fwd, err := digitalEmployeeNewForwarder(ctx, loaded)
					if err != nil {
						t.Fatal(err)
					}
					f := fwd.(*codexAppServerForwarder)
					if phase == 0 {
						mapping, cwd = f.sessions.path, f.cwd()
						if mode != "ask" {
							cwd, err = codexIsolatedConversationDir(cfg.Options.WorkDir, f.isolationScope+"\x00conversation")
							if err != nil {
								t.Fatal(err)
							}
							if !strings.Contains(filepath.Base(mapping), "codex-isolated-v1-") {
								t.Fatal("old history namespace reused")
							}
						}
					} else if f.sessions.path != mapping || f.sessions.threadID(key) != "thread" {
						t.Fatal("restart lost the existing conversation or permission namespace")
					}
					messages := []codexRPCMessage{{ID: codexIntPtr(1), Result: json.RawMessage(`{}`)}}
					id := 2
					if mode != "ask" {
						config, _ := json.Marshal(map[string]any{"config": map[string]any{"permissions": map[string]any{"strict": map[string]any{"filesystem": isolatedFilesystemFixture(), "network": map[string]bool{"enabled": false}}}}})
						messages = append(messages, codexRPCMessage{ID: codexIntPtr(id), Result: config})
						id++
					}
					result := map[string]any{"thread": map[string]string{"id": "thread"}, "cwd": cwd, "runtimeWorkspaceRoots": []string{cwd}, "approvalPolicy": "on-request", "approvalsReviewer": "user", "activePermissionProfile": map[string]string{"id": "strict"}}
					if phase == 1 {
						saved := isolatedSavedFixture(cwd)
						if mode == "ask" {
							saved.Profile = nil
							saved.Sandbox = json.RawMessage(`{"type":"workspace-write","writable_roots":[],"network_access":false,"exclude_tmpdir_env_var":true,"exclude_slash_tmp":true}`)
						} else if mode == "native-wider" {
							saved.Profile = json.RawMessage(`{"type":"disabled"}`)
						}
						meta, _ := json.Marshal(map[string]any{"type": "session_meta", "payload": map[string]string{"id": "thread"}})
						turn, _ := json.Marshal(map[string]any{"type": "turn_context", "payload": saved})
						path := filepath.Join(home, "rollout.jsonl")
						if err := os.WriteFile(path, append(append(meta, '\n'), append(turn, '\n')...), 0600); err != nil {
							t.Fatal(err)
						}
						read, _ := json.Marshal(map[string]any{"thread": map[string]string{"id": "thread", "path": path}})
						messages = append(messages, codexRPCMessage{ID: codexIntPtr(id), Result: read})
						id++
						if mode == "ask" {
							_ = json.Unmarshal(codexPermissionResumeFixture(t, "thread", saved), &result)
						}
					}
					raw, _ := json.Marshal(result)
					messages = append(messages, codexRPCMessage{ID: codexIntPtr(id), Result: raw}, codexRPCMessage{Method: "turn/completed", Params: json.RawMessage(`{"threadId":"thread","turn":{"status":"completed","items":[{"type":"agentMessage","text":"done"}]}}`)})
					writer := &bufferWriteCloser{}
					client := unitCodexClient(writer, messages...)
					testseam.Swap(t, &codexNewAppServerClient, func(context.Context, string, []string, string) (*codexAppServerClient, error) { return client, nil })
					_, err = f.forward(ctx, "conversation", "hello")
					blocked := phase == 1 && mode == "native-wider"
					if (err != nil) != blocked {
						t.Fatalf("phase=%d: %v", phase, err)
					}
					turns, resumes := 0, 0
					for _, request := range decodeCodexRequests(t, writer.String()) {
						p := request.Params
						if request.Method == "thread/start" || request.Method == "thread/resume" {
							if phase == 1 && request.Method == "thread/start" {
								t.Fatal("restart discarded existing thread")
							}
							if request.Method == "thread/resume" {
								resumes++
							}
							if p["approvalPolicy"] != "on-request" || p["approvalsReviewer"] != "user" {
								t.Fatal("restart regressed approval policy", p)
							}
							if mode == "ask" && p["sandbox"] != "workspace-write" {
								t.Fatal("ask regressed to read-only", p)
							}
						}
						if request.Method == "turn/start" {
							turns++
						}
						if mode != "ask" && (request.Method == "thread/start" || request.Method == "thread/resume" || request.Method == "turn/start") {
							if p["permissions"] != "strict" || p["sandbox"] != nil || p["config"] != nil || p["cwd"] != cwd || !reflect.DeepEqual(p["runtimeWorkspaceRoots"], []any{cwd}) {
								t.Fatal("native profile boundary changed", p)
							}
						} else if mode == "ask" && request.Method == "turn/start" {
							for _, name := range []string{"approvalPolicy", "approvalsReviewer", "sandbox", "sandboxPolicy", "permissions"} {
								if p[name] != nil {
									t.Fatal("turn overwrote restored permissions", p)
								}
							}
						}
					}
					if blocked {
						if turns != 0 || resumes != 0 || f.sessions.threadID(key) != "thread" {
							t.Fatal("wider saved permissions resumed or destroyed history")
						}
					} else if turns != 1 || (phase == 1 && resumes != 1) {
						t.Fatal("expected lifecycle request missing")
					}
				}
			})
		}
	}
}
