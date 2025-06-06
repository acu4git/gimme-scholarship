package handler

import (
	"fmt"
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

func (h *APIHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Nice deeeeesu",
	})
}

func (h *APIHandler) PostUser(c echo.Context) error {
	param := PostUserInput{}
	if err := c.Bind(&param); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": "failed to bind PostUserInput: " + err.Error(),
		})
	}

	userID, ok := c.Get(userIDKey).(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": fmt.Sprintf("failed to c.Get() %v", c.Get(userIDKey)),
		})
	}

	if err := h.repository.CreateUser(repository.UserInput{ID: userID, Email: param.Email, Level: param.Level, Grade: int64(param.Grade), AcceptEmail: param.AcceptEmail}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "failed to create user: " + err.Error(),
		})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *APIHandler) GetScholarships(c echo.Context) error {
	var pUserID *string
	userID, ok := c.Get(userIDKey).(string)
	if !ok || userID == "" {
		pUserID = nil
	} else {
		pUserID = &userID
	}

	param := GetScholarshipInput{}
	if err := c.Bind(&param); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}

	scholarships, err := h.repository.GetScholarships(repository.FilterOption{UserID: pUserID, Target: param.Target, Type: param.Type})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return c.JSONPretty(http.StatusOK, toGetScholarshipsOutput(scholarships), prettyIndent)
}
