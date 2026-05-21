package service

import (
	"context"
	"customer-registry-api/internal/model"
	"testing"
)

// MockRepository to inject into our Service for unit testing
type mockRepo struct {
	createErr error
}

func (m *mockRepo) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	return m.createErr
}
func (m *mockRepo) GetCustomerByID(ctx context.Context, id string) (*model.Customer, error) {
	return nil, nil
}
func (m *mockRepo) GetCustomerByDocument(ctx context.Context, document string) (*model.Customer, error) {
	return nil, nil
}
func (m *mockRepo) ListCustomers(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	return nil, nil
}
func (m *mockRepo) UpdateCustomerStatus(ctx context.Context, id string, status string) error {
	return nil
}

func TestCreateCustomer_Validation(t *testing.T) {
	svc := NewCustomerService(&mockRepo{})

	tests := []struct {
		name    string
		payload model.CreateCustomerRequest
		wantErr error
	}{
		{
			name: "Valid Customer",
			payload: model.CreateCustomerRequest{
				Document: "FAKE-123", Name: "John", Score: 500, RiskLevel: "LOW", IncomeRange: "1000",
			},
			wantErr: nil,
		},
		{
			name: "Invalid Score",
			payload: model.CreateCustomerRequest{
				Document: "FAKE-123", Name: "John", Score: 1500, RiskLevel: "LOW", IncomeRange: "1000",
			},
			wantErr: ErrInvalidScore,
		},
		{
			name: "Invalid Document Prefix",
			payload: model.CreateCustomerRequest{
				Document: "REAL-123", Name: "John", Score: 500, RiskLevel: "LOW", IncomeRange: "1000",
			},
			wantErr: ErrInvalidDocument,
		},
		{
			name: "Invalid Risk Level",
			payload: model.CreateCustomerRequest{
				Document: "FAKE-123", Name: "John", Score: 500, RiskLevel: "EXTREME", IncomeRange: "1000",
			},
			wantErr: ErrInvalidRisk,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateCustomer(context.Background(), tt.payload)
			if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestUpdateStatus_Validation(t *testing.T) {
	svc := NewCustomerService(&mockRepo{})

	err := svc.UpdateStatus(context.Background(), "123", "INVALID_STATUS")
	if err != ErrInvalidStatus {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}

	err = svc.UpdateStatus(context.Background(), "123", "UNDER_REVIEW")
	if err != nil {
		t.Errorf("expected no error for valid status, got %v", err)
	}
}
