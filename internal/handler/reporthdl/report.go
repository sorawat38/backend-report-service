package reporthdl

import (
	"net/http"
	"time"

	apiresponses "github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/handler/apiResponses"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/models"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	reportSrv service.Report
}

func NewHTTPHandler(reportSrv service.Report) HTTPHandler {
	return HTTPHandler{reportSrv: reportSrv}
}

func (hdl *HTTPHandler) GenerateMontlyReport(c echo.Context) error {

	err := hdl.reportSrv.GenerateMontlyReport(time.Now())
	if err != nil {
		logger.Error("can not generte monthly report", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, apiresponses.InternalError(err))
	}

	var res models.GenerateMonthlyReportResponse
	res.CommonResponse = apiresponses.SuccessResponse()

	return c.JSON(http.StatusOK, res)
}
