package api

import (
	"net/http"

	"github.com/shanehowearth/kart/api/handlers"
	"github.com/shanehowearth/kart/order"
	"github.com/shanehowearth/kart/product"
)

// RegisterRoutes register all the routes for the API.
func RegisterRoutes(
	mux *http.ServeMux,
	orderService *order.Service,
	productService *product.Service,
) {
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Order routes.
	mux.HandleFunc("GET /api/order/{id}", orderHandler.GetOrder)
	mux.HandleFunc("POST /api/order", orderHandler.CreateOrder)

	// Product routes.
	mux.HandleFunc("GET /api/product", productHandler.ListProducts)
	mux.HandleFunc("GET /api/product/{id}", productHandler.GetProduct)
}
