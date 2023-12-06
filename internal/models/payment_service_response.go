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

type GetCartByDateMonthResponse struct {
	CommonResponse
	Data []GetCartByDateMonthResponseBody `json:"data,omitempty"`
}

type GetCartByDateMonthResponseBody struct {
	CartId            string `json:"cart_id"`
	No                int    `json:"no"`
	CustomerId        string `json:"customer_id"`
	Date              string `json:"date"`
	MenuId            int    `json:"menu_id"`
	Quantity          int    `json:"quantity"`
	Status            string `json:"status"`
	Properties        string `json:"properties,omitempty"`
	AdditionalRequest string `json:"additional_request,omitempty"`
}
