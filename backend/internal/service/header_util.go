package service

import (
	"net/http"
	"strings"
)

// headerWireCasing 定义常见 Claude CLI 请求头的真实大小写，尽量贴近官方客户端的 wire format。
var headerWireCasing = map[string]string{
	"accept":                                    "Accept",
	"user-agent":                                "User-Agent",
	"x-stainless-retry-count":                   "X-Stainless-Retry-Count",
	"x-stainless-timeout":                       "X-Stainless-Timeout",
	"x-stainless-lang":                          "X-Stainless-Lang",
	"x-stainless-package-version":               "X-Stainless-Package-Version",
	"x-stainless-os":                            "X-Stainless-OS",
	"x-stainless-arch":                          "X-Stainless-Arch",
	"x-stainless-runtime":                       "X-Stainless-Runtime",
	"x-stainless-runtime-version":               "X-Stainless-Runtime-Version",
	"x-stainless-helper-method":                 "x-stainless-helper-method",
	"anthropic-dangerous-direct-browser-access": "anthropic-dangerous-direct-browser-access",
	"anthropic-version":                         "anthropic-version",
	"anthropic-beta":                            "anthropic-beta",
	"x-app":                                     "x-app",
	"content-type":                              "content-type",
	"accept-language":                           "accept-language",
	"sec-fetch-mode":                            "sec-fetch-mode",
	"accept-encoding":                           "accept-encoding",
	"authorization":                             "authorization",
	"x-claude-code-session-id":                  "X-Claude-Code-Session-Id",
	"x-client-request-id":                       "x-client-request-id",
	"content-length":                            "content-length",
	"x-api-key":                                 "x-api-key",
}

func resolveWireCasing(key string) string {
	if wk, ok := headerWireCasing[strings.ToLower(key)]; ok {
		return wk
	}
	return key
}

func setHeaderRaw(h http.Header, key, value string) {
	h.Del(key)
	if wk := resolveWireCasing(key); wk != key {
		delete(h, wk)
	}
	delete(h, key)
	h[key] = []string{value}
}

func addHeaderRaw(h http.Header, key, value string) {
	h[key] = append(h[key], value)
}

func getHeaderRaw(h http.Header, key string) string {
	if vals := h[key]; len(vals) > 0 {
		return vals[0]
	}
	if wk := resolveWireCasing(key); wk != key {
		if vals := h[wk]; len(vals) > 0 {
			return vals[0]
		}
	}
	return h.Get(key)
}
