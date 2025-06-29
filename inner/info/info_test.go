package info

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"idm/inner/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInfo(t *testing.T) {
	var a = assert.New(t)
	t.Run("get info OK response", func(t *testing.T) {
		cfg := common.Config{
			AppName:    "test",
			AppVersion: "1.0.0",
		}
		app := fiber.New()
		server := &web.Server{
			App:           app,
			GroupInternal: app.Group("/internal"),
		}
		controller := NewController(server, cfg, nil)
		controller.db = sqlx.NewDb(nil, "sqlmock")
		controller.RegisterRoutes()
		request := httptest.NewRequest(http.MethodGet, "/internal/info", nil)
		response, err := app.Test(request)
		a.Nil(err)
		a.Equal(http.StatusOK, response.StatusCode)
	})
}

func TestGetHealth(t *testing.T) {
	var a = assert.New(t)
	t.Run("get health OK response", func(t *testing.T) {
		db, _, err := sqlmock.New()
		a.Nil(err)
		defer func() {
			_ = db.Close()
		}()
		cfg := common.Config{}
		app := fiber.New()
		server := &web.Server{
			App:           app,
			GroupInternal: app.Group("/internal"),
		}
		controller := NewController(server, cfg, nil)
		controller.db = sqlx.NewDb(db, "sqlmock")
		controller.RegisterRoutes()
		request := httptest.NewRequest(http.MethodGet, "/internal/health", nil)
		response, err := app.Test(request)
		a.Nil(err)
		a.Equal(http.StatusOK, response.StatusCode)
	})
	t.Run("get health OK response", func(t *testing.T) {
		db, _, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		a.Nil(err)
		defer func() {
			_ = db.Close()
		}()
		cfg := common.Config{}
		app := fiber.New()
		server := &web.Server{
			App:           app,
			GroupInternal: app.Group("/internal"),
		}
		controller := NewController(server, cfg, nil)
		controller.db = sqlx.NewDb(db, "sqlmock")
		controller.RegisterRoutes()
		request := httptest.NewRequest(http.MethodGet, "/internal/health", nil)
		response, err := app.Test(request)
		a.Nil(err)
		a.Equal(http.StatusInternalServerError, response.StatusCode)
	})
}
