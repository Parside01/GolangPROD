package me

import (
	"log/slog"
	"net/http"
	"solution/models"
	"solution/pkg/database/postgres"
	validatate "solution/pkg/validate"
	"solution/server/controllers"
	"solution/server/middleware"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProfileController struct {
	logger *slog.Logger
	db     *postgres.PostgresDB
}

func NewProfileController(logger *slog.Logger, db *postgres.PostgresDB) *ProfileController {
	return &ProfileController{
		logger: logger,
		db:     db,
	}
}

func (s *ProfileController) GetGroup() string {
	return "api/me"
}

func (s *ProfileController) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.AuthMiddleware(s.db),
	}
}

func (s *ProfileController) GetHandlers() []controllers.ControllerHandler {
	return []controllers.ControllerHandler{
		&controllers.Handler{
			Method:  "GET",
			Path:    "/profile",
			Handler: s.getMyProfile,
		},
		&controllers.Handler{
			Method:  "PATCH",
			Path:    "/profile",
			Handler: s.updateMyProfile,
		},
	}
}

func (s *ProfileController) updateMyProfile(e echo.Context) error {
	user_id := e.Get("UserID").(string)

	req := new(struct {
		CountryCode string `json:"countryCode"`
		Phone       string `json:"phone"`
		IsPublic    bool   `json:"isPublic"`
		Image       string `json:"image"`
	})
	if err := e.Bind(&req); err != nil {
		s.logger.Error("me.getMyProfile: failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: err.Error()})
	}

	user, err := s.db.GetUserByID(user_id)
	if err != nil {
		s.logger.Error("me.getMyProfile: failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: err.Error()})
	}
	ok, _ := s.db.CountryIsExist(req.CountryCode)
	if req.CountryCode != "" && ok {
		user.CountryCode = req.CountryCode
	} else if req.CountryCode != "" {
		s.logger.Error("me.getMyProfile: invalid country code: %v", req.CountryCode)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: "invalid country code"})
	}

	var body []byte
	e.Request().Body.Read(body)
	//	а как это обрабатывать? Я же не могу знать был передан @isPublic или нет ибо если нет то он всегда false, тогда будет замена.
	//	топ самых странных решений, ахах
	if strings.Contains(string(body), "isPublic") {
		user.IsPublic = req.IsPublic
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Image != "" {
		user.Image = req.Image
	}

	if err := validatate.IsValidUser(user); err != nil {
		s.logger.Error("me.getMyProfile: invalid user: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: err.Error()})
	}

	err = s.db.UpdateUser(user)
	if err != nil {
		s.logger.Error("me.getMyProfile: failed to update user: %v", err)
		return echo.NewHTTPError(http.StatusConflict, models.ErrorResponse{Err: err.Error()})
	}
	return e.JSON(http.StatusOK, user.GetProfile())
}

func (s *ProfileController) getMyProfile(e echo.Context) error {
	userid := e.Get("UserID").(string)

	user, err := s.db.GetUserByID(userid)
	if err != nil {
		s.logger.Error("me.getMyProfile: failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, models.ErrorResponse{Err: err.Error()})
	}

	return e.JSON(http.StatusOK, user.GetProfile())
}
