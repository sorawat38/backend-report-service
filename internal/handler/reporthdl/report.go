package reporthdl

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
}

func NewHTTPHandler() HTTPHandler {
	return HTTPHandler{}
}

func (hdl *HTTPHandler) GenerateReport(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}
