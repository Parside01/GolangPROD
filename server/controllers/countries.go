package controllers

import (
	"log/slog"
	"net/http"
	"solution/pkg/database/postgres"
	validatate "solution/pkg/validate"

	"github.com/labstack/echo/v4"
)

type CountriesController struct {
	logger *slog.Logger
	db     *postgres.PostgresDB
}

func NewCountriesController(logger *slog.Logger, db *postgres.PostgresDB) *CountriesController {
	return &CountriesController{
		logger: logger,
		db:     db,
	}
}

func (c *CountriesController) GetGroup() string {
	return "/api/countries"
}

func (c *CountriesController) GetMiddleware() []echo.MiddlewareFunc {
	return nil
}

func (c *CountriesController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "GET",
			Path:    "/:alpha2",
			Handler: c.getCountries,
		},
		&Handler{
			Method:  "GET",
			Path:    "",
			Handler: c.getAllCountries,
		},
	}
}

func (c *CountriesController) getAllCountries(e echo.Context) error {
	e.Set("Content-Type", "application/json")
	region := e.QueryParam("region")

	if !validatate.IsValidRegion(region) {
		c.logger.Error("CountriesController.getAllCountries", "error", "invalid region")
		return e.JSON(http.StatusBadRequest, "invalid region")
	}

	country, err := c.db.GetCountriesByRegion(region)
	if err != nil {
		c.logger.Error("CountriesController.getAllCountries", "error", err)
		return e.JSON(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, country)
}

func (c *CountriesController) getCountries(e echo.Context) error {
	e.Set("Content-Type", "application/json")
	filter := e.Param("alpha2")

	if !validatate.IsValidAlpha2(filter) {
		c.logger.Error("CountriesController.getCountries", "error", "invalid alpha2 code")
		return e.JSON(http.StatusBadRequest, "invalid alpha2 code")
	}

	country, err := c.db.GetCountries(filter)
	if err != nil {
		c.logger.Error("CountriesController.getCountries", "error", err)
		return e.JSON(http.StatusNoContent, err)
	}

	return e.JSON(http.StatusOK, country)
}
