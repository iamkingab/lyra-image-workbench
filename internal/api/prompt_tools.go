package api

import (
	"encoding/json"
	"net/http"

	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/prompttools"
)

type PromptToolsHandler struct {
	service *prompttools.Service
}

func NewPromptToolsHandler(service *prompttools.Service) PromptToolsHandler {
	return PromptToolsHandler{service: service}
}

func (h PromptToolsHandler) TextToPrompt(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload prompttools.TextRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "请求体不是有效 JSON")
		return
	}
	record, err := h.service.TextToPrompt(r.Context(), r.Header.Get("X-Space-Token"), payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, "PROMPT_TEXT_FAILED", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "record": record})
}

func (h PromptToolsHandler) ImageToPrompt(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload prompttools.ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "请求体不是有效 JSON")
		return
	}
	record, err := h.service.ImageToPrompt(r.Context(), r.Header.Get("X-Space-Token"), payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, "PROMPT_IMAGE_FAILED", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "record": record})
}

func (h PromptToolsHandler) History(w http.ResponseWriter, r *http.Request) {
	records, err := h.service.List(r.Header.Get("X-Space-Token"), prompttools.ParseLimit(r.URL.Query().Get("limit")))
	if err != nil {
		writeSpaceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "records": records})
}

func (h PromptToolsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	record, ok, err := h.service.Delete(r.Header.Get("X-Space-Token"), r.PathValue("id"))
	if err != nil {
		writeSpaceError(w, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "PROMPT_RECORD_NOT_FOUND", "提示词记录不存在")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "record": record})
}
