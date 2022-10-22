package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/gocomerse/cmd/goComerse/app"
)

const (
	defaultConfPath      = "./config.yml,./.env"
	defaultMigrationPath = "./build/db/migrations/"
)

func main() {
	var configFiles string
	var migrationPath string

	flag.StringVar(&configFiles, "config", defaultConfPath, "comma separated list of config files to load")
	flag.StringVar(&migrationPath, "migrations", defaultMigrationPath, "path to SQL migration directory")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := new(app.Application)
	application.Init(ctx, configFiles)
	application.Migrate(migrationPath)

	application.Start(ctx)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// locking till signal received
	<-sigterm
	// start graceful shutdown
	application.Stop(ctx)
}
