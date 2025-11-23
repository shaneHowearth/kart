package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shanehowearth/kart/order"
)

// OrderHandler provides all the HTTP handlers for the Order domain.
type OrderHandler struct {
	orderService *order.Service
}

// CreateOrderRequest defines the data in the request.
type CreateOrderRequest struct {
	CouponCode string `json:"couponCode"`
	Items      []struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

// NewOrderHandler creates and initialises a new order handler.
func NewOrderHandler(osvc *order.Service) *OrderHandler {
	return &OrderHandler{orderService: osvc}
}

// CreateOrder creates a new order.
func (handler *OrderHandler) CreateOrder(writer http.ResponseWriter, request *http.Request) {
	var req CreateOrderRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate at least one item.
	// TODO: This is an assumption on my part, need to discover if this fits
	// requirements.
	if len(req.Items) == 0 {
		http.Error(writer, "order must contain at least one item", http.StatusBadRequest)
		return
	}

	// Convert to domain types.
	items := make([]order.Item, len(req.Items))
	for i, item := range req.Items {
		items[i] = order.Item{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// Create order.
	newOrder, err := handler.orderService.NewOrder(items)
	if err != nil {
		http.Error(writer, "failed to create order", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(writer).Encode(newOrder); err != nil {
		log.Printf("CreateOrder Encoding JSON failed failed: %v", err)
	}
}

// GetOrder returns an order as specified by its id.
func (handler *OrderHandler) GetOrder(writer http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	fetchedOrder, err := handler.orderService.GetOrderByID(id)
	if err != nil {
		http.Error(writer, "order not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(fetchedOrder); err != nil {
		log.Printf("GetOrder Encoding JSON failed failed: %v", err)
	}
}
