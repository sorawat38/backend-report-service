package handler

import (
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/handler/reporthdl"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, reportHandler reporthdl.HTTPHandler) {

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.RequestID(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"}, // Replace with your frontend origin(s)
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		}),
	)

	report := e.Group("/report")
	report.POST("/generate", reportHandler.GenerateMontlyReport)

}
