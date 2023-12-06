package reporthdl

import (
	"net/http"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	paymentGw gateway.PaymentService
}

func NewHTTPHandler(paymentGw gateway.PaymentService) HTTPHandler {
	return HTTPHandler{paymentGw: paymentGw}
}

func (hdl *HTTPHandler) GenerateReport(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}
