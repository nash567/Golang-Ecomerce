package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocomerse/config"
	"github.com/gocomerse/internal/logger"
	logModel "github.com/gocomerse/internal/logger/model"
	_ "github.com/lib/pq"
)

type Application struct {
	db         *sql.DB
	log        logModel.Logger
	cfg        *config.AppConfig
	httpServer *http.Server
}

func (a *Application) Init(ctx context.Context, configFiles string) {
	var err error
	// Initialize app config
	a.cfg, err = config.LoadConfig(strings.Split(configFiles, ",")...)
	if err != nil {
		log.Fatalf("Failed to load config")
	}

	// Initialize logger

	a.log, err = logger.NewZapLogger(logModel.Config{Level: a.cfg.Logger.Level})
	if err != nil {
		panic(err)
	}

	a.log = a.log.WithFields(logModel.Fields{
		"appName": a.cfg.APPName,
		"env":     a.cfg.Env,
	})
	a.httpServer = &http.Server{Addr: fmt.Sprintf("%v:%v", a.cfg.Server.Host, a.cfg.Server.Port), Handler: registerHTTPEndpoints(a)}

}

func (a *Application) Start(ctx context.Context) {
	go func() {
		a.log.WithField("addr", a.httpServer.Addr).Info("http server started")

		if err := a.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.WithError(err).Fatal("Error running API")
			} else {
				a.log.WithError(err).Info("Stopping API")
			}
		}
	}()
}

func registerHTTPEndpoints(a *Application) *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}

func (a *Application) Stop(ctx context.Context) {

	// Shutting down the http server
	if a.httpServer != nil {
		a.log.Info("Shutting down the http server...")
		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			a.log.WithError(err).Error("The http server wasn't shut down gracefully")
		}
	}
	// Shutting down the entity rpc server
	a.log.Info("Shutting down the entity rpc server...")
}
