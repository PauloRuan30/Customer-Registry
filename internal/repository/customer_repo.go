package repository

import (
	"context"
	"customer-registry-api/internal/model"
	"errors"
	"fmt"

	//Github imports
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("customer not found")
var ErrDuplicateDocument = errors.New("customer with this document already exists")

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *model.Customer) error
	GetCustomerByID(ctx context.Context, id string) (*model.Customer, error)
	GetCustomerByDocument(ctx context.Context, document string) (*model.Customer, error)
	ListCustomers(ctx context.Context, limit, offset int) ([]*model.Customer, error)
	UpdateCustomerStatus(ctx context.Context, id string, status string) error
	// DeleteCustomer(ctx context.Context, id string) error;
}

type postgresCustomerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCustomerRepository(db *pgxpool.Pool) CustomerRepository {
	return &postgresCustomerRepository{db: db}
}

// CreateCustomer insere um novo cliente no banco de dados.
func (r *postgresCustomerRepository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	query := `INSERT INTO customers (id, document, name, score, risk_level, income_range, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(ctx, query, customer.ID, customer.Document, customer.Name, customer.Score, customer.RiskLevel, customer.IncomeRange, customer.Status, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			return ErrDuplicateDocument
		}
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

// GetCustomerByID busca um cliente por ID chamando a função auxiliar scanCustomer.
func (r *postgresCustomerRepository) GetCustomerByID(ctx context.Context, id string) (*model.Customer, error) {
	return r.scanCustomer(ctx, "SELECT * FROM customers WHERE id = $1", id)
}

// GetCustomerByDocument busca um cliente por Documento chamando a função auxiliar scanCustomer.
func (r *postgresCustomerRepository) GetCustomerByDocument(ctx context.Context, document string) (*model.Customer, error) {
	return r.scanCustomer(ctx, "SELECT * FROM customers WHERE document = $1", document)
}

// ListCustomers retorna uma lista paginada de clientes.
func (r *postgresCustomerRepository) ListCustomers(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM customers ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []*model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(&c.ID, &c.Document, &c.Name, &c.Score, &c.RiskLevel, &c.IncomeRange, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, &c)
	}
	return customers, rows.Err()
}

// UpdateCustomerStatus atualiza parcialmente o status e a data de modificação.
func (r *postgresCustomerRepository) UpdateCustomerStatus(ctx context.Context, id string, status string) error {
	res, err := r.db.Exec(ctx, "UPDATE customers SET status = $1, updated_at = NOW() WHERE id = $2", status, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

/*
scanCustomer é uma função interna auxiliar (helper) criada para centralizar

	o código de leitura (Scan) de uma única linha.
	Evita duplicação no GetByID e GetByDocument.
*/
func (r *postgresCustomerRepository) scanCustomer(ctx context.Context, query string, args ...any) (*model.Customer, error) {
	var c model.Customer
	err := r.db.QueryRow(ctx, query, args...).Scan(&c.ID, &c.Document, &c.Name, &c.Score, &c.RiskLevel, &c.IncomeRange, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan customer: %w", err)
	}
	return &c, nil
}
