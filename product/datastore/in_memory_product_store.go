package datastore

import (
	"sort"
	"sync"

	"github.com/shanehowearth/kart/product"
)

// InMemoryProductStore is the concrete in-memory implementation of the product.Store interface.
// It serves as the adapter to fulfil the data contract defined by the core domain.
type InMemoryProductStore struct {
	// mu ensures the products map accesses are thread safe.
	// Only going to be usefule when mutation functions are added to the
	// interface.
	mu sync.RWMutex
	// products is the in memory store. k=Product ID, v = Product.
	products map[string]*product.Product
}

// Ensure that the InMemoryProductStore always satisfies the ProductStore
// interface.
var _ product.Store = (*InMemoryProductStore)(nil)

// NewSeededInMemoryProductStore creates and initialises a new in-memory store.
func NewSeededInMemoryProductStore() *InMemoryProductStore {
	store := &InMemoryProductStore{
		products: make(map[string]*product.Product),
	}

	for i := range SeedProducts {
		product := SeedProducts[i]
		store.products[product.ID] = &product
	}

	return store
}

// NewInMemoryProductStore creates and initialises a new in-memory store.
func NewInMemoryProductStore() *InMemoryProductStore {
	return &InMemoryProductStore{
		products: make(map[string]*product.Product),
	}
}

// GetByID retreives a product by its id.
func (imps *InMemoryProductStore) GetByID(id string) (product.Product, error) {
	// Take a read lock on the map, and release when the function exits.
	imps.mu.RLock()
	defer imps.mu.RUnlock()

	if product, ok := imps.products[id]; ok {
		return *product, nil
	}

	return product.Product{}, product.ErrNotFound
}

// List returns a list of all the products in the datastore.
func (imps *InMemoryProductStore) List() []product.Product {
	// Take a read lock on the map, and release when the function exits.
	imps.mu.RLock()
	defer imps.mu.RUnlock()

	products := make([]product.Product, 0, len(imps.products))

	// Put all the products into the slice.
	for _, product := range imps.products {
		products = append(products, *product)
	}

	// Sort the products by Name.
	// This could be changed to any other field, although Category would also
	// need to be sorted by a subfield.
	sort.Slice(products, func(i, j int) bool {
		return products[i].Name < products[j].Name
	})

	return products
}
