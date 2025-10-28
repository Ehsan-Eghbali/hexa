package app

import (
	"context"
	"dariush/config"
	"dariush/internal/core/infrastructure/db"
	"dariush/internal/presentation/http/handler/api"
	logutil "dariush/pkg/logutils"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// App encapsulates the application's core services.
type App struct {
	cfg        *config.Config
	DBRegistry *db.Registry
	Handler    *api.Handler

	sync.WaitGroup
}

func Initialize(ctx context.Context) (*App, error) {
	app := &App{}
	logutil.Init()

	// Load configuration
	err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	cfg := config.Get()
	// Initialize dbPO registry
	dbRegistry, err := db.NewDBRegistry(ctx, cfg)
	if err != nil {
		logutil.LogOnce("err_load_db_registry", err, nil)
		return nil, err
	}
	logutil.LogSuccess("load_db_registry", nil)

	app.DBRegistry = dbRegistry

	app.Handler = api.CreateHandler()
	return app, nil

}

func (app *App) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	// Wait for OS signals
	<-quitSignal

	// Kill the API Endpoints first
	app.Handler.StopServer()

	// Kill the Scheduler nad CloseAll DB connections
	app.Stop()

	close(done)
}

func (app *App) Start(ctx context.Context) {
	// Start the HTTP server in blocking mode
	app.Handler.StartServer()
}

// Stop stops the scheduler and closes database connections
func (app *App) Stop() {
	// Close database connections
	if err := app.DBRegistry.CloseAll(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error closing database connections")
	} else {
		logrus.Info("Database connections closed successfully")
	}
}
