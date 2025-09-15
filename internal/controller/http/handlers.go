package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"order-manager/internal/models"

	"github.com/go-chi/chi/v5"
)

type service interface {
	GetOrderByUID(string) (*models.Order, error)
	SaveOrder(*models.Order) error
}

type Handler struct {
	s   service
	log *slog.Logger
}

func NewHandler(s service, log *slog.Logger) *Handler {
	return &Handler{
		s:   s,
		log: log,
	}
}

// @Summary Get order by UID
// @Param order_uid path string true "Order UID"
// @Success 200 {object} models.Order
// @Failure 400 {object} error "Bad request"
// @Failure 404 {object} error "Not found"
// @Router /order/{order_uid} [get]
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "order_uid")

	order, err := h.s.GetOrderByUID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
