// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// mediaMaxDownloadBytes caps a single inbound image download (screenshots are
// well under this; the cap is a hostile-input guard).
const mediaMaxDownloadBytes = 20 << 20

// pictureDownloadCode digs the downloadCode out of a picture callback's
// loosely-typed content payload (the stream SDK models Content as
// interface{}).
func pictureDownloadCode(content interface{}) string {
	m, ok := content.(map[string]interface{})
	if !ok {
		return ""
	}
	for _, key := range []string{"downloadCode", "pictureDownloadCode"} {
		if v, ok := m[key].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// fileDownloadInfo extracts downloadCode and fileName from a file callback's
// content payload (msgtype="file"). Returns ("","") when not parseable.
func fileDownloadInfo(content interface{}) (downloadCode, fileName string) {
	m, ok := content.(map[string]interface{})
	if !ok {
		return "", ""
	}
	for _, key := range []string{"downloadCode", "fileDownloadCode"} {
		if v, ok := m[key].(string); ok && strings.TrimSpace(v) != "" {
			downloadCode = strings.TrimSpace(v)
			break
		}
	}
	if v, ok := m["fileName"].(string); ok {
		fileName = strings.TrimSpace(v)
	}
	if fileName == "" {
		fileName = "未知文件"
	}
	return downloadCode, fileName
}

// extractCallbackText pulls the visible text out of a structured-text callback
// payload (msgtype=richText / markdown / etc.) for the case where the SDK's
// data.Text.Content is empty. This matters because `dws chat message send
// --group ... --text ...` sends msgType="markdown" by default, and DingTalk's
// stream callback for markdown messages routinely leaves data.Text.Content
// blank while stashing the real body in data.Content (loosely-typed).
// Returns "" if no text-shaped field is found.
func extractCallbackText(content interface{}) string {
	switch v := content.(type) {
	case string:
		return strings.TrimSpace(v)
	case map[string]interface{}:
		// Common shapes: {"text":"..."}, {"title":"...","text":"..."},
		// {"content":"..."}, richText {"richText":[{"text":"..."}]}.
		for _, key := range []string{"text", "content", "markdown", "title"} {
			if s, ok := v[key].(string); ok && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
		if arr, ok := v["richText"].([]interface{}); ok {
			var b strings.Builder
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					if s, ok := m["text"].(string); ok {
						b.WriteString(s)
					}
				}
			}
			if s := strings.TrimSpace(b.String()); s != "" {
				return s
			}
		}
	}
	return ""
}

// downloadMessageFile resolves a chatbot media callback (picture etc.) to a
// local temp file via the robot messageFiles/download API. Error-screenshot
// questions are the top Q&A inbound; without this the connector silently
// drops every picture message.
func (c *aiCardClient) downloadMessageFile(ctx context.Context, robotCode, downloadCode string) (string, error) {
	raw, err := c.callRaw(ctx, http.MethodPost, "/v1.0/robot/messageFiles/download", map[string]any{
		"robotCode":    robotCode,
		"downloadCode": downloadCode,
	})
	if err != nil {
		return "", err
	}
	var parsed struct {
		DownloadUrl string `json:"downloadUrl"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil || strings.TrimSpace(parsed.DownloadUrl) == "" {
		return "", fmt.Errorf("messageFiles/download 未返回 downloadUrl: %s", truncateRunes(raw, 200))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.DownloadUrl, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("媒体下载 HTTP %d", resp.StatusCode)
	}
	dir := filepath.Join(os.TempDir(), "dws-connect-media")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dest := filepath.Join(dir, uuid.NewString()+mediaExt(parsed.DownloadUrl, resp.Header.Get("Content-Type")))
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, io.LimitReader(resp.Body, mediaMaxDownloadBytes)); err != nil {
		_ = os.Remove(dest)
		return "", err
	}
	return dest, nil
}

// mediaExt picks a file extension from the response content type, falling
// back to the URL path, then ".png" (DingTalk screenshots default to png).
func mediaExt(rawURL, contentType string) string {
	ct := strings.ToLower(contentType)
	switch {
	case strings.Contains(ct, "png"):
		return ".png"
	case strings.Contains(ct, "jpeg"), strings.Contains(ct, "jpg"):
		return ".jpg"
	case strings.Contains(ct, "gif"):
		return ".gif"
	case strings.Contains(ct, "webp"):
		return ".webp"
	}
	if u, err := url.Parse(rawURL); err == nil {
		if ext := strings.ToLower(filepath.Ext(u.Path)); len(ext) >= 2 && len(ext) <= 5 {
			return ext
		}
	}
	return ".png"
}
