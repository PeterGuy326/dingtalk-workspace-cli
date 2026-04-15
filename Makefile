.PHONY: build build-all test clean npm-pack package dist real-platform sign public release help

VERSION ?= 0.2.44
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

## clean: Remove build artifacts
clean:
	rm -f dws dws-darwin-* dws-linux-* dws-windows-*
	rm -rf $(DIST)
	rm -rf npm/bin

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

# $(1) = output zip path, $(2) = workspace variant: dev | real (real applies overlays/real)
define build-workspace-zip
	@if [ -d "$(WORKSPACE_DIR)" ]; then \
		_ws_tmp=$$(mktemp -d); \
		cp -r "$(WORKSPACE_DIR)/SKILL.md" "$(WORKSPACE_DIR)/references" "$(WORKSPACE_DIR)/scripts" "$$_ws_tmp/"; \
		if [ "$(strip $(2))" != "dev" ] && [ -d "$(WORKSPACE_DIR)/overlays/$(strip $(2))/references" ]; then \
			cp -Rf "$(WORKSPACE_DIR)/overlays/$(strip $(2))/references/." "$$_ws_tmp/references/"; \
		fi; \
		cd "$$_ws_tmp" && zip -qr "$(1)" SKILL.md references scripts -x '*/__MACOSX/*' '*/.DS_Store' 2>/dev/null || true; \
		rm -rf "$$_ws_tmp"; \
	else \
		echo "⚠️  $(WORKSPACE_DIR) not found, skipping workspace files"; \
	fi
endef

## package: Build all platforms and create unified release zip in ../target/
package: build-all
	@echo "📦 Packaging $(PKG_NAME).zip ..."
	@mkdir -p $(TARGET)/$(PKG_NAME)
	@rm -rf $(TARGET)/$(PKG_NAME).zip

	@cp $(DIST)/* $(TARGET)/$(PKG_NAME)/

	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/$(PKG_NAME)/dingtalk-workspace.zip,dev)

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

	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace.zip,dev)

	@mkdir -p $(TARGET)
	@for plat in $(PLATFORMS); do \
		case $$plat in \
			windows-*) bin="dws-$$plat.exe" ;; \
			*)         bin="dws-$$plat" ;; \
		esac; \
		pkg="dws_res_$$plat"; \
		mkdir -p $(TARGET)/$$pkg; \
		cp $(DIST)/$$bin $(TARGET)/$$pkg/; \
		if [ -f "$(TARGET)/dingtalk-workspace.zip" ]; then \
			cp $(TARGET)/dingtalk-workspace.zip $(TARGET)/$$pkg/; \
		fi; \
		cd $(TARGET) && zip -qr $$pkg.zip $$pkg -x '*/__MACOSX/*' '*/.DS_Store' && cd $(CURDIR); \
		rm -rf $(TARGET)/$$pkg; \
		echo "  ✅ $$pkg.zip"; \
	done

	@rm -f $(TARGET)/dingtalk-workspace.zip

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

## real-platform: REAL 打包使用 (win+mac only, with signing)
SIGN_INPUT = $(TARGET)/dws_res_mac.zip
SIGN_OUTPUT = $(TARGET)/dws_res_mac_signed.zip

real-platform:
	@$(MAKE) build-all BUILD_MODE=real
	@echo "📦 Creating platform packages..."
	@mkdir -p $(TARGET)/dws_res_win $(TARGET)/dws_res_mac
	@rm -rf $(TARGET)/dws_res_win.zip $(TARGET)/dws_res_mac.zip

	@$(call build-workspace-zip,$(CURDIR)/$(TARGET)/dingtalk-workspace.zip,real)

	@cp $(DIST)/dws-windows-amd64.exe $(TARGET)/dws_res_win/
	@if [ -f "$(TARGET)/dingtalk-workspace.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace.zip $(TARGET)/dws_res_win/; \
	fi
	@cd $(TARGET) && zip -qr dws_res_win.zip dws_res_win -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_win

	@cp $(DIST)/dws-darwin-amd64 $(TARGET)/dws_res_mac/
	@cp $(DIST)/dws-darwin-arm64 $(TARGET)/dws_res_mac/
	@if [ -f "$(TARGET)/dingtalk-workspace.zip" ]; then \
		cp $(TARGET)/dingtalk-workspace.zip $(TARGET)/dws_res_mac/; \
	fi
	@cd $(TARGET) && zip -qr dws_res_mac.zip dws_res_mac -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_mac

	@rm -f $(TARGET)/dingtalk-workspace.zip

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

	@mkdir -p $(TARGET)/dws_res_win
	@cp $(DIST)/dws-windows-amd64.exe $(TARGET)/dws_res_win/
	@if [ -f "$(TARGET)/public/skill/dingtalk-workspace.zip" ]; then \
		cp $(TARGET)/public/skill/dingtalk-workspace.zip $(TARGET)/dws_res_win/; \
	fi
	@cd $(TARGET) && zip -qr public/dws/dws_res_win.zip dws_res_win -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_win

	@mkdir -p $(TARGET)/dws_res_mac
	@cp $(DIST)/dws-darwin-arm64 $(TARGET)/dws_res_mac/
	@if [ -f "$(TARGET)/public/skill/dingtalk-workspace.zip" ]; then \
		cp $(TARGET)/public/skill/dingtalk-workspace.zip $(TARGET)/dws_res_mac/; \
	fi
	@cd $(TARGET) && zip -qr public/dws/dws_res_mac.zip dws_res_mac -x '*/__MACOSX/*' '*/.DS_Store'
	@rm -rf $(TARGET)/dws_res_mac

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
