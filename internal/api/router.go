package api

import (
	"net/http"

	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/config"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/settings"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/spaceconfig"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/spaces"
	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/uploads"
)

type Dependencies struct {
	Config      config.Config
	Settings    *settings.FileStore
	Spaces      *spaces.FileStore
	SpaceConfig *spaceconfig.Store
	Uploads     *uploads.Store
}

func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()
	health := NewHealthHandler(deps.Config)
	adminConfig := NewAdminConfigHandler(deps.Settings)
	userConfig := NewUserConfigHandler(deps.SpaceConfig)
	statusMetadata := NewStatusMetadataHandler()
	spaceHandler := NewSpaceHandler(deps.Spaces)
	uploadHandler := NewUploadHandler(deps.Uploads)

	mux.HandleFunc("GET /api/health", health.ServeHTTP)
	mux.HandleFunc("GET /api/admin/config", adminConfig.Get)
	mux.HandleFunc("POST /api/admin/config", adminConfig.Update)
	mux.HandleFunc("GET /api/config", userConfig.Get)
	mux.HandleFunc("POST /api/config", userConfig.Update)
	mux.HandleFunc("GET /api/status-metadata", statusMetadata.ServeHTTP)
	mux.HandleFunc("POST /api/spaces/session", spaceHandler.CreateSession)
	mux.HandleFunc("GET /api/spaces/session", spaceHandler.CurrentSession)
	mux.HandleFunc("DELETE /api/spaces/session", spaceHandler.DeleteSession)
	mux.HandleFunc("POST /api/uploads/reference", uploadHandler.SaveReferenceImages)

	return withCommonHeaders(mux)
}

func withCommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}
