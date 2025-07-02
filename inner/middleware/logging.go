package middleware

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"idm/inner/common"
)

func LoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := context.WithValue(c.Context(), common.LoggerKey, logger)
		c.SetContext(ctx)
		return c.Next()
	}
}

func GetLogger(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(common.LoggerKey).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}
