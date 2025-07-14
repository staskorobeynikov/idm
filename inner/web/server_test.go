package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware(t *testing.T) {
	var a = assert.New(t)
	t.Run("Middleware with Panic", func(t *testing.T) {
		app := fiber.New()
		app.Use(recover.New())
		app.Get("/panic", func(c *fiber.Ctx) error {
			panic("something went wrong")
		})
		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		resp, err := app.Test(req)
		a.Nil(err)
		a.Equal(fiber.StatusInternalServerError, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		a.Contains(string(body), "something went wrong")
	})
	t.Run("Middleware with Panic and server alive", func(t *testing.T) {
		app := fiber.New()
		app.Use(recover.New())
		app.Get("/panic", func(c *fiber.Ctx) error {
			panic("simulated panic")
		})
		app.Get("/ok", func(c *fiber.Ctx) error {
			return c.SendString("I am alive")
		})
		req1 := httptest.NewRequest(http.MethodGet, "/panic", nil)
		resp1, err := app.Test(req1)
		a.Nil(err)
		a.Equal(fiber.StatusInternalServerError, resp1.StatusCode)
		req2 := httptest.NewRequest(http.MethodGet, "/ok", nil)
		resp2, err := app.Test(req2)
		a.Nil(err)
		a.Equal(fiber.StatusOK, resp2.StatusCode)
		body, _ := io.ReadAll(resp2.Body)
		a.Equal("I am alive", string(body))
	})
}

func TestRequestMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(requestid.New())
	app.Get("/test", func(c *fiber.Ctx) error {
		id := string(c.Response().Header.Peek(fiber.HeaderXRequestID))
		return c.SendString(id)
	})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	requestID := resp.Header.Get("X-Request-ID")
	assert.NotEmpty(t, requestID)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, requestID, string(body))
}
