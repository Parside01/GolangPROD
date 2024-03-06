package me

import (
	"log/slog"
	"net/http"
	"solution/models"
	"solution/pkg/database/postgres"
	validatate "solution/pkg/validate"
	"solution/server/controllers"
	"solution/server/middleware"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordController struct {
	logger *slog.Logger
	db     *postgres.PostgresDB
}

func NewUpdatePasswordController(logger *slog.Logger, db *postgres.PostgresDB) *UpdatePasswordController {
	return &UpdatePasswordController{
		logger: logger,
		db:     db,
	}
}

func (s *UpdatePasswordController) GetGroup() string {
	return "api/me"
}

func (s *UpdatePasswordController) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.AuthMiddleware(s.db),
	}
}

func (s *UpdatePasswordController) GetHandlers() []controllers.ControllerHandler {
	return []controllers.ControllerHandler{
		&controllers.Handler{
			Method:  "POST",
			Path:    "/updatePassword",
			Handler: s.updateMyPassword,
		},
	}
}

func (s *UpdatePasswordController) updateMyPassword(e echo.Context) error {
	id := e.Get("UserID").(string)

	req := new(struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	})
	if err := e.Bind(&req); err != nil {
		s.logger.Error("me.updateMyPassword: failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: err.Error()})
	}

	if !validatate.IsValidPassword(req.NewPassword) {
		s.logger.Error("me.updateMyPassword: uncorrect password")
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrorResponse{Err: "uncorrect password"})
	}

	user, err := s.db.GetUserByID(id)
	if err != nil {
		s.logger.Error("me.updateMyPassword: failed to get user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrorResponse{Err: err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		s.logger.Error("me.updateMyPassword: invalid old password")
		return echo.NewHTTPError(http.StatusForbidden, models.ErrorResponse{Err: err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("me.updateMyPassword: failed to hash password: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrorResponse{Err: err.Error()})
	}

	user.Password = string(hash)
	if err := s.db.UpdateUser(user); err != nil {
		s.logger.Error("me.updateMyPassword: failed to update user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrorResponse{Err: err.Error()})
	}

	err = s.db.DeleteTokenByUserID(user.ID)
	if err != nil {
		s.logger.Error("me.updateMyPassword: failed to delete token: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrorResponse{Err: err.Error()})
	}

	return e.JSON(http.StatusOK, models.StatusOKResponse{Status: "ok"})
}
