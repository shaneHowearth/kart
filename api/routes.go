package api

import (
	"net/http"

	"github.com/shanehowearth/kart/api/handlers"
	"github.com/shanehowearth/kart/order"
	"github.com/shanehowearth/kart/product"
)

// Allow anyone access to the API (Cross origin requests)
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: In production, replace "http://localhost:3000" with actual domain.
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		// Set allowed methods (needed for preflight and actual request)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Set allowed headers (essential for custom headers like Authorization and Content-Type)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's an OPTIONS request, we send the headers and respond with 200 OK immediately,
		// preventing the request from reaching the actual handler.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// RegisterRoutes register all the routes for the API.
func RegisterRoutes(
	mux *http.ServeMux,
	orderService *order.Service,
	productService *product.Service,
) {
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Order routes.
	mux.Handle("GET /api/order/{id}", CORSMiddleware(http.HandlerFunc(orderHandler.GetOrder)))
	mux.Handle("POST /api/order", CORSMiddleware(http.HandlerFunc(orderHandler.CreateOrder)))
	// Allow OPTIONS in order to prevent a CORS issue.
	mux.Handle(
		"OPTIONS /api/order",
		CORSMiddleware(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			),
		),
	)

	// Product routes.
	mux.Handle("GET /api/product", CORSMiddleware(http.HandlerFunc(productHandler.ListProducts)))
	mux.Handle("GET /api/product/{id}", CORSMiddleware(http.HandlerFunc(productHandler.GetProduct)))
}
