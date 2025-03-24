package transport_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jautyw/isa-investment-funds/internal/service"
	"github.com/jautyw/isa-investment-funds/internal/transport"
	mocks "github.com/jautyw/isa-investment-funds/internal/transport/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_NewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mocks.NewMockService(ctrl)
	l := zap.NewNop()
	h := transport.NewHandler(ms, l)
	assert.NotNil(t, h)
}

func TestHandler_GetFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mocks.NewMockService(ctrl)
	l := zap.NewNop()
	h := transport.NewHandler(ms, l)

	expectedFunds := &service.Funds{
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

	transportFunds := []transport.Fund{
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
	}

	ms.EXPECT().GetFunds(gomock.Any(), "retail").Return(expectedFunds, nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/getFunds/retail", nil)
	r = mux.SetURLVars(r, map[string]string{"customer_type": "retail"})

	h.GetFunds(w, r)
	res := w.Result()

	var response transport.GetFundsResponse
	err := json.NewDecoder(res.Body).Decode(&response.Funds)
	assert.NoError(t, err)
	expectedResponse := transport.GetFundsResponse{Funds: transportFunds}
	assert.Equal(t, expectedResponse, response)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	err = res.Body.Close()
	assert.NoError(t, err)
}

func TestHandler_GetFundsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mocks.NewMockService(ctrl)
	l := zap.NewNop()
	h := transport.NewHandler(ms, l)

	expectedErrorMsg := transport.ErrGettingFunds
	ms.EXPECT().GetFunds(gomock.Any(), "retail").Return(&service.Funds{
		Funds: []service.Fund{{}},
	}, errors.New(expectedErrorMsg)).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/getFunds/retail", nil)
	r = mux.SetURLVars(r, map[string]string{"customer_type": "retail"})

	h.GetFunds(w, r)
	res := w.Result()

	bodyBytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	actualResponse := string(bodyBytes)
	assert.Contains(t, actualResponse, expectedErrorMsg)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	err = res.Body.Close()
	assert.NoError(t, err)
}

func TestHandler_GetInvestmentOverview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mocks.NewMockService(ctrl)
	l := zap.NewNop()
	h := transport.NewHandler(ms, l)

	expectedFunds := &service.Overview{
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

	transportInvestments := transport.GetInvestmentOverviewResponse{
		Investments: []transport.Investment{
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

	ms.EXPECT().GetInvestmentOverview(gomock.Any(), 10000).Return(expectedFunds, nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/getInvestmentOverview/10000", nil)
	r = mux.SetURLVars(r, map[string]string{"customer_id": "10000"})

	h.GetInvestmentOverview(w, r)
	res := w.Result()

	var response transport.GetInvestmentOverviewResponse
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	expectedResponse := transportInvestments
	assert.Equal(t, expectedResponse, response)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	err = res.Body.Close()
	assert.NoError(t, err)
}

func TestHandler_GetInvestmentOverviewBadRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mocks.NewMockService(ctrl)
	l := zap.NewNop()
	h := transport.NewHandler(ms, l)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/getInvestmentOverview/hi", nil)
	r = mux.SetURLVars(r, map[string]string{"customer_id": "hi"})

	h.GetInvestmentOverview(w, r)
	res := w.Result()

	bodyBytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	actualResponse := string(bodyBytes)
	assert.Contains(t, actualResponse, "customer_id is invalid")
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	err = res.Body.Close()
	assert.NoError(t, err)
}

func TestHandler_GetInvestmentOverviewInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mocks.NewMockService(ctrl)
	l := zap.NewNop()
	h := transport.NewHandler(ms, l)

	ms.EXPECT().GetInvestmentOverview(gomock.Any(), 10000).Return(nil, errors.New(transport.ErrGettingInvestmentOverview)).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/getInvestmentOverview/10000", nil)
	r = mux.SetURLVars(r, map[string]string{"customer_id": "10000"})

	h.GetInvestmentOverview(w, r)
	res := w.Result()

	bodyBytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	actualResponse := string(bodyBytes)
	assert.Contains(t, actualResponse, transport.ErrGettingInvestmentOverview)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	err = res.Body.Close()
	assert.NoError(t, err)
}
