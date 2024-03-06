package friends

import (
	"log/slog"
	"solution/pkg/database/postgres"
	"solution/server/controllers"
	"solution/server/middleware"

	"github.com/labstack/echo/v4"
)

type FriendsController struct {
	db     *postgres.PostgresDB
	logger *slog.Logger
}

func NewFriendsController(db *postgres.PostgresDB, logger *slog.Logger) *FriendsController {
	return &FriendsController{
		db:     db,
		logger: logger,
	}
}

func (c *FriendsController) GetGroup() string {
	return "api/friends"
}

func (c *FriendsController) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.AuthMiddleware(c.db),
	}
}

func (c *FriendsController) GetHandlers() []controllers.ControllerHandler {
	return []controllers.ControllerHandler{
		&controllers.Handler{
			Method:  "POST",
			Path:    "/add",
			Handler: c.addFriend,
		},
		&controllers.Handler{
			Method:  "POST",
			Path:    "/remove",
			Handler: c.removeFriend,
		},
		&controllers.Handler{
			Method:  "GET",
			Path:    "",
			Handler: c.getAll,
		},
	}
}
