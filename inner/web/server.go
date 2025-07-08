package web

import (
	fiberV2 "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/swaggo/fiber-swagger"
	_ "idm/docs"
)

type Server struct {
	App           *fiber.App
	GroupApiV1    fiber.Router
	GroupInternal fiber.Router
}

type ServerV2 struct {
	App *fiberV2.App
}

func registerMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(requestid.New())
}

func NewServer() *Server {
	app := fiber.New()
	registerMiddleware(app)
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://localhost:8081",
		},
	}))
	groupInternal := app.Group("/internal")
	groupApi := app.Group("/api")
	groupApiV1 := groupApi.Group("/v1")
	return &Server{
		App:           app,
		GroupApiV1:    groupApiV1,
		GroupInternal: groupInternal,
	}
}

func NewServerV2() *ServerV2 {
	app := fiberV2.New()
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	return &ServerV2{App: app}
}
