package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/acu4git/gimme-scholarship/internal/app/api/handler"
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

	// Injection
	repository, err := service.CreateRepository()
	if err != nil {
		log.Fatal(err)
	}
	handler := handler.NewAPIHandler(repository)

	registerRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}

func registerRoutes(router *echo.Echo, handler *handler.APIHandler) {
	router.GET("/scholarships", handler.GetScholarships)
}
