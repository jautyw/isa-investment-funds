package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/jautyw/isa-investment-funds/config"
	"github.com/jautyw/isa-investment-funds/internal/logger"
	"github.com/jautyw/isa-investment-funds/internal/schema"
	"github.com/jautyw/isa-investment-funds/internal/service"
	"github.com/jautyw/isa-investment-funds/internal/storage"
	"github.com/jautyw/isa-investment-funds/internal/transport"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config %e", err)
	}

	l := logger.NewLogger()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error opening postgres %v", err)
	}

	// Run migrations to ensure the tables are created or updated
	err = db.AutoMigrate(&schema.Funds{}, &schema.Orders{})
	if err != nil {
		log.Fatalf("error running migrations: %v", err)
	}

	// Add some mock data
	if err := SeedDatabase(db); err != nil {
		log.Fatalf("failed to seed db: %v", err)
	}

	// Instantiate and inject each layer of the service
	st := storage.NewStore(db)
	s := service.NewService(st)
	t := transport.NewHandler(s, l)

	r := mux.NewRouter().StrictSlash(true)
	t.HandleRequests(r)
}

// SeedDatabase exists purely for initial and subsequent local testing.
func SeedDatabase(db *gorm.DB) error {
	if err := db.Exec(fmt.Sprintf("DELETE FROM %s", "funds")).Error; err != nil {
		return fmt.Errorf("failed to clear table %s: %w", "funds", err)
	}

	if err := db.Exec(fmt.Sprintf("DELETE FROM %s", "orders")).Error; err != nil {
		return fmt.Errorf("failed to clear table %s: %w", "orders", err)
	}

	funds := []schema.Funds{
		{
			ID:           1,
			Name:         "ESG Global All Cap UCITS ETF",
			Description:  "Some fund",
			Code:         "V3AM",
			AmountGBP:    4.92,
			CustomerType: schema.Retail,
			RiskScore:    schema.Medium,
			LastUpdated:  time.Now(),
		},
		{
			ID:           2,
			Name:         "ESG Global All Cap UCITS ETF - (USD) Accumulating",
			Description:  "Some fund",
			Code:         "V3AB",
			AmountGBP:    4.92,
			CustomerType: schema.Retail,
			RiskScore:    schema.Medium,
			LastUpdated:  time.Now(),
		},
		{
			ID:           3,
			Name:         "ESG Global All Cap UCITS ETF",
			Description:  "Some fund",
			Code:         "V3AM",
			AmountGBP:    4.92,
			CustomerType: schema.Workplace,
			RiskScore:    schema.Medium,
			LastUpdated:  time.Now(),
		},
		{
			ID:           4,
			Name:         "ESG Global All Cap UCITS ETF - (USD) Accumulating",
			Description:  "Some fund",
			Code:         "V3AB",
			AmountGBP:    4.92,
			CustomerType: schema.Workplace,
			RiskScore:    schema.Medium,
			LastUpdated:  time.Now(),
		},
	}

	for _, fund := range funds {
		if err := db.Create(&fund).Error; err != nil {
			return fmt.Errorf("failed to insert fund data: %w", err)
		}
	}

	orders := []schema.Orders{
		{
			OrderID:           1,
			OrderType:         schema.Buy,
			CustomerID:        1,
			Name:              "ESG Global All Cap UCITS ETF",
			Description:       "Some fund",
			Code:              "V3AM",
			Shares:            100.0,
			PurchasedValueGBP: 492.0,
			OrderTime:         time.Now(),
		},
		{
			OrderID:           2,
			OrderType:         schema.Sell,
			CustomerID:        1,
			Name:              "ESG Global All Cap UCITS ETF",
			Description:       "Some fund",
			Code:              "V3AM",
			Shares:            50.0,
			PurchasedValueGBP: 246.0,
			OrderTime:         time.Now(),
		},
		{
			OrderID:           3,
			OrderType:         schema.Buy,
			CustomerID:        1,
			Name:              "ESG Global All Cap UCITS ETF - (USD) Accumulating",
			Description:       "Some fund",
			Code:              "V3AB",
			Shares:            100.0,
			PurchasedValueGBP: 492.0,
			OrderTime:         time.Now(),
		},
		{
			OrderID:           4,
			OrderType:         schema.Sell,
			CustomerID:        1,
			Name:              "ESG Global All Cap UCITS ETF - (USD) Accumulating",
			Description:       "Some fund",
			Code:              "V3AB",
			Shares:            50.0,
			PurchasedValueGBP: 246.0,
			OrderTime:         time.Now(),
		},
	}

	for _, order := range orders {
		if err := db.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to insert order data: %w", err)
		}
	}

	return nil
}
