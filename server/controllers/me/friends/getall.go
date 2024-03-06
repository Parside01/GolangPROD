package friends

import (
	"net/http"
	"solution/models"
	validatate "solution/pkg/validate"

	"github.com/labstack/echo/v4"
)

func (f *FriendsController) getAll(e echo.Context) error {
	user_id := e.Get("UserID").(string)

	//	не очень понятно что делать если некорректные параметры пагинации. Но вернем невный запос.
	limit, offset := e.QueryParam("limit"), e.QueryParam("offset")
	if limit == "" {
		limit = "5"
	}
	if offset == "" {
		offset = "0"
	}
	l, o, ok := validatate.IsValidPaginationParams(limit, offset)
	if !ok {
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: "invalid pagination params"})
	}

	fs, err := f.db.GetUserFriendByLimit(user_id, l, o)
	if err != nil {
		f.logger.Error("me.getAll: %v", err)
		return e.JSON(http.StatusUnauthorized, struct {
			Err string `json:"err"`
		}{Err: err.Error()})
	}

	return e.JSON(http.StatusOK, models.GetFriendResult(fs))
}
