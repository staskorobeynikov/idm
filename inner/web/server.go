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

func registerMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(requestid.New())
}

func NewServer(enableSwagger bool) *Server {
	app := fiber.New()
	registerMiddleware(app)
	if enableSwagger {
		go startSwaggerServer()
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:8081/"}, // Swagger UI порт
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

func startSwaggerServer() {
	app := fiberV2.New()
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	err := app.Listen(":8081")
	if err != nil {
		panic(err)
	}
}
