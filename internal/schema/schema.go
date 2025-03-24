package schema

import "time"

type RiskScore string
type OrderType string
type CustomerType string

const (
	Low    RiskScore = "low"
	Medium RiskScore = "medium"
	High   RiskScore = "high"

	Buy  OrderType = "buy"
	Sell OrderType = "sell"

	Retail    CustomerType = "retail"
	Workplace CustomerType = "workplace"
)

type Funds struct {
	ID           uint         `gorm:"primaryKey"`
	Name         string       `gorm:"column:name;not null"`
	Description  string       `gorm:"column:description"`
	Code         string       `gorm:"column:code;not null"`
	AmountGBP    float64      `gorm:"column:amount_gbp;not null"`
	CustomerType CustomerType `gorm:"column:customer_type;not null"`
	RiskScore    RiskScore    `gorm:"column:risk_score;not null;type:varchar(50)"`
	LastUpdated  time.Time    `gorm:"column:last_updated"`
}

type Orders struct {
	OrderID           uint      `gorm:"primaryKey"`
	OrderType         OrderType `gorm:"column:order_type;not null;type:varchar(50)"`
	CustomerID        uint      `gorm:"column:customer_id;not null"`
	Name              string    `gorm:"column:name;not null"`
	Description       string    `gorm:"column:description"`
	Code              string    `gorm:"column:code;not null"`
	Shares            float64   `gorm:"column:total_shares;not null"`
	PurchasedValueGBP float64   `gorm:"column:purchased_value_gbp;not null"`
	OrderTime         time.Time `gorm:"column:order_time;not null"`
}
