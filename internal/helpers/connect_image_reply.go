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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const imageUploadMaxBytes = 500 * 1024

var (
	imagePathPattern   = regexp.MustCompile(`((?:/[^/\s"'` + "`" + `]+)+\.(?:png|jpg|jpeg|gif|webp|bmp))`)
	bashOpenCmdPattern = regexp.MustCompile("```bash\\s*\\n\\s*open\\s+[^\\n]*\\n\\s*```")
	viewHintPattern    = regexp.MustCompile(`(?i)(?:查看|打开|预览)(?:图片|文件|它)?[：:][^\n]*(?:\n|$)`)
)

func extractImagePaths(text string) []string {
	matches := imagePathPattern.FindAllStringSubmatch(text, -1)
	seen := make(map[string]bool)
	var paths []string
	for _, m := range matches {
		p := m[1]
		if seen[p] {
			continue
		}
		if _, err := os.Stat(p); err != nil {
			continue
		}
		seen[p] = true
		paths = append(paths, p)
	}
	return paths
}

func compressImageForUpload(imagePath string) ([]byte, string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, "", err
	}

	ext := strings.ToLower(filepath.Ext(imagePath))
	if (ext == ".jpg" || ext == ".jpeg" || ext == ".gif") && len(data) <= imageUploadMaxBytes {
		return data, filepath.Base(imagePath), nil
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		if len(data) <= imageUploadMaxBytes {
			return data, filepath.Base(imagePath), nil
		}
		return nil, "", fmt.Errorf("decode image: %w", err)
	}

	baseName := filepath.Base(imagePath)
	if ext != ".jpg" && ext != ".jpeg" {
		baseName = baseName[:len(baseName)-len(ext)] + ".jpg"
	}

	for quality := 80; quality >= 10; quality -= 15 {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, "", err
		}
		if buf.Len() <= imageUploadMaxBytes {
			return buf.Bytes(), baseName, nil
		}
	}

	return nil, "", fmt.Errorf("cannot compress image below %d bytes (original: %d)", imageUploadMaxBytes, len(data))
}

type oapiMediaUploadResult struct {
	mediaID     string
	cleanID     string
	downloadURL string
}

func uploadImageToOapi(ctx context.Context, cardCli *aiCardClient, imagePath string) (*oapiMediaUploadResult, error) {
	imgData, fileName, err := compressImageForUpload(imagePath)
	if err != nil {
		return nil, fmt.Errorf("compress: %w", err)
	}

	token, err := cardCli.accessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("media", fileName)
	if err != nil {
		return nil, err
	}
	if _, err := fw.Write(imgData); err != nil {
		return nil, err
	}
	w.Close()

	uploadURL := "https://oapi.dingtalk.com/media/upload?" + url.Values{
		"access_token": {token},
		"type":         {"image"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		msg := string(body)
		if len(msg) > 200 {
			msg = msg[:200]
		}
		return nil, fmt.Errorf("upload HTTP %d: %s", resp.StatusCode, msg)
	}

	var out struct {
		MediaId string `json:"media_id"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("parse upload response: %w", err)
	}
	if out.ErrCode != 0 {
		return nil, fmt.Errorf("upload error %d: %s", out.ErrCode, out.ErrMsg)
	}
	if out.MediaId == "" {
		msg := string(body)
		if len(msg) > 200 {
			msg = msg[:200]
		}
		return nil, fmt.Errorf("upload response missing media_id: %s", msg)
	}

	cleanID := out.MediaId
	if idx := strings.LastIndex(cleanID, "/"); idx >= 0 {
		cleanID = cleanID[idx+1:]
	}
	cleanID = strings.TrimPrefix(cleanID, "@")

	result := &oapiMediaUploadResult{
		mediaID:     out.MediaId,
		cleanID:     cleanID,
		downloadURL: "https://down.dingtalk.com/media/" + cleanID,
	}
	fmt.Fprintf(os.Stderr, "[connect] 图片已上传 oapi (media_id=%s, %d bytes)\n", out.MediaId, len(imgData))
	return result, nil
}

func sendImageViaBatchSend(ctx context.Context, cardCli *aiCardClient, mediaID string, userIDs []string) error {
	token, err := cardCli.accessToken(ctx)
	if err != nil {
		return err
	}

	msgParam, _ := json.Marshal(map[string]string{"photoURL": mediaID})
	payload := map[string]any{
		"robotCode": cardCli.clientID,
		"userIds":   userIDs,
		"msgKey":    "sampleImageMsg",
		"msgParam":  string(msgParam),
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		dingtalkCardAPIBase+"/v1.0/robot/oToMessages/batchSend", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-acs-dingtalk-access-token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		msg := string(respBody)
		if len(msg) > 200 {
			msg = msg[:200]
		}
		return fmt.Errorf("batchSend HTTP %d: %s", resp.StatusCode, msg)
	}
	return nil
}

// processReplyImages uploads local images found in the reply text, replaces
// their paths with markdown image links, strips useless bash commands, and
// sends the images as separate messages to the sender. Returns the cleaned
// reply text.
func processReplyImages(ctx context.Context, cardCli *aiCardClient, senderStaffID string, replyText string) string {
	paths := extractImagePaths(replyText)
	if len(paths) == 0 {
		return replyText
	}

	cleaned := replyText
	for _, p := range paths {
		result, err := uploadImageToOapi(ctx, cardCli, p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[connect] 图片上传失败 (%s): %v\n", p, err)
			continue
		}

		cleaned = strings.ReplaceAll(cleaned, "`"+p+"`", "!["+filepath.Base(p)+"]("+result.downloadURL+")")
		cleaned = strings.ReplaceAll(cleaned, p, "!["+filepath.Base(p)+"]("+result.downloadURL+")")

		if senderStaffID != "" {
			imgCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			if err := sendImageViaBatchSend(imgCtx, cardCli, result.mediaID, []string{senderStaffID}); err != nil {
				fmt.Fprintf(os.Stderr, "[connect] 图片发送失败 (%s): %v\n", p, err)
			} else {
				fmt.Fprintf(os.Stderr, "[connect] 图片已发送 (%s → %s, mediaId=%s)\n", p, senderStaffID, result.mediaID)
			}
			cancel()
		}
	}

	cleaned = bashOpenCmdPattern.ReplaceAllString(cleaned, "")
	cleaned = viewHintPattern.ReplaceAllString(cleaned, "")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		cleaned = "图片已生成"
	}

	return cleaned
}
