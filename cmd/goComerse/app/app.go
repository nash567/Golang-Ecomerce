package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pressly/goose"

	"github.com/gocomerse/config"
	db "github.com/gocomerse/internal/db"
	"github.com/gocomerse/internal/logger"
	logModel "github.com/gocomerse/internal/logger/model"
)

type Application struct {
	// db         *sql.DB
	log        logModel.Logger
	cfg        *config.AppConfig
	httpServer *http.Server
	db         *sql.DB
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

	// Initialize database
	a.db, err = setupDB(a.cfg.Database)
	if err != nil {
		a.log.WithError(err).WithFields(logModel.Fields{
			"name":   a.cfg.Database.Name,
			"host":   a.cfg.Database.Host,
			"driver": a.cfg.Database.Driver,
			"port":   a.cfg.Database.Port,
			"user":   a.cfg.Database.User,
		}).Fatal("failed to create database connection")
	}
	a.log.WithField("host", a.cfg.Database.Host).WithField("port", a.cfg.Database.Port).Info("created database connection successfully")

	//nolint:gosec
	a.httpServer = &http.Server{Addr: fmt.Sprintf("%v:%v", a.cfg.Server.Host, a.cfg.Server.Port), Handler: registerHTTPEndpoints()}

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

func registerHTTPEndpoints() *http.ServeMux {
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

// Migrate executes SQL migrations.

func (a *Application) Migrate(migrationPath string) {
	goose.SetLogger(a.log.ToStdLogger())
	goose.SetVerbose(a.cfg.Database.Migrations.Verbose)
	if err := goose.SetDialect(a.cfg.Database.Migrations.Dialect); err != nil {
		a.log.WithError(err).Fatal("could not set dialect for sql migrations")
	}

	if err := goose.Up(a.db, migrationPath); err != nil {
		a.log.WithError(err).Fatal("could not execute migrations")
	}

}

// Verify ensures the connection is available for use.
func Verify(conn *sql.DB) error {
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// setup Db

func setupDB(cfg *db.Config) (*sql.DB, error) {
	conn, err := db.NewConnection(cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Verify(conn); err != nil {
		return nil, fmt.Errorf("failed to verify database connection: %w", err)
	}

	return conn, nil
}
