package Middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"navigation/pkg/Clients/Redis"
	"net/http"
	"time"
)

var ctx = context.Background()

func SessionCheckMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("api-key")
			if token == "" {
				return c.String(http.StatusForbidden, "Missing API Key")
			}

			exists, err := Redis.Crowdsourcing.RedisClient.Exists(ctx, token).Result()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Redis error")
			}
			if exists == 0 {
				return c.String(http.StatusForbidden, "Invalid API Key")
			}

			key := fmt.Sprintf("%s_nav_%d", token, time.Now().UnixMilli())
			err = Redis.Crowdsourcing.RedisClient.Set(ctx, key, 1, 0).Err()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Redis write error")
			}

			return next(c)
		}
	}
}