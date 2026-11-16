// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package helpers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func validateCodexIsolationOptions(opts connectAgentOptions) error {
	name := opts.CodexPermissionProfile
	if !validMachineString(name) || strings.HasPrefix(name, ":") || opts.WorkDir == "" || !opts.Memory || opts.CodexInitialPermissionMode == "bypass" {
		return fmt.Errorf("invalid native Codex isolation options")
	}
	if _, err := codexIsolatedEnvironment(nil); err != nil {
		return err
	}
	_, err := filepath.EvalSymlinks(opts.WorkDir)
	return err
}

// A dedicated home is explicit; never attach to an inherited desktop service.
func codexIsolatedEnvironment(extra []string) ([]string, error) {
	values := map[string]string{}
	for _, entry := range append(os.Environ(), extra...) {
		key, value, _ := strings.Cut(entry, "=")
		values[key] = value
	}
	home := values["CODEX_HOME"]
	canonical, err := filepath.EvalSymlinks(home)
	if home == "" || !filepath.IsAbs(home) || err != nil {
		return nil, fmt.Errorf("Codex isolation requires an existing dedicated absolute CODEX_HOME")
	}
	personal, _ := os.UserHomeDir()
	personal = filepath.Join(personal, ".codex")
	if resolved, err := filepath.EvalSymlinks(personal); err == nil {
		personal = resolved
	}
	if canonical == personal {
		return nil, fmt.Errorf("Codex isolation cannot use the default personal CODEX_HOME")
	}
	env := []string{"CODEX_HOME=" + canonical, "DWS_CODEX_ISOLATION=1"}
	// Credentials, proxy injection, desktop sockets and tool configuration are not inherited.
	for _, key := range []string{"PATH", "HOME", "TMPDIR", "TMP", "TEMP", "LANG", "LC_ALL", "SYSTEMROOT", "SystemRoot", "WINDIR"} {
		if value := values[key]; value != "" {
			env = append(env, key+"="+value)
		}
	}
	return env, nil
}

func codexIsolatedStorePath(path string, opts connectAgentOptions) string {
	home, _ := filepath.EvalSymlinks(os.Getenv("CODEX_HOME"))
	root, _ := filepath.EvalSymlinks(opts.WorkDir)
	id := sha256.Sum256([]byte(home + "\x00" + root + "\x00" + opts.CodexPermissionProfile))
	return filepath.Join(filepath.Dir(path), fmt.Sprintf("codex-isolated-v1-%x.json", id))
}

func (f *codexAppServerForwarder) checkIsolatedStore() error {
	if f.sessions == nil {
		return errCodexSavedPermissions
	}
	if f.sessions.path == "" {
		return nil
	}
	data, err := os.ReadFile(f.sessions.path)
	if os.IsNotExist(err) {
		return nil
	}
	var entries map[string]string
	if err != nil || json.Unmarshal(data, &entries) != nil || entries == nil {
		return errCodexSavedPermissions
	}
	return nil
}

func codexIsolatedConversationDir(root, conversation string) (string, error) {
	if strings.TrimSpace(conversation) == "" {
		return "", errCodexSavedPermissions
	}
	base, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	base, err = filepath.Abs(base)
	if err != nil {
		return "", err
	}
	path := filepath.Join(base, fmt.Sprintf("session-%x", sha256.Sum256([]byte(conversation))))
	if err := os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
		return "", err
	}
	info, err := os.Lstat(path)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", errCodexSavedPermissions
	}
	canonical, err := filepath.EvalSymlinks(path)
	if err != nil || canonical != path || filepath.Dir(canonical) != base {
		return "", errCodexSavedPermissions
	}
	return path, nil
}

func (c *codexAppServerClient) isolatedProfileFilesystem(ctx context.Context, name, cwd string) (map[string]string, error) {
	id := c.requestID()
	if err := c.send(map[string]any{"id": id, "method": "config/read", "params": map[string]any{"cwd": cwd, "includeLayers": false}}); err != nil {
		return nil, err
	}
	raw, err := c.waitResponse(ctx, id)
	if err != nil {
		return nil, err
	}
	var response struct { Config map[string]json.RawMessage `json:"config"` }
	if json.Unmarshal(raw, &response) != nil {
		return nil, errCodexSavedPermissions
	}
	// No implicit personal MCPs or plugins, including project-local configuration.
	for _, key := range []string{"mcp_servers", "plugins"} {
		var entries map[string]json.RawMessage
		if value := response.Config[key]; len(value) > 0 && string(value) != "null" {
			if json.Unmarshal(value, &entries) != nil || len(entries) != 0 {
				return nil, fmt.Errorf("isolated Codex configuration must not enable MCP servers or plugins")
			}
		}
	}
	var profiles map[string]struct {
		Extends *string `json:"extends"`
		Roots []string `json:"workspace_roots"`
		Filesystem map[string]json.RawMessage `json:"filesystem"`
		Network struct { Enabled *bool `json:"enabled"` } `json:"network"`
	}
	if json.Unmarshal(response.Config["permissions"], &profiles) != nil {
		return nil, errCodexSavedPermissions
	}
	profile, exists := profiles[name]
	if !exists || profile.Extends != nil || len(profile.Roots) != 0 || profile.Network.Enabled == nil || *profile.Network.Enabled {
		return nil, errCodexSavedPermissions
	}
	filesystem := map[string]string{}
	for path, raw := range profile.Filesystem {
		if string(raw) == "null" { continue }
		var access string
		if json.Unmarshal(raw, &access) != nil || (access != "deny" && access != "read" && access != "write") {
			return nil, errCodexSavedPermissions
		}
		if strings.HasPrefix(path, ":") {
			switch path { case ":root", ":minimal", ":workspace_roots", ":tmpdir", ":slash_tmp": default: return nil, errCodexSavedPermissions }
		} else {
			canonical, err := filepath.EvalSymlinks(path)
			if !filepath.IsAbs(path) || err != nil || canonical != path || access == "write" { return nil, errCodexSavedPermissions }
		}
		filesystem[path] = access
	}
	for path, access := range map[string]string{":root":"deny", ":minimal":"read", ":workspace_roots":"write", ":tmpdir":"deny", ":slash_tmp":"deny"} {
		if filesystem[path] != access { return nil, errCodexSavedPermissions }
	}
	return filesystem, nil
}

func (c *codexAppServerClient) verifyIsolatedSavedSettings(ctx context.Context, thread, cwd string, filesystem map[string]string) error {
	saved, err := c.employeeSavedSettings(ctx, thread)
	if err != nil { return err }
	if saved.Cwd != cwd || !reflect.DeepEqual(saved.Roots, []string{cwd}) || string(saved.Approval) != `"on-request"` || saved.Reviewer != "user" {
		return errCodexSavedPermissions
	}
	params, err := saved.resumeParams(thread)
	if err != nil || params["sandbox"] != nil { return errCodexSavedPermissions }
	config, ok := params["config"].(map[string]any)
	if !ok { return errCodexSavedPermissions }
	profiles := config["permissions"].(map[string]any)
	actual := profiles[params["permissions"].(string)].(map[string]any)["filesystem"]
	if !reflect.DeepEqual(actual, filesystem) { return errCodexSavedPermissions }
	return nil
}
