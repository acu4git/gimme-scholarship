package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/acu4git/gimme-scholarship/internal/domain/repository"
	"github.com/labstack/echo/v4"
	svix "github.com/svix/svix-webhooks/go"
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

func (h *APIHandler) PutUser(c echo.Context) error {
	param := PutUserInput{}
	if err := c.Bind(&param); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": "failed to bind PutUserInput: " + err.Error(),
		})
	}

	userID, ok := c.Get(userIDKey).(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": fmt.Sprintf("failed to c.Get() %v", c.Get(userIDKey)),
		})
	}

	if err := h.repository.UpdateUser(repository.UserInput{ID: userID, Level: param.Level, Grade: int64(param.Grade), AcceptEmail: param.AcceptEmail}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "failed to update user: " + err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
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

	scholarships, favoriteMap, err := h.repository.FindScholarships(repository.FilterOption{UserID: pUserID, Target: param.Target, Type: param.Type})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return c.JSONPretty(http.StatusOK, toGetScholarshipsOutput(scholarships, favoriteMap), prettyIndent)
}

func (h *APIHandler) PostFavoriteScholarship(c echo.Context) error {
	_id := c.Param("id")
	scholarshipID, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": fmt.Sprintf("failed to parse scholarship id (%s)", _id),
		})
	}

	userID, ok := c.Get(userIDKey).(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": fmt.Sprintf("failed to c.Get() %v", c.Get(userIDKey)),
		})
	}

	input := repository.UserFavoriteInput{
		Mode:          "REGISTER",
		UserID:        userID,
		ScholarshipID: scholarshipID,
	}

	if err := h.repository.UserFavoriteAction(input); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to register favorite scholarship",
			"error":   err,
		})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *APIHandler) DeleteFavoriteScholarship(c echo.Context) error {
	_id := c.Param("id")
	scholarshipID, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": fmt.Sprintf("failed to parse scholarship id (%s)", _id),
		})
	}

	userID, ok := c.Get(userIDKey).(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": fmt.Sprintf("failed to c.Get() %v", c.Get(userIDKey)),
		})
	}

	input := repository.UserFavoriteInput{
		Mode:          "DELETE",
		UserID:        userID,
		ScholarshipID: scholarshipID,
	}

	if err := h.repository.UserFavoriteAction(input); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to delete favorite scholarship",
			"error":   err,
		})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *APIHandler) HandleClerkWebhook(c echo.Context) error {
	// Webhook secretの取得
	secret := os.Getenv("CLERK_WEBHOOK_SECRET_KEY")
	if secret == "" {
		log.Println("CLERK_WEBHOOK_SECRET_KEY is not set")
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "CLERK_WEBHOOK_SECRET_KEY is not set",
		})
	}

	// Webhookインスタンス作成
	wh, err := svix.NewWebhook(secret)
	if err != nil {
		log.Println("failed to create webhook instance:", err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	// リクエストヘッダーとボディの取得
	headers := c.Request().Header
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Println("failed to read body")
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	// 検証
	if err := wh.Verify(payload, headers); err == nil {
		log.Println("failed to verify webhook:", err)
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}

	// パース処理
	var event ClerkWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}
	var clerkUser ClerkUser
	json.Unmarshal(event.Data, &clerkUser)

	// イベント毎の処理
	switch event.Type {
	case "user.created":
		log.Println("event type: user.created")
	case "user.deleted":
		log.Println("event type: user.deleted")
		if err := h.repository.DeleteUser(clerkUser.ID); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"error": err.Error(),
			})
		}
	case "user.updated":
		log.Println("event type: user.updated")
	default:
		log.Println("invalid event type")
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid event type",
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
