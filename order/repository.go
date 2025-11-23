package order

import "errors"

// Order repository errors.

//nolint:revive // Sentinal errors, no need to comment.
var (
	ErrCreateFailed = errors.New("creating order failed")
	ErrNotFound     = errors.New("order not found")
)

// Store defines the contract for persistent storage operations related
// to the Order entity.
type Store interface {
	// Create an Order.
	CreateOrder(*Order) error
	// GetByID returns an order that has the supplied ID.
	GetByID(string) (Order, error)
}
