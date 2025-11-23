package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shanehowearth/kart/product"
)

// ProductHandler provides all the HTTP handlers for the Product domain.
type ProductHandler struct {
	productService *product.Service
}

// NewProductHandler creates a new product handler.
func NewProductHandler(ps *product.Service) *ProductHandler {
	return &ProductHandler{productService: ps}
}

// ProductResponse details what data and how it is formatted is responded for a
// product - it's a DTO.
type ProductResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PriceDisplay string `json:"price"` // "$5.99"
	Category     string `json:"category"`
}

const centsPerDollar = 100

func formatPrice(cents int64) string {
	dollars := cents / centsPerDollar
	remainingCents := cents % centsPerDollar

	// Note: I provide trailing zeros, which the demonstration example does not
	// do.
	return fmt.Sprintf("$%d.%02d", dollars, remainingCents)
}

// ListProducts lists all the products.
func (h *ProductHandler) ListProducts(writer http.ResponseWriter, _ *http.Request) {
	products, err := h.productService.GetAvailableProducts()
	if err != nil {
		http.Error(writer, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	// Convert products to displayable content.
	displayableProducts := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		productResponse := ProductResponse{}
		productResponse.ID = product.ID
		productResponse.Name = product.Name
		productResponse.Category = product.Category
		productResponse.PriceDisplay = formatPrice(product.PriceCents)
		displayableProducts = append(displayableProducts, productResponse)
	}

	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(displayableProducts); err != nil {
		log.Printf("ListProducts Encoding JSON failed failed: %v", err)
	}
}

// GetProduct gets a single product.
func (h *ProductHandler) GetProduct(writer http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	fetchedProduct, err := h.productService.GetProductByID(id)
	if err != nil {
		http.Error(writer, "product not found", http.StatusNotFound)
		return
	}

	productResponse := ProductResponse{}
	productResponse.ID = fetchedProduct.ID
	productResponse.Name = fetchedProduct.Name
	productResponse.Category = fetchedProduct.Category
	productResponse.PriceDisplay = formatPrice(fetchedProduct.PriceCents)

	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(productResponse); err != nil {
		log.Printf("GetProduct Encoding JSON failed failed: %v", err)
	}
}
