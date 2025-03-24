package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (h *Handler) GetInvestmentOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("GetInvestmentOverview request made")

	vars := mux.Vars(r)
	customerID, exists := vars["customer_id"]
	if !exists || customerID == "" {
		h.Logger.Error("customer_id is missing")
		http.Error(w, "customer_id is required", http.StatusBadRequest)
		return
	}

	customerIDint, err := strconv.Atoi(customerID)
	if err != nil || customerIDint <= 0 {
		h.Logger.Error(fmt.Sprintf("%s customer_id is invalid", customerID))
		http.Error(w, fmt.Sprintf("%s customer_id is invalid", customerID), http.StatusBadRequest)
		return
	}

	overview, err := h.Service.GetInvestmentOverview(ctx, customerIDint)
	if err != nil {
		h.Logger.Error(errors.Wrap(err, ErrGettingInvestmentOverview).Error())
		http.Error(w, errors.Wrap(err, ErrGettingInvestmentOverview).Error(), http.StatusInternalServerError)
		return
	}

	investments := make([]Investment, len(overview.Investments))
	for i, o := range overview.Investments {
		investments[i] = Investment{
			Name:          o.Name,
			Description:   o.Description,
			Code:          o.Code,
			NetShares:     o.NetShares,
			NetInvestment: o.NetInvestment,
		}
	}

	response := GetInvestmentOverviewResponse{
		Investments:                investments,
		IsaAllowanceCurrentTaxYear: overview.IsaAllowanceCurrentTaxYear,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.Error(errors.Wrap(err, ErrGettingInvestmentOverview).Error())
		http.Error(w, errors.Wrap(err, ErrGettingInvestmentOverview).Error(), http.StatusInternalServerError)
	}

	h.Logger.Info("GetInvestmentOverview returned successfully")
}

type GetInvestmentOverviewResponse struct {
	Investments                []Investment `json:"investments"`
	IsaAllowanceCurrentTaxYear float64      `json:"isaAllowanceCurrentTaxYear"`
}

type Investment struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Code          string  `json:"code"`
	NetShares     float64 `json:"netShares"`
	NetInvestment float64 `json:"netInvestment"`
}
