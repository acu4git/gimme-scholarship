package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
}
