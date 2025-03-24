package storage

import (
	"context"
	"database/sql"
	"github.com/jautyw/isa-investment-funds/internal/schema"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Store struct {
	db *gorm.DB
}

// NewStore will instantiate a new instance of the Store
func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

const (
	tableFunds  = "funds"
	tableOrders = "orders"

	ErrGettingFunds                     = "error getting funds from db"
	ErrGettingInvestmentOverview        = "error getting investment overview from db"
	ErrGettingAmountSpentCurrentTaxYear = "error getting the amount spent in the current tax year"
)

var (
	// lastYearApril6 refers to the day the new tax year begins
	lastYearApril6 = time.Date(time.Now().Year()-1, 4, 6, 0, 0, 0, 0, time.UTC)
)

func (s *Store) GetFunds(ctx context.Context, customerType string) (*Funds, error) {
	var funds []Fund
	err := s.db.WithContext(ctx).Table(tableFunds).Where("customer_type = ?", customerType).Find(&funds).Error
	if err != nil {
		return nil, errors.Wrap(err, ErrGettingFunds)
	}

	return &Funds{Funds: funds}, nil
}

func (s *Store) GetInvestmentOverview(ctx context.Context, customerID int) ([]InvestmentOverview, error) {
	var investmentOverview []InvestmentOverview

	err := s.db.WithContext(ctx).
		Table(tableOrders).
		Select(`
        name, 
        description, 
        code, 
        SUM(CASE WHEN order_type = 'buy' THEN total_shares ELSE -total_shares END) AS net_shares,
        SUM(CASE WHEN order_type = 'buy' THEN purchased_value_gbp ELSE -purchased_value_gbp END) AS net_investment
    `).
		Where("customer_id = ?", customerID).
		Group("name, description, code").
		Having("SUM(CASE WHEN order_type = 'buy' THEN total_shares ELSE -total_shares END) > 0"). // Ensures only investments with positive net shares are included
		Scan(&investmentOverview).Error
	if err != nil {
		return nil, errors.Wrap(err, ErrGettingInvestmentOverview)
	}

	return investmentOverview, nil
}

func (s *Store) GetAmountSpentCurrentTaxYear(ctx context.Context, customerID int) (float64, error) {
	var allowance sql.NullFloat64
	err := s.db.WithContext(ctx).Table(tableOrders).Select("SUM(purchased_value_gbp)").Where("customer_id = ?", customerID).Where("order_type = ?", schema.Buy).Where("order_time BETWEEN ? AND ?", lastYearApril6, time.Now()).Find(&allowance).Error
	if err != nil {
		return 0, errors.Wrap(err, ErrGettingAmountSpentCurrentTaxYear)
	}

	if !allowance.Valid {
		allowance.Float64 = 0
	}
	return allowance.Float64, nil
}
