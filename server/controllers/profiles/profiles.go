package profiles

import (
	"log/slog"
	"net/http"
	"solution/models"
	"solution/pkg/database/postgres"
	"solution/server/controllers"
	"solution/server/middleware"

	"github.com/labstack/echo/v4"
)

type ProfilesController struct {
	logger *slog.Logger
	db     *postgres.PostgresDB
}

func NewProfilesController(logger *slog.Logger, db *postgres.PostgresDB) *ProfilesController {
	return &ProfilesController{
		logger: logger,
		db:     db,
	}
}

func (s *ProfilesController) GetGroup() string {
	return "api/profiles"
}

func (s *ProfilesController) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.AuthMiddleware(s.db),
	}
}

func (s *ProfilesController) GetHandlers() []controllers.ControllerHandler {
	return []controllers.ControllerHandler{
		&controllers.Handler{
			Method:  "GET",
			Path:    "/:login",
			Handler: s.getLoginProfile,
		},
	}
}

func (s *ProfilesController) getLoginProfile(e echo.Context) error {
	user_id := e.Get("UserID").(string)
	login := e.Param("login")

	user, err := s.db.GetUserByLoginIfPublic(user_id, login)
	if err != nil {
		s.logger.Error("me.getProfiles: failed to get user profile: %v", err)
		return echo.NewHTTPError(http.StatusForbidden, models.ErrorResponse{Err: err.Error()})
	}
	return e.JSON(http.StatusOK, user.GetProfile())
}
