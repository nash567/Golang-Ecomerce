package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pressly/goose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userRPC "github.com/gocomerse/api/rpc/gocomerse/user"
	"github.com/gocomerse/config"
	db "github.com/gocomerse/internal/db"
	"github.com/gocomerse/internal/logger"
	logModel "github.com/gocomerse/internal/logger/model"
	userPB "github.com/gocomerse/internal/pb/gocomerse/user"
)

const timeout = 5 * time.Second

type Application struct {
	log        logModel.Logger
	cfg        *config.AppConfig
	httpServer *http.Server
	grpcServer *grpc.Server
	services   *services
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

	a.services = buildServices(a.db)

	a.grpcServer = rgisterGRPCEndpoints(a.services, a.log)

	//nolint:gosec
	a.httpServer = &http.Server{
		Addr:              fmt.Sprintf("%v:%v", a.cfg.Server.Host, a.cfg.Server.Port),
		Handler:           registerHTTPEndpoints(ctx, a.cfg, a.log),
		ReadHeaderTimeout: timeout,
	}
}

func (a *Application) Start(ctx context.Context) {

	// grpc server
	go func() {
		listenOn := fmt.Sprintf("%v:%v", a.cfg.Server.Host, a.cfg.Server.GRPCPort)
		listener, listenerErr := net.Listen("tcp", listenOn)
		if listenerErr != nil {
			a.log.WithError(listenerErr).Fatal(fmt.Sprintf("failed to listen on %s", listenOn))
		}
		a.log.WithField("addr", listenOn).Info("grpc server started")
		if err := a.grpcServer.Serve(listener); err != nil {
			a.log.WithError(err).Fatal("failed to serve gRPC server")
		}
	}()

	// grpc gateways

	go func() {
		listenOn := fmt.Sprintf("%v:%v", a.cfg.Server.Host, a.cfg.Server.Port)

		listner, listenerErr := net.Listen("tcp",
			fmt.Sprintf("%v:%v", a.cfg.Server.Host, a.cfg.Server.Port),
		)
		if listenerErr != nil {
			a.log.WithError(listenerErr).Fatal(fmt.Sprintf("failed to listen on %s", listenOn))

		}

		a.log.WithField("addr", listenOn).Info("grpc gateway server started")

		if err := a.httpServer.Serve(listner); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.WithError(err).Fatal("Error running API")
			} else {
				a.log.WithError(err).Info("Stopping API")
			}
		}
	}()
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
	a.grpcServer.GracefulStop()
}

func rgisterGRPCEndpoints(services *services, log logModel.Logger) *grpc.Server {
	grpcServer := grpc.NewServer()
	userPB.RegisterUserServiceServer(grpcServer, userRPC.NewServer(services.UserSvc, log))
	return grpcServer
}

func registerHTTPEndpoints(ctx context.Context, cfg *config.AppConfig, log logModel.Logger) *runtime.ServeMux {
	mux := runtime.NewServeMux()
	err := userPB.RegisterUserServiceHandlerFromEndpoint(ctx,
		mux,
		fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.WithError(err).Errorf("Error while registering user service: %w", err)
	}
	return mux
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

	// to down we use gooseDown to version 0

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
