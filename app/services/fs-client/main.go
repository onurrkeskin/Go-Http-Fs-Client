package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/caarlos0/env/v6"
	"gitlab.com/onurkeskin/go-http-fs-client/app/services/fs-client/environment"
	"gitlab.com/onurkeskin/go-http-fs-client/app/services/fs-client/handlers"
	"gitlab.com/onurkeskin/go-http-fs-client/foundation/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New("fs-client-logger")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	opt := maxprocs.Logger(log.Infof)
	if _, err := maxprocs.Set(opt); err != nil {
		log.Errorf("maxprocs: %w", err)
		log.Sync()
		os.Exit(1)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// App Starting

	cfg := environment.EnvironmentConfiguration{}
	if err := env.Parse(&cfg); err != nil {
		log.Infow("Couldn't load configs", "Error", err.Error())
	} else {
		log.Infow("The environment uses configuration", zap.Object("Config", &cfg))
	}

	log.Infow("starting service", "version")
	defer log.Infow("shutdown complete")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct the mux for the API calls.
	apiMux := handlers.APIMux(handlers.ApiConfig{
		Shutdown:          shutdown,
		Log:               log,
		EnvironmentConfig: &cfg,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			log.Errorf("could not stop server gracefully: %w", err)
		}
	}
}
