// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	apperrors "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/errors"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/output"
	"github.com/spf13/cobra"
)

// EmployeeChatContext 只携带非敏感绑定元数据，正文始终经 stdin 传入。
// 它是发送范围约束，不替代 chat 自身的用户确认。
type EmployeeChatContext struct {
	AgentUUID       string `json:"agentUuid"`
	Channel         string `json:"channel"`
	BindingRevision uint64 `json:"bindingRevision"`
	ConversationID string `json:"conversationId,omitempty"`
	MessageID string `json:"messageId,omitempty"`
}

// Resolve aliases and the default profile without loading credentials or trusting argv identity claims.
func registeredChatEmployee() (digitalEmployeeBinding, bool, error) {
	profile := strings.TrimSpace(auth.RuntimeProfile())
	b, err := deapChannelLoadBinding(deapConnectConfigDir(), profile)
	if err == nil && profile != "" {
		if b.DWSProfile != profile { return b, false, fmt.Errorf("chat binding profile mismatch") }
		return b, true, nil
	}
	if err != nil && !os.IsNotExist(err) { return b, false, err }
	metadata, err := auth.ResolveProfileMetadataReadOnly(deapConnectConfigDir(), profile)
	if err != nil || metadata == nil { return b, false, err }
	exact := metadata.CorpID + ":" + metadata.UserID
	b, err = deapChannelLoadBinding(deapConnectConfigDir(), exact)
	if os.IsNotExist(err) { return b, false, nil }
	return b, err == nil, err
}

// Every reviewed chat entry that creates/forwards a message enters here before dry-run or transport.
// Unsupported employee transports are closed; existing card updates and reactions remain separate operations.
func ValidateEmployeeChatSendCommand(cmd *cobra.Command) error {
	var product *cobra.Command
	for node := cmd; node != nil; node = node.Parent() { if node.Name() == "chat" { product = node; break } }
	if product == nil { return nil }
	switch cmd.Name() {
	case "send", "send-by-bot", "send-by-webhook", "send-card", "send-a2ui-card", "reply", "forward", "combine-forward", "forward-topic",
		"+dm", "+messages-send", "+messages-reply", "+messages-send-by-bot", "+messages-batch-send-by-bot", "+messages-send-by-webhook", "+messages-send-card", "+messages-forward", "+messages-combine-forward", "+messages-forward-topic":
	default: return nil
	}
	_, employee, err := registeredChatEmployee()
	if err != nil { return apperrors.NewValidation("cannot verify chat profile registration") }
	if !employee { return nil }
	if cmd.Name() == "+messages-send" || cmd.Name() == "+messages-reply" {
		if devAppStringFlag(cmd, "employee-context") == "" { return apperrors.NewValidation("registered employee chat requires --employee-context") }
		// The shared shortcut validation checks the bound context and resolved target.
		return nil
	}
	if cmd.Name() == "send-a2ui-card" {
		if err := ValidateEmployeeChatContext(cmd, "card-create", devAppStringFlag(cmd,"conversation-id"), "", devAppStringFlag(cmd,"open-dingtalk-id"), ""); err != nil { return apperrors.NewValidation(err.Error()) }
		return nil
	}
	return apperrors.NewValidation("registered employees must use context-validated +messages-send or +messages-reply")
}

// ValidateEmployeeChatContext 在通用 chat 写入前重新读取绑定，禁止过期实例和非主管目标。
func ValidateEmployeeChatContext(cmd *cobra.Command, operation, conversation, reference, recipient, key string) error {
	raw := devAppStringFlag(cmd, "employee-context")
	if raw == "" {
		_, employee, err := registeredChatEmployee()
		if err != nil { return fmt.Errorf("cannot verify chat profile registration") }
		if employee { return fmt.Errorf("registered employee chat requires --employee-context") }
		return nil
	}
	var in EmployeeChatContext
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if len(raw) > 4096 || decoder.Decode(&in) != nil || decoder.Decode(&struct{}{}) != io.EOF || !validMachineString(in.AgentUUID) || !validMachineString(in.Channel) {
		return fmt.Errorf("invalid employee chat context")
	}
	bodyStdin, _ := cmd.Flags().GetBool("body-stdin")
	waitDelivery, _ := cmd.Flags().GetBool("wait-delivery")
	if operation != "card-create" && (!bodyStdin || !waitDelivery || !validMachineString(key)) {
		return fmt.Errorf("employee chat requires --body-stdin, --wait-delivery and an idempotency key")
	}
	b, employee, err := registeredChatEmployee()
	if err != nil || !employee || b.AgentUUID != in.AgentUUID || bindingChannel(b) != in.Channel || b.BindingRevision != in.BindingRevision || employeeBindingState(b) != "bound" || employeeDesiredState(b) != "running" {
		return fmt.Errorf("employee chat binding is no longer authorized")
	}
	switch operation {
	case "card-create":
		if recipient != "" && recipient == b.OperatorOpenDingTalkID && conversation == "" { return nil }
		if !validMachineString(in.ConversationID) || !validMachineString(in.MessageID) { return fmt.Errorf("employee card requires an accepted source message") }
		runtime := employeeRuntime{dir: digitalEmployeeRuntimeDir(b.DWSProfile)}
		data, err := os.ReadFile(runtime.recordPath(employeeEvent{ConversationID:in.ConversationID, MessageID:in.MessageID}))
		var record employeeTaskRecord
		if err != nil || json.Unmarshal(data,&record) != nil || record.ConversationID != in.ConversationID || record.MessageID != in.MessageID || record.Status != "card_creating" || record.SenderID == "" { return fmt.Errorf("employee card source is not authorized") }
		if (recipient == "" && conversation == record.ConversationID) || (conversation == "" && recipient == record.SenderID) { return nil }
		return fmt.Errorf("employee card target does not match its source")
	case "reply":
		if !validMachineString(conversation) || !validMachineString(reference) {
			return fmt.Errorf("employee reply requires an explicit conversation and reference message")
		}
	case "operator-private":
		if !validMachineString(recipient) || b.OperatorOpenDingTalkID != recipient {
			return fmt.Errorf("employee chat target does not match the operator fixed by connect")
		}
	default:
		return fmt.Errorf("unsupported employee chat operation")
	}
	return nil
}


// WriteChatDelivery 保持通用 chat 与兼容 Channel 的机器回执一致。
func WriteChatDelivery(cmd *cobra.Command, result map[string]any, conversation, key string) error {
	delivery, err := ResolveChatDelivery(cmd.Context(), result, conversation, key)
	if err != nil {
		return err
	}
	if output.UsesUnifiedResult(cmd) {
		return output.StoreResult(cmd.Context(), output.Success(delivery))
	}
	return output.WriteCommandPayload(cmd, delivery, output.FormatJSON)
}
