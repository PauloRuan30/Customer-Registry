package handler

import (
	"customer-registry-api/internal/model"
	"customer-registry-api/internal/repository"
	"customer-registry-api/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) RegisterRoutes(r chi.Router) {
	r.Post("/customers", h.Create)
	r.Get("/customers", h.List)
	r.Get("/customers/{id}", h.GetByID)
	r.Get("/customers/document/{document}", h.GetByDocument)
	r.Patch("/customers/{id}/status", h.UpdateStatus)
}

// Create godoc
// @Summary      Create a new customer record
// @Description  Validates and saves a new customer registration with an initial ACTIVE status.
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        request body model.CreateCustomerRequest true "Customer registration payload"
// @Success      201  {object}  model.Customer
// @Failure      400  {object}  model.ErrorResponse
// @Failure      409  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /customers [post]
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	customer, err := h.svc.CreateCustomer(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	customers, err := h.svc.ListCustomers(r.Context(), page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list customers", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	customer, err := h.svc.GetCustomerByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) GetByDocument(w http.ResponseWriter, r *http.Request) {
	doc := chi.URLParam(r, "document")
	customer, err := h.svc.GetCustomerByDocument(r.Context(), doc)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, customer)
}

// UpdateStatus godoc
// @Summary      Update customer status
// @Description  Partially modifies a customer profile status state (ACTIVE, INACTIVE, UNDER_REVIEW).
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Customer UUID"
// @Param        request body model.UpdateStatusRequest true "Target status payload"
// @Success      204
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /customers/{id}/status [patch]
func (h *CustomerHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req model.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, req.Status); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- HTTP Helpers ---

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{Error: msg, Details: details})
}

// handleServiceError translates domain errors into HTTP status codes
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	// 1. Check for specific sentinel errors first using errors.Is
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "Resource not found", "The requested customer does not exist")
	case errors.Is(err, repository.ErrDuplicateDocument):
		writeError(w, http.StatusConflict, "Conflict", "This document is already registered")

	// 2. Check for service validation errors
	case errors.Is(err, service.ErrInvalidDocument),
		errors.Is(err, service.ErrInvalidScore),
		errors.Is(err, service.ErrInvalidRisk),
		errors.Is(err, service.ErrInvalidStatus),
		errors.Is(err, service.ErrNameRequired):
		writeError(w, http.StatusBadRequest, "Validation failed", err.Error())

	// 3. Fallback for unexpected failures (database down, etc.)
	default:
		writeError(w, http.StatusInternalServerError, "Internal server error", err.Error())
	}
}
