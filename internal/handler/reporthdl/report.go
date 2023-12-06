package reporthdl

import (
	"errors"
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

	dateStr := c.Param("date")
	if dateStr == "" {
		err := errors.New("date is empty")
		logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, apiresponses.InvalidInputError(err))
	}

	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		logger.Error("can't parse date to format YYYY-MM-DD", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, apiresponses.InvalidInputError(err))
	}

	err = hdl.reportSrv.GenerateMontlyReport(date)
	if err != nil {
		logger.Error("can not generte monthly report", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, apiresponses.InternalError(err))
	}

	var res models.GenerateMonthlyReportResponse
	res.CommonResponse = apiresponses.SuccessResponse()

	return c.JSON(http.StatusOK, res)
}
