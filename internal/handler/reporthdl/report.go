package reporthdl

import (
	"net/http"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/service"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	reportSrv service.Report
}

func NewHTTPHandler(reportSrv service.Report) HTTPHandler {
	return HTTPHandler{reportSrv: reportSrv}
}

func (hdl *HTTPHandler) GenerateMontlyReport(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}
