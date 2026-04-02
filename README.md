# dws-wukong

DingTalk Workspace CLI 的**悟空版**——在开源版 [`dingtalk-workspace-cli`](https://github.com/DingTalk-Real-AI/dingtalk-workspace-cli) 基础上，通过覆盖层（overlay）注入钉钉内部产品命令和运行模式定制，不 fork 代码，不改开源内核。

## 与开源版的关系

```
dingtalk-workspace-cli (开源核心)          dws-wukong (悟空覆盖层)
├── CLI 引擎、命令路由                       ├── main.go (3 行完成注入)
├── MCP 通信、认证、发现                     ├── wukong/hooks.go (钩子定义)
├── pkg/edition (扩展点合约)     ◄────────   ├── wukong/products/ (22 个产品命令)
├── pkg/cli (入口封装)                       ├── Makefile (交叉编译 + 签名)
└── internal/ (不可外部导入)                 └── dingtalk-workspace/ (Agent 技能包)
```

- **依赖方向**：`dws-wukong → dingtalk-workspace-cli`，单向，不可逆。
- **复用方式**：Go Module 依赖 + `pkg/edition.Override()` 注入，开源版 `internal/` 包不暴露、不复制。
- **开源版** 输出 `Edition: open`，**悟空版** 输出 `Edition: wukong`。

## 工作原理

悟空版的 `main.go` 只有三行核心逻辑：

```go
func main() {
    edition.Override(wukong.NewHooks())   // 注入悟空钩子
    cli.SetVersion(version, buildTime, gitCommit)
    os.Exit(cli.Execute())               // 启动开源版 CLI 引擎
}
```

`edition.Override` 注入的 `Hooks` 结构体覆盖以下行为：

| 扩展点 | 悟空版行为 |
|--------|-----------|
| `Name` / `ScenarioCode` | "wukong" / "com.dingtalk.scenario.wukong" |
| `IsEmbedded` / `HideAuthLogin` | 嵌入模式，隐藏 `dws auth login` |
| `AutoPurgeToken` | token 过期自动清除 |
| `ConfigDir` | 配置目录跟随二进制位置 |
| `MergeHeaders` | 注入 x-dingtalk-* 来源标识头 |
| `StaticServers` | 硬编码 MCP 端点，跳过 Market 服务发现 |
| `VisibleProducts` | 控制帮助页面显示的产品列表 |
| `RegisterExtraCommands` | 注册 22 个产品命令 + 3 个重定向别名 |

## 产品命令

22 个产品命令 + 3 个隐藏重定向：

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

## 目录结构

```
dws-wukong/
├── main.go                     # 入口：注入 → 设版本 → 启动
├── go.mod                      # 依赖 dingtalk-workspace-cli
├── Makefile                    # 编译、打包、签名
├── wukong/
│   ├── hooks.go                # Hooks 结构体组装
│   ├── endpoints.go            # 静态 MCP 端点列表
│   ├── config.go               # 配置目录策略
│   ├── auth.go                 # 认证错误处理
│   ├── headers.go              # HTTP 头注入
│   ├── hooks_test.go           # 钩子测试
│   └── products/               # 22 个产品命令实现
│       ├── register.go         # RegisterAll 入口
│       ├── aitable.go ... sheet.go
│       ├── helpers.go          # 共享辅助
│       ├── output.go           # 输出格式化
│       └── errors.go           # 错误处理
├── dingtalk-workspace/         # Agent 技能包 (Cursor/Gemini/.agent)
│   ├── SKILL.md
│   ├── references/
│   ├── scripts/
│   └── overlays/               # dev/real overlay
├── scripts/                    # build.sh, package.sh, sign.sh
├── _docs/                      # 内部架构文档
└── entitlements.plist          # macOS 签名 entitlements
```

## 快速开始

### 前置条件

- Go 1.25+
- 开源版源码（本地开发时需与 dws-wukong 同级放置）

### 本地目录布局

```
cli-workspace/              # 或其他父目录
├── github/
│   └── dingtalk-workspace-cli/   # 开源版
├── dws-wukong/                   # 本仓库
└── go.work                       # 可选，本地 workspace（不提交）
```

`go.work` 示例：

```
go 1.25.8

use (
    ./github/dingtalk-workspace-cli
    ./dws-wukong
)
```

### 编译运行

```bash
# 当前平台
make build
./dws version

# 指定版本号
make build VERSION=0.2.31
./dws version

# 交叉编译 6 平台
make build-all VERSION=0.2.31

# 测试
make test
```

### 版本输出示例

```
Version:        0.2.31
Edition:        wukong
Build:          2026-03-30T22:45:12+0800
Commit:         abc1234
Architecture:   MCP Dynamic Aggregation
Go:             go1.25.8
```

## 构建与发布

| Make 目标 | 用途 |
|-----------|------|
| `make build` | 当前平台单二进制 |
| `make build-all` | 交叉编译 6 平台 (darwin/linux/windows × amd64/arm64) |
| `make test` | 运行所有测试 |
| `make package` | 统一 release zip |
| `make dist` | 6 平台独立包 + macOS 签名 |
| `make real-platform` | 内部打包 (win + mac + macOS 签名) |
| `make public` | 公开发布包 |
| `make release` | 准备 git release + latest.json manifest |
| `make sign` | 独立 macOS 签名 |
| `make clean` | 清理构建产物 |

版本号通过 ldflags 注入：`-X main.version`、`-X main.buildTime`、`-X main.gitCommit`。

## 依赖管理

### 本地开发

`go.work` 让两个模块共享 workspace，修改开源版代码后直接 `make build` 即可在悟空版生效。

### 发布构建

`go.mod` 引用 GitHub 上的正式 tag：

```
require github.com/DingTalk-Real-AI/dingtalk-workspace-cli v0.2.31
```

去掉 `replace` 行，CI 从 GitHub 拉取确定版本构建。

### 发布流程

```
1. 开源版合并到 main，打 tag: git tag v0.2.31
2. 悟空版: go get github.com/DingTalk-Real-AI/dingtalk-workspace-cli@v0.2.31
3. 提交 go.mod + go.sum
4. make real-platform VERSION=0.2.31
```

## 测试

```bash
go vet ./...
go test ./...
```

开源版 CI（GitHub Actions）合并到 main 后会自动触发悟空版 GitLab CI。

## 架构文档

详细的双仓库架构设计见 [`_docs/wukong-overlay-architecture.md`](_docs/wukong-overlay-architecture.md)。
