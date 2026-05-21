package handler

import (
	"bytes"
	"context"
	"customer-registry-api/internal/model"
	"customer-registry-api/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// 1. Create a dummy Mock Repository just for the handler tests
type mockRepo struct{}

func (m *mockRepo) CreateCustomer(ctx context.Context, customer *model.Customer) error { return nil }
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

func TestCreateCustomer_Handler_BadRequest(t *testing.T) {
	// 2. Setup our architecture layers
	repo := &mockRepo{}
	svc := service.NewCustomerService(repo)
	handler := NewCustomerHandler(svc)

	// 3. Setup the Chi router just like in main.go
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	// 4. Create an intentionally INVALID payload (Missing 'Name')
	payload := model.CreateCustomerRequest{
		Document:  "FAKE-123",
		Score:     500,
		RiskLevel: "LOW",
		// Name is missing!
	}
	body, _ := json.Marshal(payload)

	// 5. Create a fake HTTP POST request
	req, _ := http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 6. Create a ResponseRecorder to act as our fake browser/client
	rr := httptest.NewRecorder()

	// 7. Fire the request into the router
	router.ServeHTTP(rr, req)

	// 8. Assert that the handler properly caught the validation error and returned 400
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// 9. Assert that the error message contains the expected text
	expectedError := "name is required"
	if !bytes.Contains(rr.Body.Bytes(), []byte(expectedError)) {
		t.Errorf("handler returned unexpected body: got %v want it to contain %v", rr.Body.String(), expectedError)
	}
}
