package gateway

import (
	"time"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/models"
)

type PaymentService interface {
	GetOrderByDateMonth(date time.Time) (models.GetOrderByDateMonthResponse, error)
	GetCartById(cartId string) (models.GetCartByIdResponse, error)
}
