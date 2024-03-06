package middleware

import (
	"net/http"
	"solution/models"
	"solution/pkg/database/postgres"
	"solution/pkg/token"

	"github.com/labstack/echo/v4"
)

// Как по мне выглядит довольно удобно, просто закидываем id пользователя в конекстс и потом в контроллерах его просто получаем по ключу.
func AuthMiddleware(p *postgres.PostgresDB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := c.Request().Header.Get("Authorization")
			if t == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: "empty token"})
			}
			t, err := token.GetAuthTokenFromBearerToken(t)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: err.Error()})
			}
			info, ok := token.VerifyToken("secret", t)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: "invalid token"})
			}
			if err := p.CanUseToken(info.ID, t); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: err.Error()})
			}

			c.Set("UserID", info.ID)
			return next(c)
		}
	}
}
