package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"solution/models"
	"solution/pkg/database/postgres"
	"solution/pkg/token"
	validatate "solution/pkg/validate"
	"solution/server/controllers"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type SignInController struct {
	logger   *slog.Logger
	postgres *postgres.PostgresDB
}

func NewSignInController(logger *slog.Logger, postgres *postgres.PostgresDB) *SignInController {
	return &SignInController{
		logger:   logger,
		postgres: postgres,
	}
}

func (s *SignInController) GetGroup() string {
	return "api/auth"
}

func (s *SignInController) GetMiddleware() []echo.MiddlewareFunc {
	return nil
}

func (s *SignInController) GetHandlers() []controllers.ControllerHandler {
	return []controllers.ControllerHandler{
		&controllers.Handler{
			Method:  "POST",
			Path:    "/sign-in",
			Handler: s.signInUser,
		},
	}
}

/*
*	Сначала хотел сделать через систему сессий, но получил разочарование в виде невозможности нормального использование редиски.
*	Попытка реализации через постгрес была такой себе.
*	Поэтому будет просто токен.
 */
func (s *SignInController) signInUser(e echo.Context) error {
	var req *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := e.Bind(&req); err != nil {
		s.logger.Error("SignInController.signInUser: failed to bind request: %v", err)
		return e.JSON(http.StatusBadRequest, &Response{
			Err:    err.Error(),
			Status: 400,
			Ok:     false,
		})
	}
	if req == nil {
		s.logger.Error("SignInController.signInUser: failed to check request: %v", errors.New("wkbmwb"))
		return e.JSON(http.StatusBadRequest, &Response{
			Err:    "wbwbwb",
			Ok:     false,
			Status: http.StatusBadRequest,
		})
	}

	user, err := s.postgres.GetUserByLogin(req.Login)
	if err != nil {
		s.logger.Error("SignInController.signInUser: failed to get user: %v", err)
		return e.JSON(http.StatusUnauthorized, &Response{
			Err:    err.Error(),
			Status: http.StatusUnauthorized,
			Ok:     false,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Error("SignInController.signInUser: failed to compare password: %v", err)
		return e.JSON(http.StatusUnauthorized, &Response{
			Err:    err.Error(),
			Status: http.StatusUnauthorized,
			Ok:     false,
		})
	}

	if err := validatate.ChecStructForNil(user); err != nil {
		s.logger.Error("SignInController.signInUser: failed to check user: %v", err)
		return e.JSON(http.StatusInternalServerError, &Response{
			Err:    err.Error(),
			Ok:     false,
			Status: http.StatusInternalServerError,
		})
	}

	access, err := token.NewAccessToken("secret", user.ID)
	if err != nil {
		s.logger.Error("SignInController.signInUser: failed to create access token: %v", err)
		return e.JSON(http.StatusInternalServerError, &Response{
			Err:    err.Error(),
			Status: http.StatusUnauthorized,
			Ok:     false,
		})
	}
	err = s.postgres.WriteToken(&models.Token{
		Token:  access,
		UserID: user.ID,
	})
	if err != nil {
		s.logger.Error("SignInController.signInUser: failed to write access token: %v", err)
		return e.JSON(http.StatusInternalServerError, &Response{
			Err:    err.Error(),
			Status: http.StatusUnauthorized,
			Ok:     false,
		})
	}

	var resp struct {
		AccessToken string `json:"token"`
	}
	resp.AccessToken = access
	return e.JSON(http.StatusOK, resp)
}
