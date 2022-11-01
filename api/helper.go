package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type JsonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type CustomValidator struct {
	validator *validator.Validate
}

// Custom validator for echo
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
