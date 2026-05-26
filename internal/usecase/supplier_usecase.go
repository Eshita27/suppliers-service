package usecase

import (
	"context"
	"errors"
	"strings"

	"polyglot/suppliers/internal/domain"
)

type supplierUseCase struct {
	supplierRepo domain.SupplierRepository
}

// NewSupplierUseCase instantiates our business logic engine with its dependencies
func NewSupplierUseCase(repo domain.SupplierRepository) domain.SupplierUseCase {
	return &supplierUseCase{
		supplierRepo: repo,
	}
}

// Create enforces business workflows before handing the data to the repository
func (u *supplierUseCase) Create(ctx context.Context, supplier *domain.Supplier) error {
	// Business Rule: Ensure vendor names aren't empty or blank strings
	if strings.TrimSpace(supplier.Name) == "" {
		return errors.New("supplier name cannot be empty")
	}

	// Business Rule: Normalize email formatting
	supplier.Email = strings.ToLower(strings.TrimSpace(supplier.Email))
	if !strings.Contains(supplier.Email, "@") {
		return errors.New("invalid supplier email address format")
	}

	// Set a default tier if one isn't specified
	if supplier.Tier == "" {
		supplier.Tier = "Standard"
	}

	// Default new suppliers to active status
	supplier.IsActive = true

	// Pass the clean data contract down to the MongoDB repository layer
	return u.supplierRepo.Store(ctx, supplier)
}

// GetAll orchestrates fetching all vendors out of the persistence layer
func (u *supplierUseCase) GetAll(ctx context.Context) ([]domain.Supplier, error) {
	return u.supplierRepo.FetchAll(ctx)
}
