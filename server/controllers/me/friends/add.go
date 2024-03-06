package friends

import (
	"net/http"
	"solution/models"

	"github.com/labstack/echo/v4"
)

func (c *FriendsController) addFriend(ctx echo.Context) error {
	user_id := ctx.Get("UserID").(string)

	var req *struct {
		Login string `json:"login"`
	}
	if err := ctx.Bind(&req); err != nil {
		c.logger.Error("me.addFriend: failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: err.Error()})
	}

	currlogin, err := c.db.GetUserLoginByID(user_id)
	if err != nil {
		c.logger.Error("me.addFriend: postgres error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrorResponse{Err: "postgres error"})
	}

	if req.Login == currlogin {
		return ctx.JSON(http.StatusOK, "ok")
	}
	if err := c.db.WriteFriend(user_id, req.Login); err != nil {
		c.logger.Error("me.addFriend: failed to add friend: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrorResponse{Err: err.Error()})
	}

	return ctx.JSON(http.StatusOK, models.StatusOKResponse{Status: "ok"})
}
