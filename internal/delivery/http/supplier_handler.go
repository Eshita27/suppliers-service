package http

import (
	"net/http"

	"suppliers-api/internal/domain"

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
// @Summary      Register a new vendor
// @Description  Validates and stores a new supplier profile into the system
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        supplier  body      domain.Supplier  true  "Supplier Profiles Payload"
// @Success      201       {object}  map[string]interface{} "Returns success message and supplier data"
// @Failure      400       {object}  map[string]string "Invalid payload error"
// @Failure      422       {object}  map[string]string "Business validation logic failure"
// @Router       /suppliers [post]
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
// @Summary      Get all supply chain vendors
// @Description  Retrieves a full list of registered suppliers from MongoDB
// @Tags         suppliers
// @Produce      json
// @Success      200  {object}  map[string][]domain.Supplier "Returns wrapped supplier list array"
// @Failure      500  {object}  map[string]string "Internal database query failure"
// @Router       /suppliers [get]
func (h *SupplierHandler) FetchSuppliers(c *gin.Context) {
	suppliers, err := h.SUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suppliers matrix"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": suppliers})
}
