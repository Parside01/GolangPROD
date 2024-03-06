package friends

import (
	"net/http"
	"solution/models"

	"github.com/labstack/echo/v4"
)

func (c *FriendsController) removeFriend(ctx echo.Context) error {
	user_id := ctx.Get("UserID").(string)
	var req *struct {
		Login string `json:"login"`
	}
	if err := ctx.Bind(&req); err != nil {
		c.logger.Error("me.removeFriend: failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: err.Error()})
	}

	if err := c.db.DeleteFriendByLogin(req.Login, user_id); err != nil {
		c.logger.Error("me.removeFriend: postgres error: %v", err)
		return ctx.JSON(http.StatusOK, models.StatusOKResponse{Status: "ok"})
	}
	return ctx.JSON(http.StatusOK, models.StatusOKResponse{Status: "ok"})
}
