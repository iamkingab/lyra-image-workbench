package uploads

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/spaces"
)

const (
	MaxReferenceImages      = 8
	MaxReferenceImageBytes  = 12 * 1024 * 1024
	MaxReferenceUploadBytes = 50 * 1024 * 1024
)

var allowedImageTypes = map[string]string{
	"image/png":  "png",
	"image/jpeg": "jpg",
	"image/webp": "webp",
}

type ReferenceImage struct {
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
	FileName     string `json:"fileName"`
	Mime         string `json:"mime"`
	Size         int64  `json:"size"`
	CreatedAt    string `json:"createdAt"`
}

type Store struct {
	spaces *spaces.FileStore
}

func NewStore(spaceStore *spaces.FileStore) *Store {
	return &Store{spaces: spaceStore}
}

func (s *Store) SaveReferenceImages(spaceToken string, headers []*multipart.FileHeader) ([]ReferenceImage, error) {
	if len(headers) == 0 {
		return nil, NewUploadError("REFERENCE_IMAGE_MISSING", "请先上传图生图参考图")
	}
	if len(headers) > MaxReferenceImages {
		return nil, NewUploadError("REFERENCE_IMAGE_TOO_MANY", "参考图最多 8 张")
	}

	spaceDir, err := s.spaces.SpaceDir(spaceToken)
	if err != nil {
		return nil, err
	}
	uploadDir := filepath.Join(spaceDir, "uploads")
	if err := os.MkdirAll(uploadDir, 0o700); err != nil {
		return nil, err
	}

	items := make([]ReferenceImage, 0, len(headers))
	for _, header := range headers {
		item, err := s.saveOne(uploadDir, header)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Store) saveOne(uploadDir string, header *multipart.FileHeader) (ReferenceImage, error) {
	if header.Size > MaxReferenceImageBytes {
		return ReferenceImage{}, NewUploadError("REFERENCE_IMAGE_TOO_LARGE", "单张参考图不能超过 12MB")
	}

	file, err := header.Open()
	if err != nil {
		return ReferenceImage{}, err
	}
	defer file.Close()

	limited := io.LimitReader(file, MaxReferenceImageBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return ReferenceImage{}, err
	}
	if int64(len(data)) > MaxReferenceImageBytes {
		return ReferenceImage{}, NewUploadError("REFERENCE_IMAGE_TOO_LARGE", "单张参考图不能超过 12MB")
	}

	mime := detectMime(data, header.Header.Get("Content-Type"))
	ext, ok := allowedImageTypes[mime]
	if !ok {
		return ReferenceImage{}, NewUploadError("REFERENCE_IMAGE_TYPE_UNSUPPORTED", "参考图仅支持 PNG、JPG、WEBP")
	}

	id, err := newID()
	if err != nil {
		return ReferenceImage{}, err
	}
	now := time.Now().Format(time.RFC3339)
	originalName := safeOriginalName(header.Filename, ext)
	fileName := fmt.Sprintf("%s.%s", id, ext)
	imagePath := filepath.Join(uploadDir, fileName)
	metaPath := filepath.Join(uploadDir, fmt.Sprintf("%s.json", id))

	if err := writeFileAtomic(imagePath, data, 0o600); err != nil {
		return ReferenceImage{}, err
	}
	item := ReferenceImage{
		ID:           id,
		OriginalName: originalName,
		FileName:     fileName,
		Mime:         mime,
		Size:         int64(len(data)),
		CreatedAt:    now,
	}
	meta, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return ReferenceImage{}, err
	}
	if err := writeFileAtomic(metaPath, append(meta, '\n'), 0o600); err != nil {
		return ReferenceImage{}, err
	}
	return item, nil
}

func detectMime(data []byte, declared string) string {
	declared = strings.ToLower(strings.TrimSpace(strings.Split(declared, ";")[0]))
	if _, ok := allowedImageTypes[declared]; ok {
		return declared
	}
	return http.DetectContentType(data)
}

func newID() (string, error) {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return "ref_" + hex.EncodeToString(bytes[:]), nil
}

func safeOriginalName(name string, ext string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "." || base == string(filepath.Separator) || base == "" {
		base = "reference." + ext
	}
	base = regexp.MustCompile(`[\\/:*?"<>|]+`).ReplaceAllString(base, "-")
	base = regexp.MustCompile(`\s+`).ReplaceAllString(base, "-")
	if len(base) > 96 {
		base = base[:96]
	}
	return base
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp := fmt.Sprintf("%s.%d.tmp", path, time.Now().UnixNano())
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

type UploadError struct {
	Code    string
	Chinese string
}

func NewUploadError(code string, chinese string) UploadError {
	return UploadError{Code: code, Chinese: chinese}
}

func (e UploadError) Error() string {
	return e.Chinese
}

func AsUploadError(err error, target *UploadError) bool {
	return errors.As(err, target)
}
