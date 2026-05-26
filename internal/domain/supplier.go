package domain

import (
	"context"
	"time"
)

type Supplier struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Contact   string    `json:"contact" bson:"contact"`
	Email     string    `json:"email" bson:"email"`
	Tier      string    `json:"tier" bson:"tier"`
	IsActive  bool      `json:"is_active" bson:"is_active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// SupplierUseCase defines high-level business workflow contracts
type SupplierUseCase interface {
	Create(ctx context.Context, supplier *Supplier) error
	GetAll(ctx context.Context) ([]Supplier, error)
}

// SupplierRepository defines low-level database interaction contracts
type SupplierRepository interface {
	Store(ctx context.Context, supplier *Supplier) error
	FetchAll(ctx context.Context) ([]Supplier, error)
}
