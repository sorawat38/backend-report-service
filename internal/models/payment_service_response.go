package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type GetOrderByDateMonthResponse struct {
	CommonResponse
	Data []GetOrderByDateMonthResponseBody `json:"data,omitempty"`
}

type GetOrderByDateMonthResponseBody struct {
	OrderId      string          `json:"order_id"`
	CustomerId   string          `json:"customer_id"`
	SubTotal     decimal.Decimal `json:"sub_total"`
	DiscountCode string          `json:"discount_code,omitempty"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	DateTime     time.Time       `json:"date_time"`
	Status       string          `json:"status"`
}
