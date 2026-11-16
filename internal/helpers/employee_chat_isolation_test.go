// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");

package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/auth"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/testseam"
	"github.com/spf13/cobra"
)

func TestCrossPlatformCoverageEmployeeCommonSendGuard(t *testing.T) {
	old:=auth.RuntimeProfile();auth.SetRuntimeProfile("corp:employee");t.Cleanup(func(){auth.SetRuntimeProfile(old)})
	b:=digitalEmployeeBinding{DWSProfile:"corp:employee",AgentUUID:"employee",Channel:"codex",BindingRevision:3,BindingState:"bound",DesiredState:"running",OperatorOpenDingTalkID:"operator"}
	testseam.Swap(t,&deapChannelLoadBinding,func(string,string)(digitalEmployeeBinding,error){return b,nil})
	for _,name:=range []string{"send","send-by-bot","send-by-webhook","send-card","send-a2ui-card","reply","forward","combine-forward","forward-topic","+dm","+messages-send","+messages-reply","+messages-send-by-bot","+messages-batch-send-by-bot","+messages-send-by-webhook","+messages-send-card","+messages-forward","+messages-combine-forward","+messages-forward-topic"} {
		for _,dry:=range []bool{false,true} {
			cmd:=&cobra.Command{Use:name};cmd.Flags().Bool("dry-run",dry,"")
			root:=&cobra.Command{Use:"chat"};root.AddCommand(cmd)
			if err:=ValidateEmployeeChatSendCommand(cmd);err==nil { t.Fatalf("%s dry=%v accepted missing context",name,dry) }
		}
	}
	// No binding is the ordinary-account path, regardless of the caller's display name.
	t.Setenv("DWS_CONFIG_DIR",t.TempDir())
	testseam.Swap(t,&deapChannelLoadBinding,func(string,string)(digitalEmployeeBinding,error){return digitalEmployeeBinding{},os.ErrNotExist})
	cmd:=&cobra.Command{Use:"send"};root:=&cobra.Command{Use:"chat"};root.AddCommand(cmd)
	if err:=ValidateEmployeeChatSendCommand(cmd);err!=nil { t.Fatal("ordinary profile changed",err) }
}

func TestCrossPlatformCoverageEmployeeNativeCardSource(t *testing.T) {
	t.Setenv("DWS_CONFIG_DIR",t.TempDir())
	old:=auth.RuntimeProfile();auth.SetRuntimeProfile("corp:employee");t.Cleanup(func(){auth.SetRuntimeProfile(old)})
	b:=digitalEmployeeBinding{DWSProfile:"corp:employee",AgentUUID:"employee",Channel:"codex",BindingRevision:3,BindingState:"bound",DesiredState:"running",OperatorOpenDingTalkID:"operator"}
	testseam.Swap(t,&deapChannelLoadBinding,func(string,string)(digitalEmployeeBinding,error){return b,nil})
	r:=employeeRuntime{dir:digitalEmployeeRuntimeDir(b.DWSProfile)}
	path:=r.recordPath(employeeEvent{ConversationID:"conversation",MessageID:"message"})
	if err:=os.MkdirAll(filepath.Dir(path),0700);err!=nil { t.Fatal(err) }
	if err:=writeEmployeeJSON(path,employeeTaskRecord{ConversationID:"conversation",MessageID:"message",SenderID:"sender",Status:"card_creating"});err!=nil { t.Fatal(err) }
	for _,scenario:=range []string{"operator","sender","group","wrong-recipient","wrong-group","missing-context","forged-context","stale-context","missing-source"} {
		metadata:=EmployeeChatContext{AgentUUID:b.AgentUUID,Channel:b.Channel,BindingRevision:3,ConversationID:"conversation",MessageID:"message"}
		conversation,recipient:="","sender"
		switch scenario {
		case "operator":recipient="operator";metadata.ConversationID="";metadata.MessageID=""
		case "group":conversation="conversation";recipient=""
		case "wrong-recipient":recipient="other"
		case "wrong-group":conversation="other";recipient=""
		case "forged-context":metadata.AgentUUID="other"
		case "stale-context":metadata.BindingRevision++
		case "missing-source":metadata.MessageID="absent"
		}
		data,_:=json.Marshal(metadata)
		if scenario=="missing-context" { data=nil }
		cmd:=&cobra.Command{};cmd.Flags().String("employee-context",string(data),"")
		err:=ValidateEmployeeChatContext(cmd,"card-create",conversation,"",recipient,"")
		want:=scenario=="operator" || scenario=="sender" || scenario=="group"
		if (err==nil)!=want { t.Fatalf("%s: %v",scenario,err) }
	}
}

func TestCrossPlatformCoverageEmployeeDefaultAndAliasProfile(t *testing.T) {
	dir:=t.TempDir();t.Setenv("DWS_CONFIG_DIR",dir)
	if err:=os.WriteFile(filepath.Join(dir,"profiles.json"),[]byte(`{"version":3,"currentProfile":"corp:employee","profiles":[{"name":"employee-alias","corpId":"corp","userId":"employee"},{"name":"person","corpId":"corp","userId":"person"}]}`),0600);err!=nil { t.Fatal(err) }
	b:=digitalEmployeeBinding{SchemaVersion:1,DWSProfile:"corp:employee",AgentUUID:"employee",Channel:"codex",BindingRevision:3,BindingState:"bound",DesiredState:"running",OperatorOpenDingTalkID:"operator"}
	if err:=os.MkdirAll(filepath.Dir(digitalEmployeeBindingPath(dir,b.DWSProfile)),0700);err!=nil { t.Fatal(err) }
	data,_:=json.Marshal(b)
	if err:=os.WriteFile(digitalEmployeeBindingPath(dir,b.DWSProfile),data,0600);err!=nil { t.Fatal(err) }
	old:=auth.RuntimeProfile();t.Cleanup(func(){auth.SetRuntimeProfile(old)})
	for _,profile:=range []string{"","employee-alias","corp:employee","person"} {
		auth.SetRuntimeProfile(profile)
		cmd:=&cobra.Command{}
		err:=ValidateEmployeeChatContext(cmd,"operator-private","","","operator","key")
		if (err==nil)!=(profile=="person") { t.Fatalf("profile %q: %v",profile,err) }
	}
}
