package http

import (
	"net/http"

	"polyglot/suppliers/internal/domain"

	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	SUsecase domain.SupplierUseCase
}

// NewSupplierHandler registers the routes under a Gin router group
func NewSupplierHandler(r *gin.Engine, us domain.SupplierUseCase) {
	handler := &SupplierHandler{
		SUsecase: us,
	}

	// Register endpoints under a clean v1 api path
	api := r.Group("/api/v1")
	{
		api.POST("/suppliers", handler.CreateSupplier)
		api.GET("/suppliers", handler.FetchSuppliers)
	}
}

// CreateSupplier handles POST requests to register a new vendor
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var supplier domain.Supplier

	// Bind incoming JSON request body directly into our structural domain model
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload context"})
		return
	}

	// Delegate processing over to our business logic layer
	if err := h.SUsecase.Create(c.Request.Context(), &supplier); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Supplier registered successfully",
		"data":    supplier,
	})
}

// FetchSuppliers handles GET requests to retrieve tracked inventory vendors
func (h *SupplierHandler) FetchSuppliers(c *gin.Context) {
	suppliers, err := h.SUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suppliers matrix"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": suppliers})
}
