# dws-wukong

DingTalk Workspace CLI 的**悟空版**——在开源版 [`dingtalk-workspace-cli`](https://github.com/DingTalk-Real-AI/dingtalk-workspace-cli) 基础上，通过覆盖层（overlay）注入钉钉产品命令、A2A 协议、行为授权和运行模式定制，不 fork 代码，不改开源内核。

## 与开源版的关系

```
dingtalk-workspace-cli (开源核心)             dws-wukong (悟空覆盖层)
├── CLI 引擎、命令路由                          ├── main.go (3 行完成注入)
├── MCP 通信、认证、发现                        ├── wukong/hooks.go (钩子定义，real/dev 双模式)
├── pkg/edition (扩展点合约)      ◄────────    ├── wukong/products/ (23 个产品命令)
├── pkg/cli (入口封装)                          ├── wukong/extensions/ (5 个内部扩展命令)
└── internal/ (不可外部导入)                    ├── wukong/a2a/ (A2A Agent 协作协议)
                                                ├── wukong/pat/ (行为授权管理)
                                                ├── dingtalk-workspace/ (Agent 技能包)
                                                ├── dingtalk-products-skills/ (23 个产品级 AI 工作流技能)
                                                ├── auto-test/ (CLI↔MCP / 端到端 / 模型↔CLI 测试)
                                                └── Makefile (交叉编译 + 签名 + 发布)
```

- **依赖方向**：`dws-wukong → dingtalk-workspace-cli`，单向，不可逆。
- **复用方式**：Go Module 依赖 + `pkg/edition.Override()` 注入，开源版 `internal/` 包不暴露、不复制。
- **开源版** 输出 `Edition: open`，**悟空版** 输出 `Edition: wukong`。

## 架构

### 工作原理

悟空版的 `main.go` 只有三行核心逻辑：

```go
func main() {
    edition.Override(wukong.NewHooks(buildMode))   // 注入悟空钩子
    cli.SetVersion(version, buildTime, gitCommit)
    os.Exit(cli.Execute())                         // 启动开源版 CLI 引擎
}
```

`buildMode` 在链接时通过 ldflags 注入，`real` 为嵌入模式（隐藏 auth login），`dev` 为独立 CLI 模式。

### Hooks 扩展点

`edition.Override` 注入的 `Hooks` 结构体覆盖以下行为：

| 扩展点 | 悟空版行为 |
|--------|-----------|
| `Name` / `ScenarioCode` | "wukong" / "com.dingtalk.scenario.wukong" |
| `IsEmbedded` / `HideAuthLogin` | `buildMode=real` 时启用嵌入模式并隐藏 `dws auth login` |
| `AutoPurgeToken` | `buildMode=real` 时 token 过期自动清除 |
| `ConfigDir` | 配置目录跟随二进制位置（`DWS_CONFIG_DIR` 或 exe 旁 `.dws`） |
| `MergeHeaders` | 注入 x-dingtalk-* 来源标识头 |
| `StaticServers` | 硬编码 MCP 端点，跳过 Market 服务发现 |
| `VisibleProducts` | 控制帮助页面显示的产品列表 |
| `RegisterExtraCommands` | 四层注册：产品命令 + 扩展命令 + PAT + A2A |
| `AfterPersistentPreRun` | A2A 运行时自动装配（当命令路径为 `dws a2a` 时） |

### 四层命令注册

```
registerAll(root)
  ├── products.RegisterProducts()      23 个标准产品命令
  ├── extensions.RegisterExtensions()   5 个内部扩展命令（hidden）
  ├── pat.RegisterCommands()            行为授权管理
  ├── a2a.RegisterCommands()            A2A 协作协议
  ├── 3 个隐藏重定向别名
  ├── RegisterHintSubCmds()             子命令提示
  └── RegisterCamelCaseAliases()        驼峰别名兼容
```

## 产品命令

23 个产品命令 + 3 个隐藏重定向：

| 命令 | 说明 | 命令 | 说明 |
|------|------|------|------|
| `aitable` | AI 多维表格 | `mail` | 邮箱 |
| `calendar` | 日历 | `ding` | DING 消息 |
| `contact` | 通讯录 | `devdoc` | 开放平台文档 |
| `todo` | 待办 | `attendance` | 考勤 |
| `doc` | 文档 | `conference` | 视频会议 |
| `chat` | 群聊与机器人 | `live` | 直播 |
| `oa` | 审批 | `aiapp` | AI 应用 |
| `minutes` | 妙记 | `finance` | 财务 |
| `report` | 日志 | `law` | 法务 |
| `drive` | 网盘 | `credit` | 征信 |
| `sheet` | 电子表格 | `docparse` | 文档解析 |
| `aidesign` | AI 设计 | | |

隐藏重定向：`bot` → `chat bot search`、`approval` → `oa approval ...`、`message` → `chat message ...`

### 内部扩展命令（hidden）

| 命令 | 说明 |
|------|------|
| `outbound-call` | AI 语音外呼 |
| `discovery` | 资讯与服务发现 |
| `ai-sincere-hire` | AI 诚聘 |
| `contract` | 智能合同 |
| `oa-plus` | 增强审批 |

### A2A 协议

A2A (Agent-to-Agent) 协议支持 Agent 发现与协作通信：

```bash
dws a2a agents list                          # 列出可用 Agent
dws a2a agents info --agent ai-search        # 查看 Agent 详情
dws a2a send --agent ai-search --text "..."  # 同步发送消息
dws a2a send --agent ai-search --text "..." --stream  # SSE 流式
```

### PAT 行为授权

```bash
dws pat chmod <scope>...    # 授予指定权限
```

## 目录结构

```
dws-wukong/
├── main.go                           # 入口：注入 → 设版本 → 启动
├── go.mod                            # 依赖 dingtalk-workspace-cli
├── Makefile                          # 编译、打包、签名、发布
├── wukong/
│   ├── hooks.go                      # Hooks 结构体组装（real/dev 双模式）
│   ├── register.go                   # 四层命令注册入口
│   ├── endpoints.go                  # 静态 MCP 端点列表 + 可见产品
│   ├── config.go                     # 配置目录策略
│   ├── auth.go                       # 认证错误处理
│   ├── headers.go                    # HTTP 头注入
│   ├── a2a_runtime.go                # A2A 运行时装配
│   ├── hooks_test.go                 # 钩子合约测试
│   ├── products/                     # 23 个产品命令实现
│   │   ├── register.go               # RegisterProducts 入口
│   │   ├── helpers.go                # 共享辅助（MCP 调用路由）
│   │   ├── output.go                 # 输出格式化（JSON/表格/KV）
│   │   ├── errors.go                 # 结构化错误处理与退出码
│   │   └── aitable.go ... sheet.go   # 各产品命令文件
│   ├── extensions/                   # 内部扩展命令
│   │   ├── register.go
│   │   └── vendors/dingtalk/         # 钉钉内部产品实现
│   ├── a2a/                          # A2A 协议实现
│   │   ├── client.go                 # HTTP + SSE 客户端（重试、网关）
│   │   ├── command.go                # Cobra 命令树
│   │   └── types.go                  # 协议类型定义
│   └── pat/                          # 行为授权管理
│       ├── pat.go                    # 命令入口
│       ├── chmod.go                  # chmod 实现
│       └── helpers.go                # PAT 辅助函数
├── dingtalk-workspace/               # Agent 技能包
│   ├── SKILL.md                      # 主技能定义
│   ├── references/                   # 产品参考文档
│   ├── scripts/                      # 辅助脚本
│   ├── overlays/                     # dev/real 环境 overlay
│   ├── tests/                        # 技能测试
│   └── package.json
├── dingtalk-products-skills/         # 23 个产品级 AI 工作流技能
│   ├── wukong-doc-business-plan/     # 商业计划书生成
│   ├── wukong-doc-competitive-analysis/ # 竞品分析
│   ├── dingtalk-daily-work-summary/  # 每日工作总结
│   ├── minutes-to-doc/              # 妙记转文档
│   └── ...                           # 更多专题技能
├── tool_skills/                      # 工具级技能
│   ├── dingtalk-open-spec/           # 开放平台 API 规范探索
│   ├── skill-evaluator/              # 技能质量评估
│   └── window-invocation-count/      # 调用计数
├── cli-troubleshoot/                 # CLI 诊断排障技能
├── auto-test/                        # 自动化测试体系
│   ├── cli_to_mcp/                   # CLI ↔ MCP 集成测试
│   ├── end_to_end/                   # 端到端评测
│   └── model_to_cli/                 # 模型 → CLI 测试
├── npm/                              # macOS npm 分发包
├── scripts/                          # build.sh, package.sh, sign.sh, deploy/
├── _docs/                            # 内部架构文档
└── entitlements.plist                # macOS 签名 entitlements
```

## 快速开始

### 前置条件

- Go 1.25+
- 开源版源码（本地开发时需与 dws-wukong 同级放置）

### 本地目录布局

```
cli-workspace/              # 或其他父目录
├── dingtalk-workspace-cli/      # 开源版
├── dws-wukong/                  # 本仓库
└── go.work                      # 可选，本地 workspace（不提交）
```

`go.work` 示例：

```
go 1.25.8

use (
    ./dingtalk-workspace-cli
    ./dws-wukong
)
```

### 编译运行

```bash
# 当前平台
make build
./dws version

# 指定版本号
make build VERSION=0.2.40

# dev 模式（默认）：支持 dws auth login
make build

# real 模式：嵌入模式，隐藏 auth login
make build BUILD_MODE=real

# 交叉编译 6 平台
make build-all VERSION=0.2.40

# 测试
make test
```

### 版本输出示例

```
Version:        0.2.40
Edition:        wukong
Build:          2026-04-09T15:30:00+0800
Commit:         abc1234
Architecture:   MCP Dynamic Aggregation
Go:             go1.25.8
```

## 构建与发布

| Make 目标 | 用途 |
|-----------|------|
| `make build` | 当前平台单二进制（默认 dev 模式） |
| `make build BUILD_MODE=real` | 当前平台嵌入模式二进制 |
| `make build-all` | 交叉编译 6 平台 (darwin/linux/windows x amd64/arm64) |
| `make test` | 运行所有 Go 测试 |
| `make package` | 统一 release zip |
| `make dist` | 6 平台独立包 + macOS 签名 |
| `make real-platform` | 内部打包 (win + mac + macOS 签名，real 模式) |
| `make public` | 公开发布包（dws/ + skill/ 分目录） |
| `make release` | 准备 git release + latest.json manifest |
| `make npm-pack` | 构建 macOS npm 分发包 |
| `make sign` | 独立 macOS 签名 |
| `make clean` | 清理构建产物 |

版本号通过 ldflags 注入：`-X main.version`、`-X main.buildTime`、`-X main.gitCommit`、`-X main.buildMode`。

## 依赖管理

### 本地开发

`go.work` 让两个模块共享 workspace，修改开源版代码后直接 `make build` 即可在悟空版生效。

### 发布构建

`go.mod` 引用 GitHub 上的正式 tag：

```
require github.com/DingTalk-Real-AI/dingtalk-workspace-cli v0.2.40
```

去掉 `replace` 行，CI 从 GitHub 拉取确定版本构建。

### 发布流程

```
1. 开源版合并到 main，打 tag: git tag v0.2.40
2. 悟空版: go get github.com/DingTalk-Real-AI/dingtalk-workspace-cli@v0.2.40
3. 提交 go.mod + go.sum
4. make real-platform VERSION=0.2.40
```

## 测试

### Go 单元测试

```bash
go vet ./...
go test ./...
```

### 自动化测试体系

`auto-test/` 提供三层测试覆盖：

| 层级 | 目录 | 用途 |
|------|------|------|
| CLI ↔ MCP | `auto-test/cli_to_mcp/` | 验证 CLI 命令到 MCP 工具调用的映射正确性 |
| 端到端 | `auto-test/end_to_end/` | 模拟真实用户场景的全链路评测 |
| 模型 → CLI | `auto-test/model_to_cli/` | 验证 LLM 能正确生成 CLI 调用 |

### CI

- GitLab CI：`go test ./...` + 多平台构建
- AoneCI：CLI-to-MCP 集成测试

## Agent 技能生态

### dingtalk-workspace 主技能包

位于 `dingtalk-workspace/`，面向 AI Agent（Cursor / Gemini / Claude 等）的 DingTalk CLI 使用指南，包含：
- `SKILL.md`：主技能定义
- `references/`：各产品的详细参考文档与最佳实践
- `scripts/`：辅助执行脚本
- `overlays/`：按环境（dev / real）覆盖参考内容

### 产品级 AI 工作流技能

`dingtalk-products-skills/` 包含 23 个专题技能，每个封装了特定业务场景的 AI 工作流：

| 分类 | 技能示例 |
|------|---------|
| 文档生成 | 商业计划书、竞品分析、行业研究、技术方案、周报 |
| 会议处理 | 妙记转文档、妙记转待办、会议纪要整理 |
| 日常效率 | 每日工作总结、一键邮件、审批进度、群聊摘要 |
| 内容创作 | 社交内容、营销方案、论文深度研究、话题趋势 |

### 工具级技能

`tool_skills/` 包含基础工具能力：
- `dingtalk-open-spec/`：探索未经 CLI 封装的原生 OpenAPI
- `skill-evaluator/`：技能质量评估打分
- `window-invocation-count/`：调用频次统计

## 架构文档

详细的双仓库架构设计见 [`_docs/wukong-overlay-architecture.md`](_docs/wukong-overlay-architecture.md)。
