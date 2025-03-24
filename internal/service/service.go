//go:generate mockgen -destination=./mocks/service_mock.go -package service github.com/jautyw/isa-investment-funds/internal/service Store
package service

import (
	"context"
	"fmt"
	"github.com/jautyw/isa-investment-funds/internal/schema"
	"github.com/jautyw/isa-investment-funds/internal/storage"
	"github.com/pkg/errors"
)

type Service struct {
	store Store
}

// NewService represents a new instance of the Service
func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

const (
	ErrGettingFunds        = "error getting funds for user"
	ErrGettingOverview     = "error getting overview for user"
	ErrGettingISAAllowance = "error getting allowance for user"

	// isaAnnualGovernmentAllowance refers to the amount customers can save tax-free
	isaAnnualGovernmentAllowance = 20000
)

// Store represents a collection of methods that can be used to call the store
type Store interface {
	GetFunds(ctx context.Context, customerType string) (*storage.Funds, error)
	GetInvestmentOverview(ctx context.Context, customerID int) ([]storage.InvestmentOverview, error)
	GetAmountSpentCurrentTaxYear(ctx context.Context, customerID int) (float64, error)
}

func (s Service) GetFunds(ctx context.Context, customerType string) (*Funds, error) {
	// As desired we only want to provide information to retail customers so at this stage we block the operation.
	if customerType != string(schema.Retail) {
		return nil, errors.Wrap(errors.New(fmt.Sprintf("%s is wrong customer type", customerType)), ErrGettingFunds)
	}

	// We return all the funds available to "retail" customers
	storeFunds, err := s.store.GetFunds(ctx, customerType)
	if err != nil {
		return nil, errors.Wrap(err, ErrGettingFunds)
	}

	f := make([]Fund, len(storeFunds.Funds))
	for i, sf := range storeFunds.Funds {
		f[i] = Fund{
			Name:        sf.Name,
			Description: sf.Description,
			Code:        sf.Code,
			AmountGBP:   sf.AmountGBP,
			RiskScore:   string(sf.RiskScore),
			LastUpdated: sf.LastUpdated,
		}
	}

	funds := &Funds{Funds: f}

	return funds, nil
}

func (s Service) GetInvestmentOverview(ctx context.Context, customerID int) (*Overview, error) {
	// At the moment we only expect a customer to have invested a single type of fund, however, we should keep in mind
	// that we probably want to support multiple funds in the future.
	// TODO We should also check whether a user exists here so that we don't just return an empty array.

	investmentSummaries, err := s.store.GetInvestmentOverview(ctx, customerID)
	if err != nil {
		return nil, errors.Wrap(err, ErrGettingOverview)
	}

	// Now we check the total amount the customer has invested in the current tax year to see what remains.
	totalInvestedCurrentTaxYear, err := s.store.GetAmountSpentCurrentTaxYear(ctx, customerID)
	if err != nil {
		return nil, errors.Wrap(err, ErrGettingISAAllowance)
	}

	is := make([]InvestmentSummary, len(investmentSummaries))
	for i, sis := range investmentSummaries {
		is[i] = InvestmentSummary{
			Name:          sis.Name,
			Description:   sis.Description,
			Code:          sis.Code,
			NetShares:     sis.NetShares,
			NetInvestment: sis.NetInvestment,
		}
	}

	overview := &Overview{
		Investments: is,
	}

	// We check if they have already maxed out their annual tax free allowance
	if totalInvestedCurrentTaxYear >= isaAnnualGovernmentAllowance {
		overview.IsaAllowanceCurrentTaxYear = 0
	}

	overview.IsaAllowanceCurrentTaxYear = isaAnnualGovernmentAllowance - totalInvestedCurrentTaxYear

	return overview, nil
}
