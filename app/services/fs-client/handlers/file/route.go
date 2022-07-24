package file

import (
	"net/http"

	"gitlab.com/onurkeskin/go-http-fs-client/domain/core/file_service"
	"gitlab.com/onurkeskin/go-http-fs-client/foundation/web"
	"go.uber.org/zap"
)

type Config struct {
	ServerUrl        string
	DownloadLocation string
	Log              *zap.SugaredLogger
}

func Routes(app *web.App, cfg Config) {
	frh := FileRetriverHandler{
		fileService: file_service.NewFileService(cfg.ServerUrl, cfg.DownloadLocation),
	}
	app.Handle(http.MethodGet, "/:token", frh.RetrieveFiles)
}
