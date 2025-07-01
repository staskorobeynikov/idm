package web

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type Server struct {
	App           *fiber.App
	GroupApiV1    fiber.Router
	GroupInternal fiber.Router
}

func registerMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(requestid.New())
}

func NewServer() *Server {
	app := fiber.New()
	registerMiddleware(app)
	groupInternal := app.Group("/internal")
	groupApi := app.Group("/api")
	groupApiV1 := groupApi.Group("/v1")
	return &Server{
		App:           app,
		GroupApiV1:    groupApiV1,
		GroupInternal: groupInternal,
	}
}
