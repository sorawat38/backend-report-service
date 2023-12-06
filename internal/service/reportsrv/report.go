package reportsrv

import (
	"time"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway"
)

type service struct {
	paymentGw gateway.PaymentService
}

func New(paymentGw gateway.PaymentService) service {
	return service{paymentGw: paymentGw}
}

func (srv service) GenerateMontlyReport(date time.Time) error {
	return nil
}
