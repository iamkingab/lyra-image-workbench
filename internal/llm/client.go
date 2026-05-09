package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
}

type Message struct {
	Role    string
	Content any
}

type ImagePart struct {
	Mime string
	Data []byte
}

type Request struct {
	BaseURL     string
	APIKey      string
	Model       string
	System      string
	User        string
	Image       *ImagePart
	TimeoutSec  int
	Temperature float64
}

type Response struct {
	Text string
}

type UpstreamError struct {
	StatusCode int
	Message    string
}

func (e UpstreamError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("提示词模型请求失败：HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("提示词模型请求失败：HTTP %d：%s", e.StatusCode, e.Message)
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{}}
}

func (c *Client) Complete(ctx context.Context, req Request) (Response, error) {
	if strings.TrimSpace(req.BaseURL) == "" {
		return Response{}, errors.New("提示词模型 URL 为空")
	}
	if strings.TrimSpace(req.APIKey) == "" {
		return Response{}, errors.New("请先在当前个人空间填写 Image-2 Key，提示词工具会复用该 Key")
	}
	if strings.TrimSpace(req.Model) == "" {
		return Response{}, errors.New("提示词模型为空")
	}
	timeout := time.Duration(req.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 180 * time.Second
	}
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	body := map[string]any{
		"model": req.Model,
		"messages": []map[string]any{
			{"role": "system", "content": req.System},
			{"role": "user", "content": userContent(req.User, req.Image)},
		},
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return Response{}, err
	}
	httpReq, err := http.NewRequestWithContext(callCtx, http.MethodPost, buildURL(req.BaseURL, "chat/completions"), bytes.NewReader(payload))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Cache-Control", "no-store")

	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Response{}, UpstreamError{StatusCode: res.StatusCode, Message: readErrorMessage(res)}
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content any `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return Response{}, err
	}
	if len(out.Choices) == 0 {
		return Response{}, errors.New("提示词模型没有返回内容")
	}
	text := contentToText(out.Choices[0].Message.Content)
	if strings.TrimSpace(text) == "" {
		return Response{}, errors.New("提示词模型返回内容为空")
	}
	return Response{Text: strings.TrimSpace(text)}, nil
}

func userContent(text string, image *ImagePart) any {
	if image == nil || len(image.Data) == 0 {
		return text
	}
	mime := strings.TrimSpace(strings.Split(image.Mime, ";")[0])
	if !strings.HasPrefix(strings.ToLower(mime), "image/") {
		mime = "image/png"
	}
	return []map[string]any{
		{"type": "text", "text": text},
		{"type": "image_url", "image_url": map[string]any{"url": dataURL(mime, image.Data)}},
	}
}

func contentToText(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		parts := make([]string, 0)
		for _, item := range typed {
			if m, ok := item.(map[string]any); ok {
				if text, ok := m["text"].(string); ok {
					parts = append(parts, text)
				}
			}
		}
		return strings.Join(parts, "\n")
	default:
		data, _ := json.Marshal(typed)
		return string(data)
	}
}

func dataURL(mime string, data []byte) string {
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data)
}

func buildURL(baseURL string, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func readErrorMessage(res *http.Response) string {
	data, err := io.ReadAll(io.LimitReader(res.Body, 4096))
	if err != nil {
		return ""
	}
	var payload map[string]any
	if json.Unmarshal(data, &payload) == nil {
		if errValue, ok := payload["error"].(map[string]any); ok {
			if msg, ok := errValue["message"].(string); ok {
				return msg
			}
		}
		if msg, ok := payload["message"].(string); ok {
			return msg
		}
	}
	return strings.TrimSpace(string(data))
}
