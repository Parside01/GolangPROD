package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"solution/models"
	"solution/pkg/database/postgres"
	validatate "solution/pkg/validate"
	"solution/server/controllers"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/twharmon/gouid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrBadRequest = errors.New("empty or uncorect request body")
	ErrUserExist  = errors.New("user exist")
)

type RegisterController struct {
	logger *slog.Logger
	db     *postgres.PostgresDB
}

type Response struct {
	Err    string `json:"error"`
	Status int    `json:"status"`
	Ok     bool   `json:"ok"`
}

func NewRegisterController(logger *slog.Logger, db *postgres.PostgresDB) *RegisterController {
	return &RegisterController{
		logger: logger,
		db:     db,
	}
}

func (c *RegisterController) GetGroup() string {
	return "api/auth"
}

func (c *RegisterController) GetMiddleware() []echo.MiddlewareFunc {
	return nil
}

func (c *RegisterController) GetHandlers() []controllers.ControllerHandler {
	return []controllers.ControllerHandler{
		&controllers.Handler{
			Method:  "POST",
			Path:    "/register",
			Handler: c.registerUser,
		},
	}
}
func (c *RegisterController) registerUser(e echo.Context) error {
	var req *struct {
		Login       string `json:"login"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		CountryCode string `json:"countryCode"`
		Phone       string `json:"phone"`
		IsPublic    bool   `json:"isPublic"`
		Image       string `json:"image"`
	}
	if err := e.Bind(&req); err != nil {
		c.logger.Error("RegisterController.registerUser: failed to bind request: %v", err)
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: ErrBadRequest.Error()})
	}

	if req == nil {
		c.logger.Error("RegisterController.registerUser: uncorrect user request")
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: "uncorrect user request"})
	}

	if ok := validatate.IsValidPassword(req.Password); !ok {
		c.logger.Error("RegisterController.registerUser: invalid password")
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: "invalid password"})
	}

	encr, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.logger.Error("RegisterController.registerUser: failed to hash password: %v", err)
		return e.JSON(500, models.ErrorResponse{Err: ErrBadRequest.Error()})
	}

	userid := gouid.Bytes(16)
	user := &models.User{
		ID:          userid.String(),
		Login:       req.Login,
		Phone:       req.Phone,
		CountryCode: req.CountryCode,
		IsPublic:    req.IsPublic,
		Email:       req.Email,
		Password:    string(encr),
		CreatedAt:   time.Now(),
		Image:       req.Image,
	}

	if ok, _ := c.db.CountryIsExist(req.CountryCode); !ok {
		c.logger.Error("RegisterController.registerUser: invalid country code")
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: "invalid country code"})
	}

	if err := validatate.IsValidUser(user); err != nil {
		c.logger.Error("RegisterController.registerUser: failed to validate request", err)
		return e.JSON(http.StatusBadRequest, models.ErrorResponse{Err: ErrBadRequest.Error()})
	}

	if err = c.db.AddUser(user); err != nil {
		if err == postgres.ErrUniqueViolation {
			c.logger.Error("RegisterController.registerUser: user already exist")
			return e.JSON(http.StatusConflict, models.ErrorResponse{Err: ErrUserExist.Error()})
		}
		c.logger.Error("RegisterController.registerUser: failed to add user: %v", err)
		return e.JSON(500, models.ErrorResponse{Err: "server error"})
	}
	// обожаю костыли)
	var res struct {
		Profile *models.UserProfile `json:"profile"`
	}
	res.Profile = user.GetProfile()
	return e.JSON(http.StatusCreated, res)
}
