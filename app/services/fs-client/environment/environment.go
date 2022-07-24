package environment

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type EnvironmentConfiguration struct {
	Web struct {
		ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
		IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"20s"`
		APIHost         string        `env:"API_HOST" envDefault:"localhost:8080"`
	}
	FileServer struct {
		FileServerUrl    string `env:"FILE_SERVER_URL" envDefault:"0.0.0.0:8081"`
		DownloadLocation string `env:"DOWNLOAD_LOCATION" envDefault:"/tmp"`
	}
}

func (environmentConfiguration *EnvironmentConfiguration) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Read Timeout", environmentConfiguration.Web.ReadTimeout.String())
	enc.AddString("Write Timeout", environmentConfiguration.Web.WriteTimeout.String())
	enc.AddString("Idle Timeout", environmentConfiguration.Web.IdleTimeout.String())
	enc.AddString("Shutdown Timeout", environmentConfiguration.Web.ShutdownTimeout.String())
	enc.AddString("Client Api Host", environmentConfiguration.Web.APIHost)
	enc.AddString("File Server Host", environmentConfiguration.FileServer.FileServerUrl)
	enc.AddString("Download Location", environmentConfiguration.FileServer.DownloadLocation)

	return nil
}
