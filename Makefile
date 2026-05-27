.PHONY: build build-all test clean npm-pack package dist real-platform sign public release help sync-upstream integration-regression bundle bundle-platform dump-commands local-platform

VERSION ?= 0.2.65
BUILD_TIME := $(shell date '+%Y-%m-%dT%H:%M:%S%z')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_MODE ?= dev

LDFLAGS := -ldflags "\
	-s -w \
	-X main.version=$(VERSION) \
	-X main.buildTime=$(BUILD_TIME) \
	-X main.gitCommit=$(GIT_COMMIT) \
	-X main.buildMode=$(BUILD_MODE)"

BUILDFLAGS := -trimpath -buildmode=pie

DIST := dist
SIGN_SERVER ?= http://30.45.44.83:7770

## build: Build for current platform
build:
	go build $(BUILDFLAGS) $(LDFLAGS) -o dws .

## build-all: Build for all 6 platforms (darwin/linux/windows x amd64/arm64)
build-all: clean
	@mkdir -p $(DIST)
	@echo "🔨 Building v$(VERSION) for all platforms..."

	GOOS=darwin  GOARCH=amd64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(DIST)/dws-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(DIST)/dws-darwin-arm64 .
	GOOS=linux   GOARCH=amd64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(DIST)/dws-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(DIST)/dws-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(DIST)/dws-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(DIST)/dws-windows-arm64.exe .

	@echo ""
	@echo "✅ Built 6 binaries in $(DIST)/:"
	@ls -lh $(DIST)/
	@echo ""
	@echo "📦 v$(VERSION) | commit $(GIT_COMMIT) | $(BUILD_TIME)"

## test: Run all tests
test:
	go test ./...

## integration-regression: Pre-release checks (legacy command paths + PAT), no product code changes
integration-regression:
	@chmod +x integration/regression.sh 2>/dev/null || true
	@./integration/regression.sh

## clean: Remove build artifacts
clean:
	rm -f dws dws-darwin-* dws-linux-* dws-windows-* scripts/dump-commands/dump-commands
	rm -rf $(DIST)
	rm -rf npm/bin

## dump-commands: Build CSV exporter for full cobra tree → scripts/dump-commands/dump-commands
dump-commands:
	go build $(BUILDFLAGS) -o scripts/dump-commands/dump-commands ./scripts/dump-commands

## npm-pack: Prepare npm package with pre-built binary
npm-pack: build-all
	mkdir -p npm/bin
	@ARCH=$$(uname -m); \
	if [ "$$ARCH" = "arm64" ]; then \
		cp $(DIST)/dws-darwin-arm64 npm/bin/dws; \
	else \
		cp $(DIST)/dws-darwin-amd64 npm/bin/dws; \
	fi
	chmod +x npm/bin/dws
	cd npm && npm pack
	@echo "✅ npm package created in npm/"

# ── Workspace & packaging ────────────────────────────────────────────

WORKSPACE_DIR := ./dingtalk-workspace
TARGET := target
PKG_NAME := dws_res
ENTITLEMENTS := entitlements.plist
BUNDLED_PLUGINS ?= bundled-plugins

# $(1) = output zip path, $(2) = workspace variant: dev | real (real applies overlays/real)
# $(3) = platform: darwin | windows (for plugin binary selection; omit to exclude plugins)
define build-workspace-zip
	@if [ -d "$(WORKSPACE_DIR)" ]; then \
		_ws_tmp=$$(mktemp -d); \
		cp -r "$(WORKSPACE_DIR)/SKILL.md" "$(WORKSPACE_DIR)/references" "$(WORKSPACE_DIR)/scripts" "$$_ws_tmp/"; \
		if [ "$(strip $(2))" != "dev" ] && [ -d "$(WORKSPACE_DIR)/overlays/$(strip $(2))/references" ]; then \
			cp -Rf "$(WORKSPACE_DIR)/overlays/$(strip $(2))/references/." "$$_ws_tmp/references/"; \
		fi; \
		_zip_items="SKILL.md references scripts"; \
		if [ -n "$(strip $(3))" ] && [ -d "$(BUNDLED_PLUGINS)" ]; then \
			for _plugin_dir in "$(BUNDLED_PLUGINS)"/*/; do \
				[ -d "$$_plugin_dir" ] || continue; \
				_pname=$$(basename "$$_plugin_dir"); \
				[ -f "$$_plugin_dir/plugin.json" ] || continue; \
				_dest="$$_ws_tmp/plugins/$$_pname"; \
				mkdir -p "$$_dest/bin"; \
				cp "$$_plugin_dir/plugin.json" "$$_dest/"; \
				[ -f "$$_plugin_dir/overlay.json" ] && cp "$$_plugin_dir/overlay.json" "$$_dest/"; \
				[ -d "$$_plugin_dir/skills" ] && cp -r "$$_plugin_dir/skills" "$$_dest/"; \
				[ -d "$$_plugin_dir/references" ] && cp -r "$$_plugin_dir/references" "$$_dest/"; \
				case "$(strip $(3))" in \
					darwin) \
						for _darwinbin in "$$_plugin_dir/bin/"*-darwin; do \
							[ -f "$$_darwinbin" ] || continue; \
							_base=$$(basename "$$_darwinbin" | sed 's/-darwin$$//'); \
							cp "$$_darwinbin" "$$_dest/bin/$$_base"; \
						done; \
						chmod +x "$$_dest/bin/"* 2>/dev/null; \
						;; \
					windows) \
						for _winbin in "$$_plugin_dir/bin/"*-windows-amd64.exe; do \
							[ -f "$$_winbin" ] || continue; \
							_base=$$(basename "$$_winbin" | sed 's/-windows-amd64\.exe$$/.exe/'); \
							cp "$$_winbin" "$$_dest/bin/$$_base"; \
						done; \
						;; \
				esac; \
			done; \
			[ -d "$$_ws_tmp/plugins" ] && _zip_items="$$_zip_items plugins"; \
		fi; \
		cd "$$_ws_tmp" && zip -qr "$(1)" $$_zip_items -x '*/__MACOSX/*' '*/.DS_Store' 2>/dev/null || true; \
		rm -rf "$$_ws_tmp"; \
	else \
		echo "⚠️  $(WORKSPACE_DIR) not found, skipping workspace files"; \
	fi
endef

BUNDLE_VERSION  ?= $(VERSION)
BUNDLE_VARIANT  ?= dev
BUNDLE_OUTPUT   ?= $(TARGET)/dingtalk-workspace.zip

## bundle: Build dingtalk-workspace.zip as multi-skill bundle (manifest.json + skills/*.zip)
bundle:
	@mkdir -p $(TARGET)
	@chmod +x scripts/build-bundle.sh
	@./scripts/build-bundle.sh $(CURDIR)/$(BUNDLE_OUTPUT) $(BUNDLE_VARIANT) $(BUNDLE_VERSION)

## package: Build all platforms and create unified release zip in ../target/
package: build-all
	@echo "📦 Packaging $(PKG_NAME).zip ..."
	@mkdir -p $(TARGET)/$(PKG_NAME)
	@rm -rf $(TARGET)/$(PKG_NAME).zip

	@cp $(DIST)/* $(TARGET)/$(PKG_NAME)/

	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/$(PKG_NAME)/dingtalk-workspace.zip,dev,darwin)

	@cd $(TARGET) && zip -qr $(PKG_NAME).zip $(PKG_NAME) -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/$(PKG_NAME)
	@echo ""
	@echo "✅ $(TARGET)/$(PKG_NAME).zip created:"
	@ls -lh $(TARGET)/$(PKG_NAME).zip
	@echo ""
	@echo "📂 Contents:"
	@unzip -l $(TARGET)/$(PKG_NAME).zip | head -30

# ── Platform definitions ─────────────────────────────────────────────
PLATFORMS := darwin-arm64 darwin-amd64 linux-amd64 linux-arm64 windows-amd64 windows-arm64

## dist: Build, package all 6 platforms individually, sign macOS binaries
dist: build-all
	@echo "📦 Creating all platform packages..."

	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-darwin.zip,dev,darwin)
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-windows.zip,dev,windows)
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-linux.zip,dev)

	@mkdir -p $(TARGET)
	@for plat in $(PLATFORMS); do \
		case $$plat in \
			windows-*) bin="dws-$$plat.exe"; _ws="dingtalk-workspace-windows.zip" ;; \
			darwin-*)  bin="dws-$$plat";     _ws="dingtalk-workspace-darwin.zip" ;; \
			*)         bin="dws-$$plat";     _ws="dingtalk-workspace-linux.zip" ;; \
		esac; \
		pkg="dws_res_$$plat"; \
		mkdir -p $(TARGET)/$$pkg; \
		cp $(DIST)/$$bin $(TARGET)/$$pkg/; \
		if [ -f "$(TARGET)/$$_ws" ]; then \
			cp $(TARGET)/$$_ws $(TARGET)/$$pkg/dingtalk-workspace.zip; \
		fi; \
		cd $(TARGET) && zip -qr $$pkg.zip $$pkg -x '*/__MACOSX/*' '*/.DS_Store' && cd $(CURDIR); \
		rm -rf $(TARGET)/$$pkg; \
		echo "  ✅ $$pkg.zip"; \
	done

	@rm -f $(TARGET)/dingtalk-workspace-darwin.zip $(TARGET)/dingtalk-workspace-windows.zip $(TARGET)/dingtalk-workspace-linux.zip

	@echo ""
	@echo "📦 All platform packages:"
	@ls -lh $(TARGET)/dws_res_*.zip

	@if [ -f "$(ENTITLEMENTS)" ]; then \
		echo ""; \
		for mac_pkg in darwin-arm64 darwin-amd64; do \
			pkg_zip="$(TARGET)/dws_res_$$mac_pkg.zip"; \
			signed_zip="$(TARGET)/dws_res_$${mac_pkg}_signed.zip"; \
			echo "🔏 Signing $$pkg_zip via $(SIGN_SERVER)..."; \
			curl -X POST $(SIGN_SERVER)/sign \
				-F "file=@$$pkg_zip" \
				-F "entitlements=@$(ENTITLEMENTS)" \
				-o "$$signed_zip" \
				--fail --silent --show-error && \
			mv "$$signed_zip" "$$pkg_zip" && \
			echo "  ✅ $$pkg_zip signed"; \
		done; \
	else \
		echo ""; \
		echo "⚠️  $(ENTITLEMENTS) not found, skipping macOS signing"; \
	fi

	@echo ""
	@echo "🎉 Done! v$(VERSION) | commit $(GIT_COMMIT) | $(BUILD_TIME)"
	@echo ""
	@echo "📦 Final packages:"
	@ls -lh $(TARGET)/dws_res_*.zip

# ── Upstream sync ─────────────────────────────────────────────────────
CLI_DIR := ../dingtalk-workspace-cli
CLI_UPSTREAM_REMOTE := upstream
CLI_UPSTREAM_URL := https://github.com/DingTalk-Real-AI/dingtalk-workspace-cli.git
CLI_UPSTREAM_TAG ?= v1.0.29
CLI_RELEASE_BRANCH ?= release/$(CLI_UPSTREAM_TAG)

## sync-upstream: Recreate a fresh release branch from upstream tag in CLI repo
sync-upstream:
	@echo "🔄 Syncing $(CLI_DIR): recreate $(CLI_RELEASE_BRANCH) from $(CLI_UPSTREAM_TAG)..."
	@cd $(CLI_DIR) && \
		if ! git remote get-url $(CLI_UPSTREAM_REMOTE) >/dev/null 2>&1; then \
			echo "  Adding remote $(CLI_UPSTREAM_REMOTE) → $(CLI_UPSTREAM_URL)"; \
			git remote add $(CLI_UPSTREAM_REMOTE) $(CLI_UPSTREAM_URL); \
		fi && \
		git fetch $(CLI_UPSTREAM_REMOTE) --tags && \
		if ! git diff-index --quiet HEAD --; then \
			echo "⚠️  Discarding uncommitted tracked changes (remote tag wins):"; \
			git status --short --untracked-files=no; \
			git reset --hard HEAD; \
		fi && \
		if [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
			echo "🧹 Removing untracked files left over from previous branches:"; \
			git ls-files --others --exclude-standard; \
			git clean -fd -e /build -e /dist -e /.idea; \
		fi && \
		if [ "$$(git symbolic-ref --short -q HEAD)" = "$(CLI_RELEASE_BRANCH)" ]; then \
			echo "  Currently on $(CLI_RELEASE_BRANCH), detaching before recreate"; \
			git checkout --quiet --detach; \
		fi && \
		if git show-ref --verify --quiet refs/heads/$(CLI_RELEASE_BRANCH); then \
			echo "  Deleting existing $(CLI_RELEASE_BRANCH)"; \
			git branch -D $(CLI_RELEASE_BRANCH); \
		fi && \
		echo "  Creating $(CLI_RELEASE_BRANCH) from $(CLI_UPSTREAM_TAG)" && \
		git checkout -b $(CLI_RELEASE_BRANCH) refs/tags/$(CLI_UPSTREAM_TAG) && \
		echo "✅ $(CLI_DIR) on $(CLI_RELEASE_BRANCH) @ $(CLI_UPSTREAM_TAG)"

## real-platform: REAL 打包使用 (win+mac only, with signing)
SIGN_INPUT = $(TARGET)/dws_res_mac.zip
SIGN_OUTPUT = $(TARGET)/dws_res_mac_signed.zip

real-platform: sync-upstream
	@$(MAKE) build-all BUILD_MODE=real
	@echo "📦 Creating platform packages..."
	@mkdir -p $(TARGET)/dws_res_win $(TARGET)/dws_res_mac
	@rm -rf $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@# ── Mac workspace zip (with universal plugin binaries) ──
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-mac.zip,real,darwin)

	@# ── Win workspace zip (with Windows plugin binaries) ──
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-win.zip,real,windows)

	@cp $(DIST)/dws-windows-amd64.exe $(TARGET)/dws_res_win/
	@if [ -f "$(TARGET)/dingtalk-workspace-win.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-win.zip $(TARGET)/dws_res_win/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_win.zip dws_res_win -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_win

	@cp $(DIST)/dws-darwin-amd64 $(TARGET)/dws_res_mac/
	@cp $(DIST)/dws-darwin-arm64 $(TARGET)/dws_res_mac/
	@if [ -f "$(TARGET)/dingtalk-workspace-mac.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dws_res_mac/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_mac.zip dws_res_mac -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_mac

	@rm -f $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dingtalk-workspace-win.zip

	@echo ""
	@echo "✅ Platform packages created:"
	@ls -lh $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@if [ -f "$(ENTITLEMENTS)" ]; then \
		echo ""; \
		echo "🔏 Signing macOS binary via $(SIGN_SERVER)..."; \
		curl -X POST $(SIGN_SERVER)/sign \
			-F "file=@$(SIGN_INPUT)" \
			-F "entitlements=@$(ENTITLEMENTS)" \
			-o $(SIGN_OUTPUT) \
			--fail --silent --show-error && \
		mv $(SIGN_OUTPUT) $(SIGN_INPUT) && \
		echo "✅ macOS package signed successfully" && \
		ls -lh $(SIGN_INPUT); \
	else \
		echo ""; \
		echo "⚠️  $(ENTITLEMENTS) not found, skipping macOS signing"; \
	fi

	@echo ""
	@echo "📂 Windows package contents:"
	@unzip -l $(TARGET)/dws_res_win.zip
	@echo ""
	@echo "📂 macOS package contents (signed):"
	@unzip -l $(TARGET)/dws_res_mac.zip

## local-platform: 本地 real 部署 (嵌入悟空运行时；不 sync-upstream，使用本地 CLI 分支)
## 必须出 real 产物：dev 二进制走 keychain，在悟空沙箱里读不到 DEK。
local-platform:
	@$(MAKE) build-all BUILD_MODE=real
	@echo "📦 Creating platform packages (local real, no sync-upstream)..."
	@mkdir -p $(TARGET)/dws_res_win $(TARGET)/dws_res_mac
	@rm -rf $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@# ── Mac workspace zip (with macOS plugin binaries) ──
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-mac.zip,real,darwin)

	@# ── Win workspace zip (with Windows plugin binaries) ──
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-win.zip,real,windows)

	@cp $(DIST)/dws-windows-amd64.exe $(TARGET)/dws_res_win/
	@if [ -f "$(TARGET)/dingtalk-workspace-win.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-win.zip $(TARGET)/dws_res_win/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_win.zip dws_res_win -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_win
	@cp $(DIST)/dws-darwin-amd64 $(TARGET)/dws_res_mac/
	@cp $(DIST)/dws-darwin-arm64 $(TARGET)/dws_res_mac/
	@if [ -f "$(TARGET)/dingtalk-workspace-mac.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dws_res_mac/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_mac.zip dws_res_mac -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_mac
	@rm -f $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dingtalk-workspace-win.zip
	@echo ""
	@echo "✅ Local platform packages created (no signing):"
	@ls -lh $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

## bundle-platform: REAL 打包 (bundle 版: dingtalk-workspace.zip 为多 skill bundle)
bundle-platform: sync-upstream
	@$(MAKE) build-all BUILD_MODE=real
	@echo "📦 Creating platform packages (bundle)..."
	@mkdir -p $(TARGET)/dws_res_win $(TARGET)/dws_res_mac
	@rm -rf $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@chmod +x scripts/build-bundle.sh
	@./scripts/build-bundle.sh $(CURDIR)/$(TARGET)/dingtalk-workspace.zip real $(VERSION)

	@# ── Platform workspace zips (with bundled plugins) ──
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-mac.zip,real,darwin)
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-win.zip,real,windows)

	@cp $(DIST)/dws-windows-amd64.exe $(TARGET)/dws_res_win/
	@if [ -f "$(TARGET)/dingtalk-workspace-win.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-win.zip $(TARGET)/dws_res_win/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_win.zip dws_res_win -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_win

	@cp $(DIST)/dws-darwin-amd64 $(TARGET)/dws_res_mac/
	@cp $(DIST)/dws-darwin-arm64 $(TARGET)/dws_res_mac/
	@if [ -f "$(TARGET)/dingtalk-workspace-mac.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dws_res_mac/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_mac.zip dws_res_mac -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_mac

	@rm -f $(TARGET)/dingtalk-workspace.zip $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dingtalk-workspace-win.zip

	@echo ""
	@echo "✅ Bundle platform packages created:"
	@ls -lh $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@if [ -f "$(ENTITLEMENTS)" ]; then \
		echo ""; \
		echo "🔏 Signing macOS binary via $(SIGN_SERVER)..."; \
		curl -X POST $(SIGN_SERVER)/sign \
			-F "file=@$(SIGN_INPUT)" \
			-F "entitlements=@$(ENTITLEMENTS)" \
			-o $(SIGN_OUTPUT) \
			--fail --silent --show-error && \
		mv $(SIGN_OUTPUT) $(SIGN_INPUT) && \
		echo "✅ macOS package signed successfully" && \
		ls -lh $(SIGN_INPUT); \
	else \
		echo ""; \
		echo "⚠️  $(ENTITLEMENTS) not found, skipping macOS signing"; \
	fi

	@echo ""
	@echo "📂 Windows package contents:"
	@unzip -l $(TARGET)/dws_res_win.zip
	@echo ""
	@echo "📂 macOS package contents (signed):"
	@unzip -l $(TARGET)/dws_res_mac.zip

## sign: Sign macOS binary standalone (if real-platform already ran)
sign:
	@if [ ! -f "$(SIGN_INPUT)" ]; then \
		echo "❌ $(SIGN_INPUT) not found. Run 'make real-platform' first."; \
		exit 1; \
	fi
	@if [ ! -f "$(ENTITLEMENTS)" ]; then \
		echo "❌ $(ENTITLEMENTS) not found."; \
		exit 1; \
	fi
	@echo "🔏 Signing macOS binary via $(SIGN_SERVER)..."
	@curl -X POST $(SIGN_SERVER)/sign \
		-F "file=@$(SIGN_INPUT)" \
		-F "entitlements=@$(ENTITLEMENTS)" \
		-o $(SIGN_OUTPUT) \
		--fail --silent --show-error
	@mv $(SIGN_OUTPUT) $(SIGN_INPUT)
	@echo "✅ Signed: $(SIGN_INPUT)"
	@ls -lh $(SIGN_INPUT)

## public: Build and create public release directory with dws/ and skill/ subdirs
public: build-all
	@echo "📦 Creating public release..."
	@rm -rf $(TARGET)/public
	@mkdir -p $(TARGET)/public/dws $(TARGET)/public/skill

	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/public/skill/dingtalk-workspace.zip,dev)
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-mac.zip,dev,darwin)
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-win.zip,dev,windows)

	@mkdir -p $(TARGET)/dws_res_win
	@cp $(DIST)/dws-windows-amd64.exe $(TARGET)/dws_res_win/
	@if [ -f "$(TARGET)/dingtalk-workspace-win.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-win.zip $(TARGET)/dws_res_win/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr public/dws/dws_res_win.zip dws_res_win -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_win

	@mkdir -p $(TARGET)/dws_res_mac
	@cp $(DIST)/dws-darwin-amd64 $(TARGET)/dws_res_mac/
	@cp $(DIST)/dws-darwin-arm64 $(TARGET)/dws_res_mac/
	@if [ -f "$(TARGET)/dingtalk-workspace-mac.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dws_res_mac/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr public/dws/dws_res_mac.zip dws_res_mac -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_mac

	@rm -f $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dingtalk-workspace-win.zip

	@cp README.md $(TARGET)/public/

	@echo ""
	@echo "✅ Public release created:"
	@find $(TARGET)/public -type f | sort
	@echo ""
	@ls -lh $(TARGET)/public/dws/ $(TARGET)/public/skill/

	@PUBLIC_ZIP="dws-v$(VERSION)-$$(date '+%Y%m%d%H%M%S').zip"; \
	cd $(TARGET) && zip -qr "$$PUBLIC_ZIP" public -x '*/__MACOSX/*' '*/.DS_Store' && \
	rm -rf public && \
	echo "" && \
	echo "📦 $$PUBLIC_ZIP:" && \
	ls -lh "$$PUBLIC_ZIP"

RELEASE_DIR := releases/current

## release: Build all, generate manifest, prepare git release
release: build-all
	@echo "📦 Preparing release v$(VERSION)..."
	@mkdir -p $(RELEASE_DIR)
	@cp $(DIST)/dws-darwin-amd64 $(RELEASE_DIR)/
	@cp $(DIST)/dws-darwin-arm64 $(RELEASE_DIR)/
	@cp $(DIST)/dws-linux-amd64 $(RELEASE_DIR)/
	@cp $(DIST)/dws-linux-arm64 $(RELEASE_DIR)/
	@cp $(DIST)/dws-windows-amd64.exe $(RELEASE_DIR)/
	@cp $(DIST)/dws-windows-arm64.exe $(RELEASE_DIR)/
	@$(call build-workspace-zip,$(CURDIR)/$(RELEASE_DIR)/skill-pack.zip,dev)
	@printf '{\n  "version": "%s",\n  "date": "%s",\n  "changelog": "",\n  "assets": {\n    "darwin-amd64": "releases/current/dws-darwin-amd64",\n    "darwin-arm64": "releases/current/dws-darwin-arm64",\n    "linux-amd64": "releases/current/dws-linux-amd64",\n    "linux-arm64": "releases/current/dws-linux-arm64",\n    "windows-amd64": "releases/current/dws-windows-amd64.exe",\n    "windows-arm64": "releases/current/dws-windows-arm64.exe"\n  },\n  "skillPack": "releases/current/skill-pack.zip"\n}\n' \
		"$(VERSION)" "$$(date '+%Y-%m-%d')" > releases/latest.json
	@echo ""
	@echo "✅ Release v$(VERSION) prepared:"
	@ls -lh $(RELEASE_DIR)/
	@echo ""
	@echo "📌 发布步骤:"
	@echo "  1. 编辑 releases/latest.json 的 changelog 字段（可选）"
	@echo "  2. git add releases/ && git commit -m 'Release v$(VERSION)'"
	@echo "  3. git tag v$(VERSION)"
	@echo "  4. git push && git push --tags"

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //; s/: /\t/' | column -t -s $$'\t'

# ════════════════════════════════════════════════════════════════════════════
# CI / Jenkins 自动化专用块
# ════════════════════════════════════════════════════════════════════════════
# 设计原则：
#   1. 与 real-platform/sync-upstream 完全解耦，互不影响
#   2. 假设 caller（Jenkins shell）已经把 CLI_DIR 准备好（例如通过 codeload tarball）
#   3. 不做任何 git 操作（GitHub git 协议在阿里内网不稳定）
#   4. 默认 BUILD_MODE=real + 签名（如有 entitlements.plist）
#
# 使用：
#   make ci-platform VERSION=0.2.59 CI_CLI_TARBALL_URL=https://codeload.github.com/.../tar.gz/refs/tags/v1.0.19
#   或者 caller 自己准备好 CLI_DIR 后：
#   make ci-platform VERSION=0.2.59
# ════════════════════════════════════════════════════════════════════════════

CI_CLI_TARBALL_URL  ?= https://codeload.github.com/DingTalk-Real-AI/dingtalk-workspace-cli/tar.gz/refs/heads/main
CI_BUILD_MODE       ?= real
CI_WORKSPACE_VARIANT ?= real

## ci-prepare-cli: 通过 codeload tarball 把 CLI_DIR 拉到位（不走 git 协议）
ci-prepare-cli:
	@if [ -d "$(CLI_DIR)" ] && [ -n "$$(ls -A $(CLI_DIR) 2>/dev/null)" ]; then \
		echo "✅ $(CLI_DIR) already populated, skip download"; \
	else \
		echo "==> Downloading CLI tarball from $(CI_CLI_TARBALL_URL)"; \
		_tmp=$$(mktemp -d); \
		curl -fL --connect-timeout 30 --max-time 600 \
			-o "$$_tmp/cli.tar.gz" "$(CI_CLI_TARBALL_URL)" || { echo "❌ tarball download failed"; rm -rf "$$_tmp"; exit 1; }; \
		tar -xzf "$$_tmp/cli.tar.gz" -C "$$_tmp"; \
		_extracted=$$(find "$$_tmp" -maxdepth 1 -mindepth 1 -type d | head -1); \
		[ -n "$$_extracted" ] || { echo "❌ no directory found in tarball"; ls -la "$$_tmp"; rm -rf "$$_tmp"; exit 1; }; \
		mkdir -p "$$(dirname $(CLI_DIR))"; \
		rm -rf "$(CLI_DIR)"; \
		mv "$$_extracted" "$(CLI_DIR)"; \
		rm -rf "$$_tmp"; \
		echo "✅ $(CLI_DIR) ready ($$(du -sh $(CLI_DIR) | cut -f1))"; \
	fi

## ci-platform: CI 专用打包（不依赖 sync-upstream，BUILD_MODE=real，含签名，含平台插件）
ci-platform: ci-prepare-cli
	@$(MAKE) build-all BUILD_MODE=$(CI_BUILD_MODE)
	@echo "📦 [CI] Creating platform packages..."
	@mkdir -p $(TARGET)/dws_res_win $(TARGET)/dws_res_mac
	@rm -rf $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@# ── Mac workspace zip (with macOS plugin binaries) ──
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-mac.zip,$(CI_WORKSPACE_VARIANT),darwin)

	@# ── Win workspace zip (with Windows plugin binaries) ──
	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace-win.zip,$(CI_WORKSPACE_VARIANT),windows)

	@cp $(DIST)/dws-windows-amd64.exe $(TARGET)/dws_res_win/
	@if [ -f "$(TARGET)/dingtalk-workspace-win.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-win.zip $(TARGET)/dws_res_win/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_win.zip dws_res_win -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_win

	@cp $(DIST)/dws-darwin-amd64 $(TARGET)/dws_res_mac/
	@cp $(DIST)/dws-darwin-arm64 $(TARGET)/dws_res_mac/
	@if [ -f "$(TARGET)/dingtalk-workspace-mac.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dws_res_mac/dingtalk-workspace.zip; \
	fi
	@cd $(TARGET) && zip -qr dws_res_mac.zip dws_res_mac -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_mac

	@rm -f $(TARGET)/dingtalk-workspace-mac.zip $(TARGET)/dingtalk-workspace-win.zip

	@echo ""
	@echo "✅ [CI] Platform packages created:"
	@ls -lh $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@if [ -f "$(ENTITLEMENTS)" ]; then \
		echo ""; \
		echo "🔏 [CI] Signing macOS binary via $(SIGN_SERVER)..."; \
		curl -X POST $(SIGN_SERVER)/sign \
			-F "file=@$(SIGN_INPUT)" \
			-F "entitlements=@$(ENTITLEMENTS)" \
			-o $(SIGN_OUTPUT) \
			--fail --silent --show-error && \
		mv $(SIGN_OUTPUT) $(SIGN_INPUT) && \
		echo "✅ macOS package signed successfully" && \
		ls -lh $(SIGN_INPUT); \
	else \
		echo ""; \
		echo "⚠️  $(ENTITLEMENTS) not found, skipping macOS signing"; \
	fi

	@echo ""
	@echo "📂 Windows package contents:"
	@unzip -l $(TARGET)/dws_res_win.zip
	@echo ""
	@echo "📂 macOS package contents:"
	@unzip -l $(TARGET)/dws_res_mac.zip

.PHONY: ci-prepare-cli ci-platform
