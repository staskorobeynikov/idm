package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg := common.GetConfig(".env")
	var logger = common.NewLogger(cfg)
	defer func() { _ = logger.Sync() }()
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error closing db", zap.Error(err))
		}
	}()
	var server = build(cfg, logger, db)
	go func() {
		var err = server.App.Listen(":8080")
		if err != nil {
			logger.Panic("http server error: %s", zap.Error(err))
		}
	}()
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go gracefulShutdown(server, wg, logger)
	wg.Wait()
	logger.Info("Graceful shutdown complete.")
}

func gracefulShutdown(
	server *web.Server,
	wg *sync.WaitGroup,
	logger *common.Logger,
) {
	defer wg.Done()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()
	<-ctx.Done()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown with error", zap.Error(err))
	}
	logger.Info("Server exiting")
}

func build(cfg common.Config, logger *common.Logger, db *sqlx.DB) *web.Server {
	var server = web.NewServer()
	var employeeRepo = employee.NewRepository(db)
	var roleRepo = role.NewRepository(db)
	var vld = validator.New()
	var employeeService = employee.NewService(employeeRepo, vld)
	var employeeController = employee.NewController(server, employeeService, logger)
	employeeController.RegisterRoutes()
	var roleService = role.NewService(roleRepo, vld)
	var roleController = role.NewController(server, roleService, logger)
	roleController.RegisterRoutes()
	var infoController = info.NewController(server, cfg, db, logger)
	infoController.RegisterRoutes()
	return server
}
