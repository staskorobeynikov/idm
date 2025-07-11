package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"idm/inner/common"
)

func LoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.Context(), common.LoggerKey, logger)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func GetLogger(ctx *fiber.Ctx) *zap.Logger {
	if l, ok := ctx.UserContext().Value(common.LoggerKey).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}
