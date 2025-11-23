package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shanehowearth/kart/api"
	"github.com/shanehowearth/kart/order"
	"github.com/shanehowearth/kart/order/datastore/inmemoryorderdatastore"
	"github.com/shanehowearth/kart/product"
	inmemoryproductdatastore "github.com/shanehowearth/kart/product/datastore"
)

const (
	readTimeoutSeconds       = 5
	readHeaderTimeoutSeconds = 3
	writeTimeoutSeconds      = 10
	idleTimeoutSeconds       = 120
)

func main() {
	// Initialise dependencies.
	productStore := inmemoryproductdatastore.NewSeededInMemoryProductStore()

	productService, err := product.NewProductService(productStore)
	if err != nil {
		// Cannot continue, panic with the error message.
		log.Fatalf("Failed to initialize product service: %v", err)
	}

	orderStore := inmemoryorderdatastore.NewInMemoryOrderStore()

	orderService, err := order.NewOrderService(orderStore, productService)
	if err != nil {
		log.Fatalf("Failed to initialize order service: %v", err)
	}

	// Routes.
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, orderService, productService)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       readTimeoutSeconds * time.Second,
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
		WriteTimeout:      writeTimeoutSeconds * time.Second,
		IdleTimeout:       idleTimeoutSeconds * time.Second,
	}

	log.Println("Starting server on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
