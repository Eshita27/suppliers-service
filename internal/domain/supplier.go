package domain

import (
	"context"
	"time"
)

// Supplier represents the core data model for a supply chain vendor.
type Supplier struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Contact   string    `json:"contact" bson:"contact"`
	Email     string    `json:"email" bson:"email"`
	Tier      string    `json:"tier" bson:"tier"` // e.g., Tier-1 Premium, Tier-2 Standard
	IsActive  bool      `json:"is_active" bson:"is_active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// SupplierUseCase defines the business contract our API will execute.
type SupplierUseCase interface {
	Create(ctx context.Context, supplier *Supplier) error
	GetAll(ctx context.Context) ([]Supplier, error)
}

// SupplierRepository defines the persistence contract for MongoDB.
type SupplierRepository interface {
	Store(ctx context.Context, supplier *Supplier) error
	FetchAll(ctx context.Context) ([]Supplier, error)
}
