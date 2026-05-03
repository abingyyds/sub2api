package routes

import (
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/httpclient"

	"github.com/gin-gonic/gin"
)

type downloadAsset struct {
	FileName    string
	UpstreamURL string
	ContentType string
}

var tutorialDownloadAssets = map[string]downloadAsset{
	"nodejs-windows-x64-msi": {
		FileName:    "node-v24.15.0-x64.msi",
		UpstreamURL: "https://nodejs.org/dist/v24.15.0/node-v24.15.0-x64.msi",
		ContentType: "application/x-msi",
	},
	"nodejs-windows-arm64-zip": {
		FileName:    "node-v24.15.0-win-arm64.zip",
		UpstreamURL: "https://nodejs.org/dist/v24.15.0/node-v24.15.0-win-arm64.zip",
		ContentType: "application/zip",
	},
	"nodejs-macos-intel-pkg": {
		FileName:    "node-v24.15.0.pkg",
		UpstreamURL: "https://nodejs.org/dist/v24.15.0/node-v24.15.0.pkg",
		ContentType: "application/octet-stream",
	},
	"nodejs-macos-arm64-tar-gz": {
		FileName:    "node-v24.15.0-darwin-arm64.tar.gz",
		UpstreamURL: "https://nodejs.org/dist/v24.15.0/node-v24.15.0-darwin-arm64.tar.gz",
		ContentType: "application/gzip",
	},
	"nodejs-linux-x64-tar-xz": {
		FileName:    "node-v24.15.0-linux-x64.tar.xz",
		UpstreamURL: "https://nodejs.org/dist/v24.15.0/node-v24.15.0-linux-x64.tar.xz",
		ContentType: "application/x-xz",
	},
	"nodejs-linux-arm64-tar-xz": {
		FileName:    "node-v24.15.0-linux-arm64.tar.xz",
		UpstreamURL: "https://nodejs.org/dist/v24.15.0/node-v24.15.0-linux-arm64.tar.xz",
		ContentType: "application/x-xz",
	},
	"cc-switch-windows-x64-msi": {
		FileName:    "CC-Switch-v3.13.0-Windows.msi",
		UpstreamURL: "https://github.com/farion1231/cc-switch/releases/download/v3.13.0/CC-Switch-v3.13.0-Windows.msi",
		ContentType: "application/octet-stream",
	},
	"cc-switch-windows-x64-portable-zip": {
		FileName:    "CC-Switch-v3.13.0-Windows-Portable.zip",
		UpstreamURL: "https://github.com/farion1231/cc-switch/releases/download/v3.13.0/CC-Switch-v3.13.0-Windows-Portable.zip",
		ContentType: "application/zip",
	},
	"cc-switch-macos-universal-dmg": {
		FileName:    "CC-Switch-v3.13.0-macOS.dmg",
		UpstreamURL: "https://github.com/farion1231/cc-switch/releases/download/v3.13.0/CC-Switch-v3.13.0-macOS.dmg",
		ContentType: "application/x-apple-diskimage",
	},
	"cc-switch-linux-x64-deb": {
		FileName:    "CC-Switch-v3.13.0-Linux-x86_64.deb",
		UpstreamURL: "https://github.com/farion1231/cc-switch/releases/download/v3.13.0/CC-Switch-v3.13.0-Linux-x86_64.deb",
		ContentType: "application/vnd.debian.binary-package",
	},
	"cc-switch-linux-arm64-deb": {
		FileName:    "CC-Switch-v3.13.0-Linux-arm64.deb",
		UpstreamURL: "https://github.com/farion1231/cc-switch/releases/download/v3.13.0/CC-Switch-v3.13.0-Linux-arm64.deb",
		ContentType: "application/vnd.debian.binary-package",
	},
	"cc-switch-linux-x64-appimage": {
		FileName:    "CC-Switch-v3.13.0-Linux-x86_64.AppImage",
		UpstreamURL: "https://github.com/farion1231/cc-switch/releases/download/v3.13.0/CC-Switch-v3.13.0-Linux-x86_64.AppImage",
		ContentType: "application/octet-stream",
	},
	"cc-switch-linux-arm64-appimage": {
		FileName:    "CC-Switch-v3.13.0-Linux-arm64.AppImage",
		UpstreamURL: "https://github.com/farion1231/cc-switch/releases/download/v3.13.0/CC-Switch-v3.13.0-Linux-arm64.AppImage",
		ContentType: "application/octet-stream",
	},
	"codex-windows-x64-exe": {
		FileName:    "codex-x86_64-pc-windows-msvc.exe",
		UpstreamURL: "https://github.com/openai/codex/releases/download/rust-v0.128.0/codex-x86_64-pc-windows-msvc.exe",
		ContentType: "application/octet-stream",
	},
	"codex-windows-arm64-exe": {
		FileName:    "codex-aarch64-pc-windows-msvc.exe",
		UpstreamURL: "https://github.com/openai/codex/releases/download/rust-v0.128.0/codex-aarch64-pc-windows-msvc.exe",
		ContentType: "application/octet-stream",
	},
	"codex-macos-arm64-dmg": {
		FileName:    "codex-aarch64-apple-darwin.dmg",
		UpstreamURL: "https://github.com/openai/codex/releases/download/rust-v0.128.0/codex-aarch64-apple-darwin.dmg",
		ContentType: "application/x-apple-diskimage",
	},
	"codex-macos-intel-dmg": {
		FileName:    "codex-x86_64-apple-darwin.dmg",
		UpstreamURL: "https://github.com/openai/codex/releases/download/rust-v0.128.0/codex-x86_64-apple-darwin.dmg",
		ContentType: "application/x-apple-diskimage",
	},
	"codex-linux-x64-tar-gz": {
		FileName:    "codex-x86_64-unknown-linux-gnu.tar.gz",
		UpstreamURL: "https://github.com/openai/codex/releases/download/rust-v0.122.0/codex-x86_64-unknown-linux-gnu.tar.gz",
		ContentType: "application/gzip",
	},
	"codex-linux-arm64-tar-gz": {
		FileName:    "codex-aarch64-unknown-linux-gnu.tar.gz",
		UpstreamURL: "https://github.com/openai/codex/releases/download/rust-v0.122.0/codex-aarch64-unknown-linux-gnu.tar.gz",
		ContentType: "application/gzip",
	},
	"cherry-studio-windows-x64-setup-exe": {
		FileName:    "Cherry-Studio-1.9.2-x64-setup.exe",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-x64-setup.exe",
		ContentType: "application/octet-stream",
	},
	"cherry-studio-windows-x64-portable-exe": {
		FileName:    "Cherry-Studio-1.9.2-x64-portable.exe",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-x64-portable.exe",
		ContentType: "application/octet-stream",
	},
	"cherry-studio-windows-arm64-setup-exe": {
		FileName:    "Cherry-Studio-1.9.2-arm64-setup.exe",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-arm64-setup.exe",
		ContentType: "application/octet-stream",
	},
	"cherry-studio-macos-intel-dmg": {
		FileName:    "Cherry-Studio-1.9.2-x64.dmg",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-x64.dmg",
		ContentType: "application/x-apple-diskimage",
	},
	"cherry-studio-macos-arm64-dmg": {
		FileName:    "Cherry-Studio-1.9.2-arm64.dmg",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-arm64.dmg",
		ContentType: "application/x-apple-diskimage",
	},
	"cherry-studio-linux-x64-appimage": {
		FileName:    "Cherry-Studio-1.9.2-x86_64.AppImage",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-x86_64.AppImage",
		ContentType: "application/octet-stream",
	},
	"cherry-studio-linux-x64-deb": {
		FileName:    "Cherry-Studio-1.9.2-amd64.deb",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-amd64.deb",
		ContentType: "application/vnd.debian.binary-package",
	},
	"cherry-studio-linux-arm64-appimage": {
		FileName:    "Cherry-Studio-1.9.2-arm64.AppImage",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-arm64.AppImage",
		ContentType: "application/octet-stream",
	},
	"cherry-studio-linux-arm64-deb": {
		FileName:    "Cherry-Studio-1.9.2-arm64.deb",
		UpstreamURL: "https://github.com/CherryHQ/cherry-studio/releases/download/v1.9.2/Cherry-Studio-1.9.2-arm64.deb",
		ContentType: "application/vnd.debian.binary-package",
	},
}

// RegisterCommonRoutes 注册通用路由（健康检查、状态等）
func RegisterCommonRoutes(r *gin.Engine, cfg *config.Config) {
	downloadClient, err := httpclient.GetClient(httpclient.Options{
		Timeout:  30 * time.Minute,
		ProxyURL: cfg.Update.ProxyURL,
	})
	if err != nil {
		log.Printf("[Downloads] failed to initialize shared HTTP client: %v", err)
		downloadClient = &http.Client{Timeout: 30 * time.Minute}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	serveDownload := func(c *gin.Context) {
		asset, ok := tutorialDownloadAssets[c.Param("id")]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "download asset not found"})
			return
		}

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, asset.UpstreamURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upstream request"})
			return
		}
		req.Header.Set("User-Agent", "Sub2API-Downloads")
		if rangeHeader := c.GetHeader("Range"); rangeHeader != "" {
			req.Header.Set("Range", rangeHeader)
		}
		if ifRange := c.GetHeader("If-Range"); ifRange != "" {
			req.Header.Set("If-Range", ifRange)
		}

		resp, err := downloadClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch upstream download"})
			return
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":         "upstream download failed",
				"upstream_code": resp.StatusCode,
			})
			return
		}

		contentType := asset.ContentType
		if upstreamType := resp.Header.Get("Content-Type"); upstreamType != "" {
			contentType = upstreamType
		}

		headers := map[string]string{
			"Content-Disposition": contentDispositionHeader(asset.FileName),
			"Cache-Control":       "private, no-store, max-age=0",
			"X-Accel-Buffering":   "no",
		}
		for _, key := range []string{"Accept-Ranges", "Content-Length", "Content-Range", "ETag", "Last-Modified"} {
			if value := resp.Header.Get(key); value != "" {
				headers[key] = value
			}
		}

		c.DataFromReader(resp.StatusCode, resp.ContentLength, contentType, resp.Body, headers)
	}

	r.GET("/downloads/:id", serveDownload)
	r.HEAD("/downloads/:id", serveDownload)

	// Claude Code 遥测日志（忽略，直接返回200）
	eventLoggingHandler := func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
	r.POST("/api/event_logging/batch", eventLoggingHandler)
	r.GET("/api/event_logging/batch", eventLoggingHandler)

	// Setup status endpoint (always returns needs_setup: false in normal mode)
	// This is used by the frontend to detect when the service has restarted after setup
	r.GET("/setup/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"needs_setup": false,
				"step":        "completed",
			},
		})
	})
}

func contentDispositionHeader(fileName string) string {
	baseName := path.Base(fileName)
	return `attachment; filename="` + baseName + `"; filename*=UTF-8''` + url.PathEscape(baseName)
}
