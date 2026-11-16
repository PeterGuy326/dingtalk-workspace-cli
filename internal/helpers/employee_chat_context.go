// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/corecmd/contractfinal"
	apperrors "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/errors"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/output"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// EmployeeChatContext 只携带非敏感绑定元数据，正文始终经 stdin 传入。
// 它是发送范围约束，不替代 chat 自身的用户确认。
type EmployeeChatContext struct {
	AgentUUID       string `json:"agentUuid"`
	Channel         string `json:"channel"`
	BindingRevision uint64 `json:"bindingRevision"`
	ConversationID  string `json:"conversationId,omitempty"`
	MessageID       string `json:"messageId,omitempty"`
	ApprovalID      string `json:"approvalId,omitempty"`
}

// Resolve aliases and the default profile without loading credentials or trusting argv identity claims.
func registeredChatEmployee() (digitalEmployeeBinding, bool, error) {
	profile := strings.TrimSpace(auth.RuntimeProfile())
	b, err := deapChannelLoadBinding(deapConnectConfigDir(), profile)
	if err == nil && profile != "" {
		if b.DWSProfile != profile {
			return b, false, fmt.Errorf("chat binding profile mismatch")
		}
		return b, true, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return b, false, err
	}
	metadata, err := auth.ResolveProfileMetadataReadOnly(deapConnectConfigDir(), profile)
	if err != nil || metadata == nil {
		return b, false, err
	}
	exact := metadata.CorpID + ":" + metadata.UserID
	b, err = deapChannelLoadBinding(deapConnectConfigDir(), exact)
	if os.IsNotExist(err) {
		return b, false, nil
	}
	return b, err == nil, err
}

// EmployeeProfileIsRegistered uses the same local identity discovery as send validation.
func EmployeeProfileIsRegistered() (bool, error) {
	_, registered, err := registeredChatEmployee()
	return registered, err
}

// Employee writes enter here before preview or transport; only reviewed bound operations pass.
func ValidateEmployeeChatSendCommand(cmd *cobra.Command) error {
	var product *cobra.Command
	generic := false
	for node := cmd; node != nil; node = node.Parent() {
		if node.Name() == "api" || node.Name() == "mcp" || node.Name() == "plugin" || cmdutil.IsPluginSourced(node) || cmdutil.IsEnvelopeSourced(node) {
			generic = true
		}
		if node.Name() == "chat" {
			product = node
		}
	}
	if generic {
		employee, err := EmployeeProfileIsRegistered()
		if err != nil {
			return apperrors.NewValidation("cannot verify chat profile registration")
		}
		if employee {
			return apperrors.NewValidation("registered employees cannot use raw API, generic MCP or plugin command transports")
		}
		return nil
	}
	if product == nil {
		return nil
	}
	_, employee, err := registeredChatEmployee()
	if err != nil {
		return apperrors.NewValidation("cannot verify chat profile registration")
	}
	if !employee {
		return nil
	}
	if cmd.Name() == "+messages-send" || cmd.Name() == "+messages-reply" {
		if devAppStringFlag(cmd, "employee-context") == "" {
			return apperrors.NewValidation("registered employee chat requires --employee-context")
		}
		// The shared shortcut validation checks the bound context and resolved target.
		return nil
	}
	if cmd.Name() == "send-a2ui-card" || cmd.Name() == "update-a2ui-card" || cmd.Name() == "add-text-emotion" || cmd.Name() == "remove-text-emotion" {
		operation := "card-create"
		if cmd.Name() == "update-a2ui-card" {
			operation = "card-update"
		}
		if cmd.Name() == "add-text-emotion" || cmd.Name() == "remove-text-emotion" {
			operation = "feedback"
		}
		if err := ValidateEmployeeChatContext(cmd, operation, devAppStringFlag(cmd, "conversation-id"), devAppStringFlag(cmd, "message-id"), devAppStringFlag(cmd, "open-dingtalk-id"), ""); err != nil {
			return apperrors.NewValidation(err.Error())
		}
		return nil
	}
	if cmd.Name() == "create-text-emotion" {
		label := devAppStringFlag(cmd, "text")
		if label == devAppStringFlag(cmd, "emotion-name") && validEmployeeFeedbackLabel(label) && devAppStringFlag(cmd, "background-id") == "im_bg_1" {
			return nil
		}
	}
	if final, ok := contractfinal.RuntimeContractFinal(cmd); ok && final.Safety != nil && final.Safety.Effect == "read" {
		return nil
	}
	return apperrors.NewValidation("registered employees must use context-validated chat operations")
}

// ValidateEmployeeChatContext 在通用 chat 写入前重新读取绑定，禁止过期实例和非主管目标。
func ValidateEmployeeChatContext(cmd *cobra.Command, operation, conversation, reference, recipient, key string) error {
	raw := devAppStringFlag(cmd, "employee-context")
	if raw == "" {
		_, employee, err := registeredChatEmployee()
		if err != nil {
			return fmt.Errorf("cannot verify chat profile registration")
		}
		if employee {
			return fmt.Errorf("registered employee chat requires --employee-context")
		}
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
	if (operation == "reply" || operation == "operator-private") && (!bodyStdin || !waitDelivery || !validMachineString(key)) {
		return fmt.Errorf("employee chat requires --body-stdin, --wait-delivery and an idempotency key")
	}
	b, employee, err := registeredChatEmployee()
	if err != nil || !employee || b.AgentUUID != in.AgentUUID || bindingChannel(b) != in.Channel || b.BindingRevision != in.BindingRevision || employeeBindingState(b) != "bound" || employeeDesiredState(b) != "running" {
		return fmt.Errorf("employee chat binding is no longer authorized")
	}
	switch operation {
	case "card-update":
		biz := devAppStringFlag(cmd, "biz-id")
		if !validMachineString(biz) {
			return fmt.Errorf("employee card update requires its recorded bizId")
		}
		if in.ApprovalID != "" {
			if !validMachineString(in.ApprovalID) || strings.ContainsAny(in.ApprovalID, "/\\") || in.ApprovalID == "." || in.ApprovalID == ".." {
				return fmt.Errorf("invalid approval identity")
			}
			data, err := os.ReadFile(filepath.Join(digitalEmployeeRuntimeDir(b.DWSProfile), "approvals", in.ApprovalID+".json"))
			var record employeeApprovalRecord
			if err != nil || json.Unmarshal(data, &record) != nil || record.ID != in.ApprovalID || record.BizID != biz {
				return fmt.Errorf("employee approval card is not owned by this runtime")
			}
			return nil
		}
		if !validMachineString(in.ConversationID) || !validMachineString(in.MessageID) {
			return fmt.Errorf("employee card update requires its source")
		}
		runtime := employeeRuntime{dir: digitalEmployeeRuntimeDir(b.DWSProfile)}
		data, err := os.ReadFile(runtime.recordPath(employeeEvent{ConversationID: in.ConversationID, MessageID: in.MessageID}))
		var record employeeTaskRecord
		if err != nil || json.Unmarshal(data, &record) != nil || record.ConversationID != in.ConversationID || record.MessageID != in.MessageID || record.BizID != biz || record.Status != "card_updating" {
			return fmt.Errorf("employee reply card is not owned by this source")
		}
		return nil
	case "feedback":
		if conversation != in.ConversationID || reference != in.MessageID || !validMachineString(conversation) || !validMachineString(reference) {
			return fmt.Errorf("employee feedback source mismatch")
		}
		feedback := employeeFeedback{dir: digitalEmployeeRuntimeDir(b.DWSProfile)}
		data, err := os.ReadFile(feedback.path(employeeEvent{ConversationID: conversation, MessageID: reference}))
		var record employeeFeedbackRecord
		if err != nil || json.Unmarshal(data, &record) != nil || record.ConversationID != conversation || record.MessageID != reference {
			return fmt.Errorf("employee feedback source not recorded")
		}
		for _, emotion := range record.Owned {
			if emotion.ID == devAppStringFlag(cmd, "emotion-id") && emotion.Label == devAppStringFlag(cmd, "emotion-name") && emotion.Label == devAppStringFlag(cmd, "text") && validEmployeeFeedbackLabel(emotion.Label) && devAppStringFlag(cmd, "background-id") == "im_bg_1" {
				return nil
			}
		}
		return fmt.Errorf("employee feedback does not match its owned status")
	case "card-create":
		if recipient != "" && recipient == b.OperatorOpenDingTalkID && conversation == "" {
			return nil
		}
		if !validMachineString(in.ConversationID) || !validMachineString(in.MessageID) {
			return fmt.Errorf("employee card requires an accepted source message")
		}
		runtime := employeeRuntime{dir: digitalEmployeeRuntimeDir(b.DWSProfile)}
		data, err := os.ReadFile(runtime.recordPath(employeeEvent{ConversationID: in.ConversationID, MessageID: in.MessageID}))
		var record employeeTaskRecord
		if err != nil || json.Unmarshal(data, &record) != nil || record.ConversationID != in.ConversationID || record.MessageID != in.MessageID || record.Status != "card_creating" || record.SenderID == "" {
			return fmt.Errorf("employee card source is not authorized")
		}
		if (recipient == "" && conversation == record.ConversationID) || (conversation == "" && recipient == record.SenderID) {
			return nil
		}
		return fmt.Errorf("employee card target does not match its source")
	case "reply":
		if !validMachineString(conversation) || !validMachineString(reference) {
			return fmt.Errorf("employee reply requires an explicit conversation and reference message")
		}
		runtime := employeeRuntime{dir: digitalEmployeeRuntimeDir(b.DWSProfile)}
		data, err := os.ReadFile(runtime.recordPath(employeeEvent{ConversationID: conversation, MessageID: reference}))
		var record employeeTaskRecord
		if err != nil || json.Unmarshal(data, &record) != nil || record.ConversationID != conversation || record.MessageID != reference {
			return fmt.Errorf("employee reply must reference an accepted source message")
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
