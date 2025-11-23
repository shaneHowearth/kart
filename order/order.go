package order

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/shanehowearth/kart/internal/validation"
	"github.com/shanehowearth/kart/product"
)

// ProductGetter defines the contract for fetching product details.
type ProductGetter interface {
	GetProductByID(id string) (product.Product, error)
}

// Service provides business logic for order operations.
type Service struct {
	repo          Store
	productGetter ProductGetter
}

// ErrCannotCreateOrderService - Error if Order cannot be created.
var ErrCannotCreateOrderService = errors.New("cannot create order service")

// Item provides the structure to hold order item details.
type Item struct {
	ProductID string
	Quantity  int // Note: Current implementation only supports integer quantities only.
	// Fractional quantities would require schema changes.
}

// Order is the structure to hold the Order details.
type Order struct {
	ID       string
	Items    []Item
	Products []ProductReference
}

// ProductReference is the value object within the Order aggregate.
type ProductReference struct {
	ID         string
	Name       string
	PriceCents int64
	Category   string
}

// NewOrderService - create a new instance of a order service.
func NewOrderService(repo Store, productGetter ProductGetter) (*Service, error) {
	if validation.IsNil(repo) {
		return nil, fmt.Errorf("%w order store is nil", ErrCannotCreateOrderService)
	}

	if validation.IsNil(productGetter) {
		return nil, fmt.Errorf("%w product getter is nil", ErrCannotCreateOrderService)
	}

	return &Service{
		repo:          repo,
		productGetter: productGetter,
	}, nil
}

// NewOrder creates a new order.
func (svc *Service) NewOrder(
	items []Item,
) (Order, error) {
	// Order must have at least 1 item.
	// TODO ensure that this matches expected business requirements.
	if len(items) < 1 {
		return Order{}, fmt.Errorf("%w no items", ErrCreateFailed)
	}

	// Fetch current product information for this order.
	productReferences := make([]ProductReference, 0, len(items))

	for idx := range items {
		productInfo, err := svc.productGetter.GetProductByID(items[idx].ProductID)
		if err != nil {
			// TODO: Not sure if this is a catastrophic error, or not.  Am
			// treating it as catastrophic because order fulfilment, and
			// billing, will be compromised.
			return Order{}, fmt.Errorf("%w product %s not found: %v", ErrCreateFailed, items[idx].ProductID, err)
		}

		productReferences = append(productReferences, ProductReference{
			ID:         productInfo.ID,
			Name:       productInfo.Name,
			PriceCents: productInfo.PriceCents,
			Category:   productInfo.Category,
		})
	}

	orderID := uuid.New().String()

	// Persist the order.
	newOrder := Order{
		ID:       orderID,
		Items:    items,
		Products: productReferences,
	}

	err := svc.repo.CreateOrder(&newOrder)
	if err != nil {
		// TODO: Need clarification on surfacing repository errors to the caller.
		// log full error here and return simpler error.
		log.Printf("%v with repository error %v", ErrCreateFailed, err)
		return Order{}, fmt.Errorf("%w repository error", ErrCreateFailed)
	}

	return newOrder, nil
}

// GetOrderByID gets a single order by id.
func (svc *Service) GetOrderByID(id string) (Order, error) {
	order, err := svc.repo.GetByID(id)

	return order, err
}
