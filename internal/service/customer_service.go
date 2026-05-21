package service

import (
	"context"
	"customer-registry-api/internal/model"
	"customer-registry-api/internal/repository"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDocument = errors.New("document must start with FAKE-")
	ErrInvalidScore    = errors.New("score must be between 0 and 1000")
	ErrInvalidRisk     = errors.New("risk_level must be LOW, MEDIUM or HIGH")
	ErrInvalidStatus   = errors.New("status must be ACTIVE, INACTIVE or UNDER_REVIEW")
	ErrNameRequired    = errors.New("name is required")
)

type CustomerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (service *CustomerService) CreateCustomer(ctx context.Context, req model.CreateCustomerRequest) (*model.Customer, error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	customer := &model.Customer{
		ID:          uuid.New().String(),
		Document:    strings.ToUpper(req.Document),
		Name:        req.Name,
		Score:       req.Score,
		RiskLevel:   strings.ToUpper(req.RiskLevel),
		IncomeRange: req.IncomeRange,
		Status:      "ACTIVE", //Default status for new customers
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := service.repo.CreateCustomer(ctx, customer); err != nil {
		if errors.Is(err, repository.ErrDuplicateDocument) {
			return nil, fmt.Errorf("document %s already registered", customer.Document)
		}
		return nil, fmt.Errorf("service error: %w", err)
	}
	return customer, nil
}

func (service *CustomerService) GetCustomerByID(ctx context.Context, id string) (*model.Customer, error) {
	return service.repo.GetCustomerByID(ctx, id)
}

func (service *CustomerService) GetCustomerByDocument(ctx context.Context, doc string) (*model.Customer, error) {
	return service.repo.GetCustomerByDocument(ctx, strings.ToUpper(doc))
}

func (service *CustomerService) ListCustomers(ctx context.Context, page, limit int) ([]*model.Customer, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return service.repo.ListCustomers(ctx, limit, offset)
}

func (service *CustomerService) UpdateStatus(ctx context.Context, id, status string) error {
	status = strings.ToUpper(status)
	if !isValidStatus(status) {
		return ErrInvalidStatus
	}
	return service.repo.UpdateCustomerStatus(ctx, id, status)
}

// --- Validation Helpers ---

func validateCreateRequest(req model.CreateCustomerRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrNameRequired
	}
	if !regexp.MustCompile(`^FAKE-\d+$`).MatchString(req.Document) {
		return ErrInvalidDocument
	}
	if req.Score < 0 || req.Score > 1000 {
		return ErrInvalidScore
	}
	if !isValidRisk(req.RiskLevel) {
		return ErrInvalidRisk
	}
	return nil
}

func isValidRisk(r string) bool {
	r = strings.ToUpper(r)
	return r == "LOW" || r == "MEDIUM" || r == "HIGH"
}

func isValidStatus(s string) bool {
	return s == "ACTIVE" || s == "INACTIVE" || s == "UNDER_REVIEW"
}
