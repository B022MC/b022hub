package service

import (
	"time"

	infraerrors "github.com/B022MC/b022hub/internal/pkg/errors"
)

const (
	PaymentProviderLinuxDoCredit = "linuxdo_credit"

	PaymentOrderStatusPending = "pending"
	PaymentOrderStatusPaid    = "paid"
)

var (
	ErrPaymentOrderNotFound     = infraerrors.NotFound("PAYMENT_ORDER_NOT_FOUND", "payment order not found")
	ErrLinuxDoCreditDisabled    = infraerrors.Forbidden("LINUXDO_CREDIT_DISABLED", "linuxdo credit payment is disabled")
	ErrLinuxDoCreditInvalidAmt  = infraerrors.BadRequest("LINUXDO_CREDIT_INVALID_AMOUNT", "invalid payment amount")
	ErrLinuxDoCreditBadConfig   = infraerrors.ServiceUnavailable("LINUXDO_CREDIT_CONFIG_INVALID", "linuxdo credit payment is not configured correctly")
	ErrLinuxDoCreditBadSign     = infraerrors.BadRequest("LINUXDO_CREDIT_INVALID_SIGNATURE", "invalid linuxdo credit signature")
	ErrLinuxDoCreditOrderDenied = infraerrors.Forbidden("LINUXDO_CREDIT_ORDER_DENIED", "payment order does not belong to current user")
)

type LinuxDoCreditConfig struct {
	Enabled      bool
	ClientID     string
	ClientSecret string
	ExchangeRate float64
}

type PaymentOrder struct {
	ID                 int64      `json:"id"`
	Provider           string     `json:"provider"`
	OutTradeNo         string     `json:"out_trade_no"`
	ProviderTradeNo    string     `json:"provider_trade_no,omitempty"`
	UserID             int64      `json:"user_id"`
	Title              string     `json:"title"`
	Amount             float64    `json:"amount"`
	CreditedAmount     float64    `json:"credited_amount"`
	Status             string     `json:"status"`
	RawProviderPayload string     `json:"raw_provider_payload,omitempty"`
	PaidAt             *time.Time `json:"paid_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type LinuxDoCreditCheckout struct {
	Order       *PaymentOrder
	CheckoutURL string
	Fields      map[string]string
}
