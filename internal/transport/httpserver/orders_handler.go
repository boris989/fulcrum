package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/boris989/fulcrum/internal/orders/app"
)

type OrdersHandler struct {
	svc *app.Service
}

func RegisterOrders(mux *http.ServeMux, svc *app.Service) {
	h := &OrdersHandler{svc: svc}

	mux.HandleFunc("/orders", h.handleCreate)
	mux.HandleFunc("/orders/", h.handlePay)
}

type createOrderRequest struct {
	Amount int64 `json:"amount"`
}

type createOrderResponse struct {
	ID string `json:"id"`
}

func (h *OrdersHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req createOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	order, err := h.svc.CreateOrder(r.Context(), req.Amount)
	if err != nil {
		mapError(w, err)
		return
	}

	resp := createOrderResponse{ID: order.ID()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *OrdersHandler) handlePay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "pay" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id := parts[0]

	if err := h.svc.PayOrder(r.Context(), id); err != nil {
		mapError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func mapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, app.ErrOptimisticLock):
		http.Error(w, "conflict", http.StatusConflict)
	case err.Error() == "order not found":
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
