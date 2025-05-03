package handler

import (
	"net/http"

	"github.com/acu4git/gimme-scholarship/internal/domain/repository"
	"github.com/labstack/echo/v4"
)

type APIHandler struct {
	repository repository.Repository
}

func NewAPIHandler(repository repository.Repository) *APIHandler {
	return &APIHandler{
		repository: repository,
	}
}

func (h *APIHandler) GetScholarships(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
