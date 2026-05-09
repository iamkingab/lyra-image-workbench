package pixhost

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	UploadURL = "https://api.pixhost.to/images"
	MaxBytes  = 10 * 1024 * 1024
)

var allowedTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/gif":  "gif",
}

type Client struct {
	httpClient *http.Client
}

type Result struct {
	ShowURL  string `json:"showUrl"`
	ThumbURL string `json:"thumbUrl,omitempty"`
	Name     string `json:"name,omitempty"`
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: 120 * time.Second}}
}

func (c *Client) UploadFile(ctx context.Context, path string, mime string, fileName string) (Result, error) {
	normalizedMime := normalizeMime(mime)
	ext, ok := allowedTypes[normalizedMime]
	if !ok {
		return Result{}, errors.New("PiXhost 仅支持 JPG、PNG、GIF 图片")
	}
	info, err := os.Stat(path)
	if err != nil {
		return Result{}, err
	}
	if info.Size() > MaxBytes {
		return Result{}, errors.New("PiXhost 单张图片最大 10MB")
	}
	file, err := os.Open(path)
	if err != nil {
		return Result{}, err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="img"; filename="%s"`, escapeFilename(normalizeFileName(fileName, ext))))
	header.Set("Content-Type", normalizedMime)
	part, err := writer.CreatePart(header)
	if err != nil {
		return Result{}, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return Result{}, err
	}
	_ = writer.WriteField("content_type", "0")
	_ = writer.WriteField("max_th_size", "420")
	if err := writer.Close(); err != nil {
		return Result{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, UploadURL, &body)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Result{}, errors.New(readErrorMessage(res))
	}

	var payload struct {
		ShowURL  string `json:"show_url"`
		ThumbURL string `json:"th_url"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return Result{}, err
	}
	if strings.TrimSpace(payload.ShowURL) == "" {
		return Result{}, errors.New("PiXhost 未返回图片 URL")
	}
	return Result{
		ShowURL:  toDirectImageURL(payload.ShowURL),
		ThumbURL: normalizePublicURL(payload.ThumbURL),
		Name:     payload.Name,
	}, nil
}

func normalizeMime(mime string) string {
	mime = strings.ToLower(strings.TrimSpace(strings.Split(mime, ";")[0]))
	if mime == "image/jpg" {
		return "image/jpeg"
	}
	return mime
}

func normalizeFileName(name string, ext string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "." || base == string(filepath.Separator) || base == "" {
		base = "ai-image." + ext
	}
	base = regexp.MustCompile(`[\\/:*?"<>|]+`).ReplaceAllString(base, "-")
	base = regexp.MustCompile(`\s+`).ReplaceAllString(base, "-")
	if len(base) > 96 {
		base = base[:96]
	}
	if regexp.MustCompile(`\.[a-z0-9]{2,5}$`).MatchString(strings.ToLower(base)) {
		return base
	}
	return base + "." + ext
}

func escapeFilename(value string) string {
	return strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(value)
}

func normalizePublicURL(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "//") {
		return "https:" + value
	}
	return value
}

func toDirectImageURL(value string) string {
	normalized := normalizePublicURL(value)
	parsed, err := http.NewRequest(http.MethodGet, normalized, nil)
	if err != nil || parsed.URL == nil {
		return normalized
	}
	match := regexp.MustCompile(`^/show/([^/]+)/(.+)$`).FindStringSubmatch(parsed.URL.Path)
	if len(match) == 3 && strings.HasSuffix(strings.ToLower(parsed.URL.Hostname()), "pixhost.to") {
		return fmt.Sprintf("https://img2.pixhost.to/images/%s/%s", match[1], match[2])
	}
	return normalized
}

func readErrorMessage(res *http.Response) string {
	data, err := io.ReadAll(io.LimitReader(res.Body, 4096))
	if err != nil {
		return fmt.Sprintf("PiXhost 上传失败：HTTP %d", res.StatusCode)
	}
	var payload map[string]any
	if json.Unmarshal(data, &payload) == nil {
		for _, key := range []string{"message", "error"} {
			if value, ok := payload[key].(string); ok && strings.TrimSpace(value) != "" {
				return value
			}
		}
	}
	if msg := strings.TrimSpace(string(data)); msg != "" {
		return msg
	}
	return fmt.Sprintf("PiXhost 上传失败：HTTP %d", res.StatusCode)
}
