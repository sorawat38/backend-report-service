package main

import (
	"net/http"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/config"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/handler"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/labstack/echo/v4"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger.InitLog(cfg.Log)
	defer logger.CloseLogger()

	// Init repository

	// Starting server
	e := echo.New()
	handler.InitRoute(e)

	logger.Infof("Starting server on port %v...\n", cfg.App.Port)
	if err := e.Start(":" + cfg.App.Port); err != http.ErrServerClosed {
		logger.Fatal(err.Error())
	}
}
