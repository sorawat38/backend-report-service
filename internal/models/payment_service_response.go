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
	CartId       string          `json:"cart_id"`
	CustomerId   string          `json:"customer_id"`
	SubTotal     decimal.Decimal `json:"sub_total"`
	DiscountCode string          `json:"discount_code,omitempty"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	DateTime     time.Time       `json:"date_time"`
	Status       string          `json:"status"`
}

type GetCartByIdResponse struct {
	CommonResponse
	Data []GetCartByIdResponseBody `json:"data,omitempty"`
}

type GetCartByIdResponseBody struct {
	CartId            string `json:"cart_id"`
	No                int    `json:"no"`
	CustomerId        string `json:"customer_id"`
	Date              string `json:"date"`
	MenuId            string `json:"menu_id"`
	Quantity          int    `json:"quantity"`
	Status            string `json:"status"`
	Properties        string `json:"properties,omitempty"`
	AdditionalRequest string `json:"additional_request,omitempty"`
}

type GetDiscountByCodeResponse struct {
	CommonResponse
	Data GetDiscountByCodeResponseBody `json:"data,omitempty"`
}

type GetDiscountByCodeResponseBody struct {
	Type  string          `json:"type"`
	Value decimal.Decimal `json:"value"`
}
