package prompttools

import "time"

type Mode string

const (
	ModeTextToPrompt  Mode = "text-to-prompt"
	ModeImageToPrompt Mode = "image-to-prompt"
)

type Source struct {
	Type     string `json:"type"`
	UploadID string `json:"uploadId,omitempty"`
	TaskID   string `json:"taskId,omitempty"`
	Index    int    `json:"index,omitempty"`
}

type TextRequest struct {
	Input    string `json:"input"`
	Style    string `json:"style"`
	Ratio    string `json:"ratio"`
	Language string `json:"language"`
	Target   string `json:"target"`
}

type ImageRequest struct {
	Source   Source `json:"source"`
	Language string `json:"language"`
	Target   string `json:"target"`
}

type Record struct {
	ID              string         `json:"id"`
	Mode            Mode           `json:"mode"`
	Input           string         `json:"input,omitempty"`
	Style           string         `json:"style,omitempty"`
	Ratio           string         `json:"ratio,omitempty"`
	Language        string         `json:"language,omitempty"`
	Target          string         `json:"target,omitempty"`
	Source          Source         `json:"source,omitempty"`
	SourceImageURL  string         `json:"sourceImageUrl,omitempty"`
	FlatPrompt      string         `json:"flatPrompt"`
	NegativePrompt  string         `json:"negativePrompt,omitempty"`
	MustKeep        []string       `json:"mustKeep,omitempty"`
	Avoid           []string       `json:"avoid,omitempty"`
	JSONDescription map[string]any `json:"jsonDescription,omitempty"`
	Raw             string         `json:"raw,omitempty"`
	Model           string         `json:"model"`
	ElapsedMs       int64          `json:"elapsedMs"`
	CreatedAt       time.Time      `json:"createdAt"`
}
