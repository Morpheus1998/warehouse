package products

import (
	"net/http"

	"github.com/warehouse/app/store"
)

type Handler struct {
	ProductsStore store.ProductsStore
}

func NewHandler() *Handler {
	return &Handler{}
}

// CreateOrUpdateProducts is http api /product
func (h *Handler) CreateOrUpdateProducts(w http.ResponseWriter, r *http.Request) {}

// SellProduct is http api DELETE /product/{productId}
func (h *Handler) SellProduct(w http.ResponseWriter, r *http.Request) {}

// GetAllProductsWithStock is http api GET /product
func (h *Handler) GetAllProductsWithStock(w http.ResponseWriter, r *http.Request) {}
