package storage

import (
	"isa-investment-funds/internal/schema"
	"time"
)

type Funds struct {
	Funds []Fund
}

type Fund struct {
	Name        string           `gorm:"name"`
	Description string           `gorm:"description"`
	Code        string           `gorm:"code"`
	AmountGBP   float64          `gorm:"amount_gbp"`
	RiskScore   schema.RiskScore `gorm:"risk_score"`
	LastUpdated time.Time        `gorm:"last_updated"`
}

type Overview struct {
	Investments []InvestmentOverview
}

type InvestmentOverview struct {
	Name          string  `gorm:"column:name"`
	Description   string  `gorm:"column:description"`
	Code          string  `gorm:"column:code"`
	NetShares     float64 `gorm:"column:net_shares"`
	NetInvestment float64 `gorm:"column:net_investment"`
}

type Order struct {
	OrderID         float64          `gorm:"order_id"`
	Name            string           `gorm:"name"`
	PurchaseTime    time.Time        `gorm:"purchase_time"`
	OrderType       schema.OrderType `gorm:"orderType"`
	SharesPurchased float64          `gorm:"shares_purchased"`
	AmountGBP       float64          `gorm:"amount_gbp"`
}
