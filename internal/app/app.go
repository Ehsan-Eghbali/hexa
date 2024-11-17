package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"hexagonal/internal/adapter/http/handler/api"
	"hexagonal/internal/adapter/repository/task"
	"hexagonal/internal/core/infrastructure/config"
	"hexagonal/internal/core/infrastructure/db"
	logutil "hexagonal/internal/core/infrastructure/log"
	"hexagonal/internal/core/service/taskSRV"
	"os"
	"sync"
)

// App encapsulates the application's core services.
type App struct {
	cfg         *config.Config
	DBRegistry  *db.Registry
	TaskService *taskSRV.TaskService
	Handler     *api.Handler

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

	cfg := config.Instance
	// Initialize dbPO registry
	dbRegistry, err := db.NewDBRegistry(cfg)
	if err != nil {
		logutil.LogOnce("err_load_db_registry", err, nil)
		return nil, err
	}
	logutil.LogSuccess("load_db_registry", nil)

	taskRepo := task.NewRedisTaskRepository(dbRegistry)
	taskService := taskSRV.NewTaskService(taskRepo)
	app.DBRegistry = dbRegistry
	app.TaskService = taskService

	app.Handler = api.CreateHandler(cfg, taskService)
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
