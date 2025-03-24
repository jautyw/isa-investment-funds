package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (h *Handler) GetFunds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Add("Content-Type", "application/json")

	h.Logger.Info("GetFunds request made")

	vars := mux.Vars(r)
	customerType := vars["customer_type"]
	if customerType == "" {
		h.Logger.Error(customerType + "is empty")
		http.Error(w, customerType+"is empty", http.StatusBadRequest)
		return
	}

	getFunds, err := h.Service.GetFunds(ctx, customerType)
	if err != nil {
		h.Logger.Error(errors.Wrap(err, ErrGettingFunds).Error())
		http.Error(w, errors.Wrap(err, ErrGettingFunds).Error(), http.StatusInternalServerError)
		return
	}

	var getFundsResponse GetFundsResponse
	funds := make([]Fund, len(getFunds.Funds))
	for i, f := range getFunds.Funds {
		funds[i] = Fund{
			Name:        f.Name,
			Description: f.Description,
			Code:        f.Code,
			AmountGBP:   f.AmountGBP,
			RiskScore:   f.RiskScore,
			LastUpdated: f.LastUpdated,
		}
	}

	getFundsResponse.Funds = funds

	resBytes, err := json.Marshal(getFundsResponse.Funds)
	if err != nil {
		h.Logger.Error(errors.Wrap(err, ErrGettingFunds).Error())
		http.Error(w, errors.Wrap(err, ErrGettingFunds).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resBytes); err != nil {
		h.Logger.Error(errors.Wrap(err, ErrGettingFunds).Error())
		http.Error(w, errors.Wrap(err, ErrGettingFunds).Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("GetFunds returned successfully")
}

type GetFundsResponse struct {
	Funds []Fund `json:"funds"`
}

type Fund struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	AmountGBP   float64   `json:"amountGBP"`
	RiskScore   string    `json:"riskScore"`
	LastUpdated time.Time `json:"lastUpdated"`
}
