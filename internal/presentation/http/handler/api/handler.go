package api

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hexagonal/internal/adapter/http/middleware"
	"hexagonal/internal/core/infrastructure/config"
	"hexagonal/internal/core/service/taskSRV"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultTimeOutForGracefulShutDown = 5 * time.Second
	IdleTimeout                       = time.Second * 60
	ReadTimeout                       = time.Second * 15
	WriteTimeout                      = time.Second * 15
)

type Handler struct {
	Cfg         *config.Config
	TaskService *taskSRV.TaskService
	httpServer  *http.Server
}

// CreateHandler creates a new instance of the Handler
func CreateHandler(cfg *config.Config, taskService *taskSRV.TaskService) *Handler {
	return &Handler{
		Cfg:         cfg,
		TaskService: taskService,
	}
}

// SetupRouter initializes the Gin router and applies middlewares
func (h *Handler) SetupRouter() *gin.Engine {

	gin.SetMode(selectMode(h.Cfg.App.Debug))

	// Initialize the router
	router := gin.New()

	// Apply common middlewares
	router.Use(cors.Default())
	router.Use(gin.Recovery())
	router.Use(middleware.SomeMiddleWare())
	h.RegisterHexaRoutes(router)

	return router
}

// StartServer starts the http server in blocking mode
func (h *Handler) StartServer() {

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(h.Cfg.App.Port),
		Handler:      h.SetupRouter(),
		WriteTimeout: WriteTimeout,
		ReadTimeout:  ReadTimeout,
		IdleTimeout:  IdleTimeout,
	}

	h.httpServer = server

	err := h.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatal(err)
	}

	// Code reaches here after HTTP server shutdown
	logrus.Info("[Handler] HTTP REST Server is shutting down!")
}

// StopServer handles the http server in graceful shutdown
func (h *Handler) StopServer() {
	ctxTimeout, cancelTimeout := context.WithTimeout(
		context.Background(),
		DefaultTimeOutForGracefulShutDown,
	)

	defer cancelTimeout()

	h.httpServer.SetKeepAlivesEnabled(false)

	if err := h.httpServer.Shutdown(ctxTimeout); err != nil {
		logrus.Error(err)
	}

	logrus.Info("[Handler] HTTP REST Server graceful shutdown completed")
}

// selectMode selects the Gin mode based on debug flag
func selectMode(debug bool) string {
	if debug {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}
