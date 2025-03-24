//go:generate mockgen -destination=./mocks/handler_mock.go -package transport github.com/jautyw/isa-investment-funds/internal/transport Service
package transport

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jautyw/isa-investment-funds/internal/service"
	"go.uber.org/zap"
	"log"
	"net/http"
)

// Handler represents a class that communicates with the service layer
type Handler struct {
	Service Service
	Logger  *zap.Logger
}

// NewHandler will instantiate a new instance of Service
func NewHandler(s Service, l *zap.Logger) *Handler {
	return &Handler{
		Service: s,
		Logger:  l,
	}
}

// Service represents a type that can be used to call the service
type Service interface {
	GetFunds(ctx context.Context, customerType string) (*service.Funds, error)
	GetInvestmentOverview(ctx context.Context, customerID int) (*service.Overview, error)
}

// HandleRequests refers to a collection of endpoints within the service
func (h *Handler) HandleRequests(m *mux.Router) {
	m.HandleFunc("/getFunds/{customer_type}", h.GetFunds).Methods(http.MethodGet)
	m.HandleFunc("/getInvestmentOverview/{customer_id}", h.GetInvestmentOverview).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", m))
}

const (
	ErrGettingFunds              = "/getFunds error"
	ErrGettingInvestmentOverview = "/getInvestmentOverview error"
)
