package service

import (
	"isa-investment-funds/internal/schema"
	"time"
)

type Funds struct {
	Funds []Fund
}

type Fund struct {
	Name        string
	Description string
	Code        string
	AmountGBP   float64
	RiskScore   string
	LastUpdated time.Time
}

type Overview struct {
	Investments                []InvestmentSummary
	IsaAllowanceCurrentTaxYear float64
}

type InvestmentSummary struct {
	Name          string
	Description   string
	Code          string
	NetShares     float64
	NetInvestment float64
}

type Order struct {
	OrderID         float64
	OrderType       schema.OrderType
	Name            string
	PurchaseTime    time.Time
	SharesPurchased float64
	AmountGBP       float64
}
