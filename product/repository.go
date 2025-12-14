package product

import (
	"errors"
)

// Product repository errors.

// ErrNotFound is returned when a product cannot be found by ID.
var ErrNotFound = errors.New("product not found")

// Store defines the contract for persistent storage operations related
// to the Product entity.
type Store interface {
	// Get a list of products by their ids.
	GetByIDs(ids []string) ([]Product, []string, error)

	// List all products.
	List() []Product

	// TODO - Creation and modification of products is not mentioned in the
	// requirements, but in a production example there would be methods defined
	// here for those purposes.
}
