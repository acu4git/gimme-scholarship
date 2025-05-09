package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/acu4git/gimme-scholarship/internal/app/api/handler"
	"github.com/acu4git/gimme-scholarship/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port = flag.Int("port", 3000, "port number")
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	repository, err := service.CreateRepository()
	if err != nil {
		log.Fatal(err)
	}
	handler := handler.NewAPIHandler(repository)

	e.GET("/scholarships", handler.GetScholarships)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
