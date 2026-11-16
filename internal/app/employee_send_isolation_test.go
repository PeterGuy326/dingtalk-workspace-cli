// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package app

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/helpers"
	"github.com/spf13/cobra"
)

func TestEmployeeSendGuardCoversAssembledChatTree(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("DWS_CONFIG_DIR", dir)
	old := auth.RuntimeProfile()
	auth.SetRuntimeProfile("corp:employee")
	t.Cleanup(func() { auth.SetRuntimeProfile(old) })
	path := filepath.Join(dir, "digital-employees", fmt.Sprintf("%x.json", sha256.Sum256([]byte("corp:employee"))))
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	binding := map[string]any{"schemaVersion": 1, "dwsProfile": "corp:employee", "agentUuid": "employee", "operatorOpenDingTalkId": "operator", "channel": "codex", "bindingRevision": 3, "bindingState": "bound", "desiredState": "running"}
	data, _ := json.Marshal(binding)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	root := newRootPresentationCommand()
	chat, _, err := root.Find([]string{"chat"})
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		for _, child := range cmd.Commands() {
			walk(child)
		}
		if !cmd.Runnable() || cmd.HasSubCommands() {
			return
		}
		name := cmd.Name()
		// These are new-message transports; status queries and sender searches are reads.
		outbound := name == "send" || name == "reply" || name == "forward" || name == "combine-forward" || name == "forward-topic" || strings.HasPrefix(name, "send-") || strings.HasPrefix(name, "+messages-send") || strings.HasPrefix(name, "+messages-batch-send") || name == "+messages-reply" || strings.HasPrefix(name, "+messages-forward") || name == "+messages-combine-forward"
		if !outbound {
			return
		}
		count++
		if err := helpers.ValidateEmployeeChatSendCommand(cmd); err == nil {
			t.Errorf("outbound leaf escaped employee guard: %s", cmd.CommandPath())
		}
	}
	walk(chat)
	if count < 20 {
		t.Fatalf("unexpectedly incomplete chat tree: %d outbound leaves", count)
	}
}

func TestEmployeeRawCallsStopAtPublicPreflight(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("DWS_CONFIG_DIR", dir)
	path := filepath.Join(dir, "digital-employees", fmt.Sprintf("%x.json", sha256.Sum256([]byte("corp:employee"))))
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"schemaVersion":1,"dwsProfile":"corp:employee","agentUuid":"employee","operatorOpenDingTalkId":"operator","channel":"codex","bindingRevision":3,"bindingState":"bound","desiredState":"running"}`), 0600); err != nil {
		t.Fatal(err)
	}
	old := auth.RuntimeProfile()
	t.Cleanup(func() { auth.SetRuntimeProfile(old) })
	for _, args := range [][]string{
		{"api", "GET", "/v1.0/synthetic", "--dry-run"},
		{"mcp", "published", "invoke", "2480", "send_personal_message", "--params", "{}", "--dry-run"},
		{"chat", "message", "send", "--open-dingtalk-id", "other", "--content", "synthetic", "--dry-run"},
	} {
		root := newRootCommandWithMode(context.Background(), nil, false, false, false)
		root.SetOut(io.Discard)
		root.SetErr(io.Discard)
		root.SetArgs(append([]string{"--profile", "corp:employee"}, args...))
		if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "registered employee") {
			t.Fatalf("%v escaped preflight: %v", args, err)
		}
	}
}
