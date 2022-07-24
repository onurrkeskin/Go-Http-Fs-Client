package handlers

import (
	"net/http"
	"os"

	"gitlab.com/onurkeskin/go-http-fs-client/app/services/fs-client/environment"
	"gitlab.com/onurkeskin/go-http-fs-client/app/services/fs-client/handlers/file"
	"gitlab.com/onurkeskin/go-http-fs-client/domain/mid"
	"gitlab.com/onurkeskin/go-http-fs-client/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type ApiConfig struct {
	Shutdown          chan os.Signal
	Log               *zap.SugaredLogger
	EnvironmentConfig *environment.EnvironmentConfiguration
}

func APIMux(cfg ApiConfig) http.Handler {
	var app *web.App = web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Panics(),
	)

	file.Routes(app, file.Config{
		Log:              cfg.Log,
		ServerUrl:        cfg.EnvironmentConfig.FileServer.FileServerUrl,
		DownloadLocation: cfg.EnvironmentConfig.FileServer.DownloadLocation,
	})

	return app
}
