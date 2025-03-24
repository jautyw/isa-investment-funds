package storage_test

import (
	"context"
	"fmt"
	"github.com/jautyw/isa-investment-funds/config"
	"github.com/jautyw/isa-investment-funds/internal/schema"
	"github.com/jautyw/isa-investment-funds/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

func setupDB(ctx context.Context) (*gorm.DB, func(), testcontainers.Container) {
	ctx = context.Background()

	cfg := config.Config{
		Host:     "localhost",
		User:     "user",
		Password: "password",
		Database: "investments",
		Port:     "9920",
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", cfg.Port)},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
			"POSTGRES_DB":       cfg.Database,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}

	err = db.AutoMigrate(&schema.Funds{}, &schema.Orders{}) // Example model
	if err != nil {
		log.Fatalf("Failed to migrate database: %s", err)
	}

	teardown := func() {
		err := container.Terminate(ctx)
		if err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}

	return db, teardown, container
}

func cleanDB(db *gorm.DB) error {
	if err := db.Exec(fmt.Sprintf("DELETE FROM %s", "funds")).Error; err != nil {
		return fmt.Errorf("failed to clear table %s: %w", "funds", err)
	}

	if err := db.Exec(fmt.Sprintf("DELETE FROM %s", "orders")).Error; err != nil {
		return fmt.Errorf("failed to clear table %s: %w", "orders", err)
	}

	return nil
}

func TestStoreSuite(t *testing.T) {

}

func TestStore_GetFunds(t *testing.T) {
	ctx := context.Background()
	db, teardown, _ := setupDB(ctx)
	defer teardown()

	err := cleanDB(db)
	assert.NoError(t, err)

	fund := schema.Funds{
		ID:           1,
		Name:         "ESG Global All Cap UCITS ETF",
		Description:  "Some fund",
		Code:         "V3AM",
		AmountGBP:    4.92,
		CustomerType: schema.Retail,
		RiskScore:    schema.Medium,
		LastUpdated:  time.Date(time.Now().Year()-1, 1, 0, 0, 0, 0, 0, time.Local),
	}

	s := storage.NewStore(db)
	err = db.Create(&fund).Error
	assert.NoError(t, err)

	funds, err := s.GetFunds(ctx, "retail")
	assert.NoError(t, err)
	assert.Equal(t, &storage.Funds{Funds: []storage.Fund{
		{
			Name:        fund.Name,
			Description: fund.Description,
			Code:        fund.Code,
			AmountGBP:   fund.AmountGBP,
			RiskScore:   fund.RiskScore,
			LastUpdated: fund.LastUpdated,
		},
	}}, funds)
}

func TestStore_GetInvestmentOverview(t *testing.T) {
	ctx := context.Background()
	db, teardown, _ := setupDB(ctx)
	defer teardown()

	err := cleanDB(db)
	assert.NoError(t, err)

	order := schema.Orders{
		OrderID:           1,
		OrderType:         schema.Buy,
		CustomerID:        11,
		Name:              "ESG Global All Cap UCITS ETF",
		Description:       "Some fund",
		Code:              "V3AM",
		Shares:            4,
		PurchasedValueGBP: 200,
		OrderTime:         time.Date(time.Now().Year()-1, 1, 0, 0, 0, 0, 0, time.Local),
	}

	s := storage.NewStore(db)
	err = db.Create(&order).Error
	assert.NoError(t, err)

	investments, err := s.GetInvestmentOverview(ctx, 11)
	assert.NoError(t, err)
	assert.Equal(t, []storage.InvestmentOverview{
		{
			Name:          "ESG Global All Cap UCITS ETF",
			Description:   "Some fund",
			Code:          "V3AM",
			NetShares:     4,
			NetInvestment: 200,
		}}, investments)
}

func TestStore_GetAmountSpentCurrentTaxYear(t *testing.T) {
	ctx := context.Background()
	db, teardown, _ := setupDB(ctx)
	defer teardown()

	err := cleanDB(db)
	assert.NoError(t, err)

	order := schema.Orders{
		OrderID:           1,
		OrderType:         schema.Buy,
		CustomerID:        11,
		Name:              "ESG Global All Cap UCITS ETF",
		Description:       "Some fund",
		Code:              "V3AM",
		Shares:            4,
		PurchasedValueGBP: 200,
		OrderTime:         time.Date(time.Now().Year(), 1, 0, 0, 0, 0, 0, time.Local),
	}

	s := storage.NewStore(db)
	err = db.Create(&order).Error
	assert.NoError(t, err)

	allowance, err := s.GetAmountSpentCurrentTaxYear(ctx, 11)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), allowance)
}

func TestStore_GetAmountSpentCurrentTaxYearFullAllowance(t *testing.T) {
	ctx := context.Background()
	db, teardown, _ := setupDB(ctx)
	defer teardown()

	err := cleanDB(db)
	assert.NoError(t, err)

	order := schema.Orders{
		OrderID:           1,
		OrderType:         schema.Buy,
		CustomerID:        11,
		Name:              "ESG Global All Cap UCITS ETF",
		Description:       "Some fund",
		Code:              "V3AM",
		Shares:            4,
		PurchasedValueGBP: 200,
		// Order made before the current tax year so allowance should be 0
		OrderTime: time.Date(time.Now().Year()-1, 1, 0, 0, 0, 0, 0, time.Local),
	}

	s := storage.NewStore(db)
	err = db.Create(&order).Error
	assert.NoError(t, err)

	allowance, err := s.GetAmountSpentCurrentTaxYear(ctx, 11)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), allowance)
}

func TestStore_GetAmountSpentCurrentTaxYearError(t *testing.T) {
	ctx := context.Background()
	db, teardown, container := setupDB(ctx)
	defer teardown()

	err := cleanDB(db)
	assert.NoError(t, err)

	order := schema.Orders{
		OrderID:           1,
		OrderType:         schema.Buy,
		CustomerID:        11,
		Name:              "ESG Global All Cap UCITS ETF",
		Description:       "Some fund",
		Code:              "V3AM",
		Shares:            4,
		PurchasedValueGBP: 200,
		OrderTime:         time.Date(time.Now().Year()-1, 1, 0, 0, 0, 0, 0, time.Local),
	}

	s := storage.NewStore(db)
	err = db.Create(&order).Error
	assert.NoError(t, err)

	err = container.Terminate(ctx) // Stops the database container
	require.NoError(t, err)

	_, err = s.GetAmountSpentCurrentTaxYear(ctx, 11)
	assert.Error(t, err)
	assert.ErrorContains(t, err, storage.ErrGettingAmountSpentCurrentTaxYear)
}
