package prompttools

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/config"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/jobs"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/llm"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/output"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/settings"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/spaceconfig"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/uploads"
)

type Service struct {
	store       *Store
	settings    *settings.FileStore
	spaceConfig *spaceconfig.Store
	uploads     *uploads.Store
	jobs        *jobs.Manager
	output      *output.Store
	llm         *llm.Client
}

func NewService(store *Store, settingsStore *settings.FileStore, spaceConfig *spaceconfig.Store, uploadStore *uploads.Store, jobManager *jobs.Manager, outputStore *output.Store, llmClient *llm.Client) *Service {
	return &Service{
		store:       store,
		settings:    settingsStore,
		spaceConfig: spaceConfig,
		uploads:     uploadStore,
		jobs:        jobManager,
		output:      outputStore,
		llm:         llmClient,
	}
}

func (s *Service) TextToPrompt(ctx context.Context, spaceToken string, req TextRequest) (Record, error) {
	input := strings.TrimSpace(req.Input)
	if input == "" {
		return Record{}, errors.New("请输入需要扩写的文字想法")
	}
	apiKey, err := s.apiKey(spaceToken)
	if err != nil {
		return Record{}, err
	}
	started := time.Now()
	resp, err := s.llm.Complete(ctx, llm.Request{
		BaseURL:     s.settings.Get().NewAPIBaseURL,
		APIKey:      apiKey,
		Model:       config.DefaultPromptModel,
		System:      textSystemPrompt(),
		User:        textUserPrompt(input, strings.TrimSpace(req.Style), strings.TrimSpace(req.Ratio), strings.TrimSpace(req.Target)),
		TimeoutSec:  config.DefaultPromptTimeoutSec,
		Temperature: 0.4,
	})
	if err != nil {
		return Record{}, err
	}
	parsed := parsePromptJSON(resp.Text)
	record, err := s.newRecord(spaceToken, Record{
		Mode:           ModeTextToPrompt,
		Input:          input,
		Style:          firstString(parsed, "style", strings.TrimSpace(req.Style)),
		Ratio:          firstString(parsed, "ratio", strings.TrimSpace(req.Ratio)),
		Language:       defaultString(req.Language, "zh"),
		Target:         defaultString(req.Target, "image-2"),
		FlatPrompt:     firstString(parsed, "flatPrompt", resp.Text),
		NegativePrompt: firstString(parsed, "negativePrompt", ""),
		MustKeep:       stringSlice(parsed["mustKeep"]),
		Raw:            resp.Text,
		Model:          config.DefaultPromptModel,
		ElapsedMs:      time.Since(started).Milliseconds(),
	})
	if err != nil {
		return Record{}, err
	}
	return record, s.store.Save(spaceToken, record)
}

func (s *Service) ImageToPrompt(ctx context.Context, spaceToken string, req ImageRequest) (Record, error) {
	started := time.Now()
	path, mime, sourceURL, err := s.resolveImageSource(spaceToken, req.Source)
	if err != nil {
		return Record{}, err
	}
	apiKey, err := s.apiKey(spaceToken)
	if err != nil {
		return Record{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Record{}, err
	}
	if len(data) == 0 {
		return Record{}, errors.New("图片内容为空")
	}
	resp, err := s.llm.Complete(ctx, llm.Request{
		BaseURL:    s.settings.Get().NewAPIBaseURL,
		APIKey:     apiKey,
		Model:      config.DefaultPromptModel,
		System:     imageSystemPrompt(),
		User:       imageUserPrompt(strings.TrimSpace(req.Target)),
		Image:      &llm.ImagePart{Mime: mime, Data: data},
		TimeoutSec: config.DefaultPromptTimeoutSec,
	})
	if err != nil {
		return Record{}, err
	}
	parsed := parsePromptJSON(resp.Text)
	record, err := s.newRecord(spaceToken, Record{
		Mode:            ModeImageToPrompt,
		Language:        defaultString(req.Language, "zh"),
		Target:          defaultString(req.Target, "image-2"),
		Source:          req.Source,
		SourceImageURL:  sourceURL,
		FlatPrompt:      firstString(parsed, "flatPrompt", resp.Text),
		NegativePrompt:  firstString(parsed, "negativePrompt", ""),
		MustKeep:        stringSlice(parsed["mustKeep"]),
		Avoid:           stringSlice(parsed["avoid"]),
		JSONDescription: objectMap(parsed["jsonDescription"]),
		Raw:             resp.Text,
		Model:           config.DefaultPromptModel,
		ElapsedMs:       time.Since(started).Milliseconds(),
	})
	if err != nil {
		return Record{}, err
	}
	return record, s.store.Save(spaceToken, record)
}

func (s *Service) List(spaceToken string, limit int) ([]Record, error) {
	return s.store.List(spaceToken, limit)
}

func (s *Service) Delete(spaceToken string, id string) (Record, bool, error) {
	return s.store.Delete(spaceToken, id)
}

func (s *Service) apiKey(spaceToken string) (string, error) {
	cfg, err := s.spaceConfig.Get(spaceToken)
	if err != nil {
		return "", err
	}
	return cfg.APIKey, nil
}

func (s *Service) resolveImageSource(spaceToken string, source Source) (string, string, string, error) {
	switch strings.TrimSpace(source.Type) {
	case "upload":
		item, path, err := s.uploads.GetReferenceImage(spaceToken, source.UploadID)
		if err != nil {
			return "", "", "", err
		}
		return path, item.Mime, "", nil
	case "result":
		job, ok, err := s.jobs.Get(spaceToken, source.TaskID)
		if err != nil {
			return "", "", "", err
		}
		if !ok {
			return "", "", "", errors.New("任务不存在")
		}
		for _, result := range job.Results {
			if result.Index == source.Index && result.OK && result.ImageURL != "" {
				path, mime, err := s.output.ResolveURL(result.ImageURL)
				return path, mime, result.ImageURL, err
			}
		}
		return "", "", "", errors.New("任务图片不存在")
	default:
		return "", "", "", errors.New("图片来源无效")
	}
}

func (s *Service) newRecord(spaceToken string, record Record) (Record, error) {
	if _, err := s.spaceConfig.Get(spaceToken); err != nil {
		return Record{}, err
	}
	id, err := newRecordID()
	if err != nil {
		return Record{}, err
	}
	record.ID = id
	record.CreatedAt = time.Now()
	record.FlatPrompt = strings.TrimSpace(record.FlatPrompt)
	record.NegativePrompt = strings.TrimSpace(record.NegativePrompt)
	if record.FlatPrompt == "" {
		return Record{}, errors.New("提示词模型没有返回可用提示词")
	}
	return record, nil
}

func newRecordID() (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return "ppt_" + time.Now().Format("20060102150405") + "_" + hex.EncodeToString(bytes[:]), nil
}

func parsePromptJSON(text string) map[string]any {
	candidates := []string{strings.TrimSpace(text)}
	if fenced := extractFence(text); fenced != "" {
		candidates = append([]string{fenced}, candidates...)
	}
	if object := extractJSONObject(text); object != "" {
		candidates = append([]string{object}, candidates...)
	}
	for _, candidate := range candidates {
		var out map[string]any
		if json.Unmarshal([]byte(candidate), &out) == nil {
			return out
		}
	}
	return map[string]any{"flatPrompt": strings.TrimSpace(text)}
}

func extractFence(text string) string {
	re := regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")
	match := re.FindStringSubmatch(text)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractJSONObject(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		return text[start : end+1]
	}
	return ""
}

func firstString(values map[string]any, key string, fallback string) string {
	if value, ok := values[key].(string); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return strings.TrimSpace(fallback)
}

func stringSlice(value any) []string {
	switch typed := value.(type) {
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if text := fmt.Sprint(item); strings.TrimSpace(text) != "" {
				out = append(out, strings.TrimSpace(text))
			}
		}
		return out
	case []string:
		return typed
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil
		}
		return []string{strings.TrimSpace(typed)}
	default:
		return nil
	}
}

func objectMap(value any) map[string]any {
	if typed, ok := value.(map[string]any); ok {
		return typed
	}
	return nil
}

func defaultString(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func ParseLimit(value string) int {
	limit, _ := strconv.Atoi(value)
	if limit <= 0 {
		return 50
	}
	return limit
}
