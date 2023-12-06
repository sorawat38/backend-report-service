package reportsrv

import (
	"time"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/jung-kurt/gofpdf"
	"go.uber.org/zap"
)

type service struct {
	paymentGw gateway.PaymentService
}

func New(paymentGw gateway.PaymentService) service {
	return service{paymentGw: paymentGw}
}

func (srv service) GenerateMontlyReport(date time.Time) error {

	// Get orders from payment service
	// _, err := srv.paymentGw.GetOrderByDateMonth(date)
	// if err != nil {
	// 	logger.Error("can't get orders by date month", zap.Error(err))
	// 	return err
	// }

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, World!")

	if err := pdf.OutputFileAndClose("hello.pdf"); err != nil {
		logger.Error("can't generate report", zap.Error(err))
		return err
	}

	return nil
}
