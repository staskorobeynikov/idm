package main

import (
	"context"
	"crypto/tls"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"idm/docs"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/middleware"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// @title IDM API documentation
// @description  API for managing IDM service
// @host localhost:8080
// @BasePath /api/v1
// @schemes https
// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl http://localhost:9990/realms/idm/protocol/openid-connect/token
// @scope openid Access everything
func main() {
	cfg := common.GetConfig(".env")
	docs.SwaggerInfo.Version = cfg.AppVersion
	var logger = common.NewLogger(cfg)
	cer, err := tls.LoadX509KeyPair(cfg.SslSert, cfg.SslKey)
	if err != nil {
		logger.Panic("failed certificate loading: %s", zap.Error(err))
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
	defer func() { _ = logger.Sync() }()
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error closing db", zap.Error(err))
		}
	}()
	var server = build(cfg, logger, db)
	go func() {
		ln, err := tls.Listen("tcp", ":8080", tlsConfig)
		if err != nil {
			logger.Panic("failed TLS listener creating: %s", zap.Error(err))
		}
		err = server.App.Listener(ln)
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

func build(
	cfg common.Config,
	logger *common.Logger,
	db *sqlx.DB,
) *web.Server {
	var server = web.NewServer()
	server.App.Use(requestid.New())
	server.App.Use(middleware.LoggerMiddleware(logger))
	server.App.Use(recover.New())
	server.App.Use("/swagger/*", swagger.HandlerDefault)
	server.GroupApiV1.Use(web.AuthMiddleware(logger))
	var employeeRepo = employee.NewRepository(db)
	var roleRepo = role.NewRepository(db)
	var vld = validator.New()
	var employeeService = employee.NewService(employeeRepo, vld)
	var employeeController = employee.NewController(server, employeeService)
	employeeController.RegisterRoutes()
	var roleService = role.NewService(roleRepo, vld)
	var roleController = role.NewController(server, roleService, logger)
	roleController.RegisterRoutes()
	var infoController = info.NewController(server, cfg, db, logger)
	infoController.RegisterRoutes()
	return server
}
