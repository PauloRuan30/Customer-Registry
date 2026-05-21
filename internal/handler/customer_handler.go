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
// @Summary      Criar um novo cliente
// @Description  Regista um novo cliente fictício. O documento é validado e deve iniciar com "FAKE-".
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        request body model.CreateCustomerRequest true "Dados de registo do cliente"
// @Success      201  {object}  model.Customer
// @Failure      400  {object}  model.ErrorResponse
// @Failure      409  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /customers [post]
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Payload JSON inválido", err.Error())
		return
	}

	customer, err := h.svc.CreateCustomer(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, customer)
}

// List godoc
// @Summary      Listar clientes
// @Description  Retorna uma lista paginada de todos os clientes registados no sistema.
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        page    query     int  false  "Número da página" default(1)
// @Param        limit   query     int  false  "Limite de itens por página" default(20)
// @Success      200  {array}   model.Customer
// @Failure      500  {object}  model.ErrorResponse
// @Router       /customers [get]
func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	customers, err := h.svc.ListCustomers(r.Context(), page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Falha ao listar clientes", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, customers)
}

// GetByID godoc
// @Summary      Procurar cliente por ID
// @Description  Recupera as informações completas de um cliente utilizando o seu identificador UUID interno.
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID do Cliente"
// @Success      200  {object}  model.Customer
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /customers/{id} [get]
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	customer, err := h.svc.GetCustomerByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, customer)
}

// GetByDocument godoc
// @Summary      Procurar cliente por Documento
// @Description  Recupera um cliente utilizando o seu documento fictício (ex: FAKE-00001).
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        document  path      string  true  "Documento fictício do Cliente"
// @Success      200       {object}  model.Customer
// @Failure      404       {object}  model.ErrorResponse
// @Failure      500       {object}  model.ErrorResponse
// @Router       /customers/document/{document} [get]
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
// @Summary      Atualizar estado do cliente
// @Description  Modifica apenas o status operacional interno de um cliente específico (ACTIVE, INACTIVE, UNDER_REVIEW).
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "UUID do Cliente"
// @Param        request  body      model.UpdateStatusRequest true "Payload com o novo estado"
// @Success      204      "No Content - Atualização bem-sucedida"
// @Failure      400      {object}  model.ErrorResponse
// @Failure      404      {object}  model.ErrorResponse
// @Failure      500      {object}  model.ErrorResponse
// @Router       /customers/{id}/status [patch]
func (h *CustomerHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req model.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Payload JSON inválido", err.Error())
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, req.Status); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- HTTP Helpers ---

// writeJSON formata as respostas de sucesso
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError formata o envelope de erros padronizado
func writeError(w http.ResponseWriter, status int, msg, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{Error: msg, Details: details})
}

// handleServiceError traduz os erros de domínio (Service/Repo) em códigos HTTP coerentes
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "Recurso não encontrado", "O cliente requisitado não existe")
	case errors.Is(err, repository.ErrDuplicateDocument):
		writeError(w, http.StatusConflict, "Conflito", "Este documento já está registado no sistema")
	case errors.Is(err, service.ErrInvalidDocument),
		errors.Is(err, service.ErrInvalidScore),
		errors.Is(err, service.ErrInvalidRisk),
		errors.Is(err, service.ErrInvalidStatus),
		errors.Is(err, service.ErrNameRequired):
		writeError(w, http.StatusBadRequest, "Falha na validação de dados", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "Erro interno do servidor", err.Error())
	}
}
