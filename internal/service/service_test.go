package service_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jautyw/isa-investment-funds/internal/service"
	mocks "github.com/jautyw/isa-investment-funds/internal/service/mocks"
	"github.com/jautyw/isa-investment-funds/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStore(ctrl)
	h := service.NewService(ms)
	assert.NotNil(t, h)
}

func TestService_GetFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStore(ctrl)
	h := service.NewService(ms)
	assert.NotNil(t, h)

	ctx := context.Background()

	storeFunds := &storage.Funds{
		Funds: []storage.Fund{
			{
				Name:        "ESG Global All Cap UCITS ETF",
				Description: "Some desc",
				Code:        "V3AM",
				AmountGBP:   4.92,
				RiskScore:   "medium",
				LastUpdated: time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				Name:        "ESG Global All Cap UCITS ETF - (USD) Accumulating",
				Description: "Some desc 2",
				Code:        "V3AB",
				AmountGBP:   4.92,
				RiskScore:   "medium",
				LastUpdated: time.Time{},
			},
		},
	}

	serviceFunds := &service.Funds{
		Funds: []service.Fund{
			{
				Name:        "ESG Global All Cap UCITS ETF",
				Description: "Some desc",
				Code:        "V3AM",
				AmountGBP:   4.92,
				RiskScore:   "medium",
				LastUpdated: time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				Name:        "ESG Global All Cap UCITS ETF - (USD) Accumulating",
				Description: "Some desc 2",
				Code:        "V3AB",
				AmountGBP:   4.92,
				RiskScore:   "medium",
				LastUpdated: time.Time{},
			},
		},
	}

	ms.EXPECT().GetFunds(ctx, "retail").Return(storeFunds, nil).Times(1)

	f, err := h.GetFunds(ctx, "retail")
	assert.NoError(t, err)
	assert.Equal(t, serviceFunds, f)
}

func TestService_GetFundsInvalidCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStore(ctrl)
	h := service.NewService(ms)
	assert.NotNil(t, h)

	ctx := context.Background()

	_, err := h.GetFunds(ctx, "workplace")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "workplace is wrong customer type")
	assert.ErrorContains(t, err, service.ErrGettingFunds)
}

func TestService_GetFundsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStore(ctrl)
	h := service.NewService(ms)
	assert.NotNil(t, h)

	ctx := context.Background()

	ms.EXPECT().GetFunds(ctx, "retail").Return(nil, errors.New(service.ErrGettingFunds)).Times(1)

	_, err := h.GetFunds(ctx, "retail")
	assert.Error(t, err)
	assert.ErrorContains(t, err, service.ErrGettingFunds)
}

func TestService_GetInvestmentOverview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStore(ctrl)
	h := service.NewService(ms)
	assert.NotNil(t, h)

	ctx := context.Background()

	storeFunds := []storage.InvestmentOverview{
		{
			Name:          "ESG Global All Cap UCITS ETF",
			Description:   "Some desc",
			Code:          "V3AM",
			NetShares:     5,
			NetInvestment: 24.6,
		},
		{
			Name:          "ESG Global All Cap UCITS ETF - (USD) Accumulating",
			Description:   "Some desc 2",
			Code:          "V3AB",
			NetShares:     10,
			NetInvestment: 49.2,
		},
	}

	expectedOverview := &service.Overview{
		Investments: []service.InvestmentSummary{
			{
				Name:          "ESG Global All Cap UCITS ETF",
				Description:   "Some desc",
				Code:          "V3AM",
				NetShares:     5,
				NetInvestment: 24.6,
			},
			{
				Name:          "ESG Global All Cap UCITS ETF - (USD) Accumulating",
				Description:   "Some desc 2",
				Code:          "V3AB",
				NetShares:     10,
				NetInvestment: 49.2,
			},
		},
		IsaAllowanceCurrentTaxYear: 19926.2,
	}

	ms.EXPECT().GetInvestmentOverview(ctx, 10000).Return(storeFunds, nil).Times(1)
	ms.EXPECT().GetAmountSpentCurrentTaxYear(ctx, 10000).Return(73.8, nil).Times(1)

	overview, err := h.GetInvestmentOverview(ctx, 10000)
	assert.NoError(t, err)
	assert.Equal(t, expectedOverview, overview)
}

func TestService_GetInvestmentOverviewError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStore(ctrl)
	h := service.NewService(ms)
	assert.NotNil(t, h)

	ctx := context.Background()
	ms.EXPECT().GetInvestmentOverview(ctx, 10000).Return(nil, errors.New(service.ErrGettingOverview)).Times(1)

	_, err := h.GetInvestmentOverview(ctx, 10000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, service.ErrGettingOverview)
}

func TestService_GetInvestmentISAAllowanceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockStore(ctrl)
	h := service.NewService(ms)
	assert.NotNil(t, h)

	ctx := context.Background()

	storeFunds := []storage.InvestmentOverview{
		{
			Name:          "ESG Global All Cap UCITS ETF",
			Description:   "Some desc",
			Code:          "V3AM",
			NetShares:     5,
			NetInvestment: 24.6,
		},
		{
			Name:          "ESG Global All Cap UCITS ETF - (USD) Accumulating",
			Description:   "Some desc 2",
			Code:          "V3AB",
			NetShares:     10,
			NetInvestment: 49.2,
		},
	}

	ms.EXPECT().GetInvestmentOverview(ctx, 10000).Return(storeFunds, nil).Times(1)
	ms.EXPECT().GetAmountSpentCurrentTaxYear(ctx, 10000).Return(float64(0), errors.New(service.ErrGettingISAAllowance)).Times(1)

	_, err := h.GetInvestmentOverview(ctx, 10000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, service.ErrGettingISAAllowance)
}
