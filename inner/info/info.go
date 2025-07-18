package info

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
	"idm/inner/web"
	"time"
)

type Controller struct {
	server *web.Server
	cfg    common.Config
	db     *sqlx.DB
	logger *common.Logger
}

func NewController(
	server *web.Server,
	cfg common.Config,
	db *sqlx.DB,
	logger *common.Logger,
) *Controller {
	return &Controller{
		server: server,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

type Response struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupInternal.Get("/info", c.GetInfo)
	c.server.GroupInternal.Get("/health", c.GetHealth)
}

func (c *Controller) GetInfo(ctx *fiber.Ctx) error {
	var err = ctx.Status(fiber.StatusOK).JSON(&Response{
		Name:    c.cfg.AppName,
		Version: c.cfg.AppVersion,
	})
	if err != nil {
		c.logger.Error("error returning info")
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning info")
	}
	return nil
}

func (c *Controller) GetHealth(ctx *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.db.PingContext(ctxWithTimeout); err != nil {
		c.logger.Error("error returning health")
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	if err := ctx.Status(fiber.StatusOK).SendString("OK"); err != nil {
		c.logger.Error("error returning health")
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
