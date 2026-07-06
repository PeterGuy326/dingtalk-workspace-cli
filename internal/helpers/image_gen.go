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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/executor"
	"github.com/spf13/cobra"
)

const (
	imageGenDefaultModel   = "gemini-3.1-flash-image-preview"
	imageGenDefaultBaseURL = "https://idealab.alibaba-inc.com/api/openai/v1"
	imageGenDefaultAPIKey  = "797c5d7a0043d5343d2b52522a448e86"
	imageGenDefaultOutput  = "generated-image.png"
	imageGenRequestTimeout = 120 * time.Second
)

type imageHandler struct{}

func (imageHandler) Name() string { return "image" }

func (imageHandler) Command(_ executor.Runner) *cobra.Command {
	return newImageCommand()
}

func init() {
	RegisterPublic(func() Handler {
		return imageHandler{}
	})
}

func newImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "AI 图片生成（调用 idealab Gemini 图像模型）",
	}
	cmd.AddCommand(newImageGenerateCommand())
	return cmd
}

func newImageGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "根据 prompt 生成图片",
		Example: `  dws image generate --prompt "一只可爱的猫咪"
  dws image generate --prompt "red square on white background" --output /tmp/square.png
  dws image generate --prompt "landscape" --model gemini-3.1-flash-image-preview`,
		RunE: runImageGenerate,
	}
	cmd.Flags().String("prompt", "", "图片生成描述（必填）")
	cmd.Flags().String("output", imageGenDefaultOutput, "输出文件路径")
	cmd.Flags().String("model", imageGenDefaultModel, "图像生成模型")
	cmd.Flags().String("api-key", "", "API key（默认读 env DWS_IDEALAB_API_KEY 或 opencode 配置）")
	cmd.Flags().String("base-url", "", "API base URL（默认读 env DWS_IDEALAB_BASE_URL 或 opencode 配置）")
	_ = cmd.MarkFlagRequired("prompt")
	return cmd
}

func runImageGenerate(cmd *cobra.Command, _ []string) error {
	prompt, _ := cmd.Flags().GetString("prompt")
	output, _ := cmd.Flags().GetString("output")
	model, _ := cmd.Flags().GetString("model")
	apiKey, _ := cmd.Flags().GetString("api-key")
	baseURL, _ := cmd.Flags().GetString("base-url")

	if apiKey == "" {
		apiKey = resolveImageAPIKey()
	}
	if baseURL == "" {
		baseURL = resolveImageBaseURL()
	}
	if apiKey == "" {
		return fmt.Errorf("API key not found: set --api-key, env DWS_IDEALAB_API_KEY, or configure in opencode config")
	}
	if baseURL == "" {
		baseURL = imageGenDefaultBaseURL
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "generating image: model=%s prompt=%.60s...\n", model, prompt)

	imgData, err := callImageGenAPI(cmd.Context(), baseURL, apiKey, model, prompt)
	if err != nil {
		return fmt.Errorf("image generation failed: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(output, imgData, 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "image saved: %s (%d bytes)\n", output, len(imgData))
	return nil
}

func callImageGenAPI(ctx context.Context, baseURL, apiKey, model, prompt string) ([]byte, error) {
	reqBody := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 4096,
	}
	bodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, imageGenRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(respBody))
	}

	return extractImageFromResponse(respBody)
}

func extractImageFromResponse(body []byte) ([]byte, error) {
	var resp struct {
		Choices []struct {
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := resp.Choices[0].Message.Content

	var parts []struct {
		ImageURL struct {
			URL string `json:"url"`
		} `json:"image_url"`
	}
	if err := json.Unmarshal(content, &parts); err == nil {
		for _, part := range parts {
			if part.ImageURL.URL != "" {
				return decodeBase64Image(part.ImageURL.URL)
			}
		}
	}

	var textContent string
	if err := json.Unmarshal(content, &textContent); err == nil {
		return nil, fmt.Errorf("model returned text instead of image: %.200s", textContent)
	}

	return nil, fmt.Errorf("no image found in response")
}

func decodeBase64Image(data string) ([]byte, error) {
	if idx := strings.Index(data, ","); idx != -1 && strings.HasPrefix(data, "data:") {
		data = data[idx+1:]
	}
	return base64.StdEncoding.DecodeString(data)
}

func resolveImageAPIKey() string {
	if v := os.Getenv("DWS_IDEALAB_API_KEY"); v != "" {
		return v
	}
	return readOpencodeConfigAPIKey()
}

func resolveImageBaseURL() string {
	if v := os.Getenv("DWS_IDEALAB_BASE_URL"); v != "" {
		return v
	}
	return readOpencodeConfigBaseURL()
}

func readOpencodeConfigAPIKey() string {
	cfg := readOpencodeConfig()
	if cfg == nil {
		return imageGenDefaultAPIKey
	}
	return cfg.apiKey
}

func readOpencodeConfigBaseURL() string {
	cfg := readOpencodeConfig()
	if cfg == nil {
		return imageGenDefaultBaseURL
	}
	return cfg.baseURL
}

type opencodeIdealabConfig struct {
	apiKey  string
	baseURL string
}

func readOpencodeConfig() *opencodeIdealabConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := filepath.Join(home, ".config", "opencode", "opencode.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var cfg struct {
		Provider map[string]struct {
			Options struct {
				APIKey  string `json:"apiKey"`
				BaseURL string `json:"baseURL"`
			} `json:"options"`
		} `json:"provider"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}

	if idealab, ok := cfg.Provider["idealab"]; ok {
		return &opencodeIdealabConfig{
			apiKey:  idealab.Options.APIKey,
			baseURL: idealab.Options.BaseURL,
		}
	}
	return nil
}
