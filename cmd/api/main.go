package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/acu4git/gimme-scholarship/internal/app/api/handler"
	"github.com/acu4git/gimme-scholarship/internal/infra/mailer"
	"github.com/acu4git/gimme-scholarship/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port = flag.Int("port", 8080, "port number")
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	// CORS
	origins := []string{"*"}
	if os.Getenv("ALLOW_ORIGINS") != "" {
		origins = strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
	}))

	// jwt auth
	skipPaths := []string{
		"/health",
		"/scholarships",
		"/webhooks/clerk/users",
	}
	optionalPaths := []string{
		"/scholarships",
	}
	auth := handler.NewAuth(skipPaths, optionalPaths)
	e.Use(auth.ClerkJWTMiddleware())
	e.Use(auth.OptionalJWTMiddleware()) // optional認証

	// Injection
	ctx := context.Background()
	repository, err := service.CreateRepository()
	if err != nil {
		log.Fatal(err)
	}
	mailer, err := mailer.NewSESMailer(ctx, "no-reply@kit-gimme-scholarship.com")

	handler := handler.NewAPIHandler(repository, mailer)

	registerRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}

func registerRoutes(router *echo.Echo, handler *handler.APIHandler) {
	router.GET("/scholarships", handler.GetScholarships)
	router.POST("/scholarships/:id/favorite", handler.PostFavoriteScholarship)
	router.DELETE("/scholarships/:id/favorite", handler.DeleteFavoriteScholarship)

	router.POST("/users", handler.PostUser)
	router.PUT("/users", handler.PutUser)

	router.POST("/webhooks/clerk/users", handler.HandleClerkWebhook)

	router.GET("/health", handler.HealthCheck)
}
