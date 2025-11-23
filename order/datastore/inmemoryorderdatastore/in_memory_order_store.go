package inmemoryorderdatastore

import (
	"fmt"
	"sync"

	"github.com/shanehowearth/kart/order"
)

// InMemoryOrderStore is the concrete in-memory implementation of the order.Store interface.
// It serves as the adapter to fulfil the data contract defined by the core domain.
type InMemoryOrderStore struct {
	// mu ensures the products map accesses are thread safe.
	// Only going to be useful when mutation functions are added to the
	// interface.
	mu sync.RWMutex
	// orders is the in memory store. k=OrderID, v = Order.
	orders map[string]*order.Order
}

// Ensure that the InMemoryOrderStore always satisfies the Store
// interface.
var _ order.Store = (*InMemoryOrderStore)(nil)

// NewInMemoryOrderStore creates and initialises a new in-memory store.
func NewInMemoryOrderStore() *InMemoryOrderStore {
	return &InMemoryOrderStore{
		orders: make(map[string]*order.Order),
	}
}

// CreateOrder - stores a new order.
// Note: This trusts the service layer to provide a valid order, with valid ID.
// The repository layer cannot do any checks on the ID because that is a
// business logic responsibility.
func (imos *InMemoryOrderStore) CreateOrder(
	newOrder *order.Order,
) error {
	imos.mu.Lock()
	defer imos.mu.Unlock()

	// Does the orderID already exist.
	if _, ok := imos.orders[newOrder.ID]; ok {
		return fmt.Errorf("%w order with that ID already exists", order.ErrCreateFailed)
	}

	// Save.
	imos.orders[newOrder.ID] = newOrder

	return nil
}

// GetByID gets an order by id.
func (imos *InMemoryOrderStore) GetByID(orderID string) (order.Order, error) {
	imos.mu.RLock()
	defer imos.mu.RUnlock()

	// Does the orderID already exist.
	if fetched, ok := imos.orders[orderID]; ok {
		return *fetched, nil
	}

	return order.Order{}, fmt.Errorf("%w no order with ID %s",
		order.ErrNotFound,
		orderID,
	)
}
