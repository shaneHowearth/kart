package product

import (
	"errors"

	"github.com/shanehowearth/kart/internal/validation"
)

// ErrNilProductRepo - Error if productStore is nil.
var ErrNilProductRepo = errors.New("product store is nil")

// Service provides business logic for product operations.
type Service struct {
	repo Store
}

// Product is the core domain entity.
type Product struct {
	ID string // Int is implied in the design, but string gives us the
	// flexibility of slugs or UUID, as well as ints that are bigger than int64.
	Name       string
	PriceCents int64 // Price is stored as whole cents, to prevent float math problems.
	Category   string
}

// NewProductService - create a new instance of a product service.
func NewProductService(repo Store) (*Service, error) {
	if validation.IsNil(repo) {
		return nil, ErrNilProductRepo
	}

	return &Service{repo: repo}, nil
}

// GetAvailableProducts gets all available products.
func (ps *Service) GetAvailableProducts() ([]Product, error) {
	products := ps.repo.List()

	// TODO: Business logic would go here - eg. Filtering out discontinued
	// products.

	return products, nil
}

// GetProductsByIDs gets the product details identified by the list of
// ProductIds..
func (ps *Service) GetProductsByIDs(productIds []string) ([]Product, []string, error) {
	products, missed, err := ps.repo.GetByIDs(productIds)

	// TODO: Business logic would go here.

	return products, missed, err
}
