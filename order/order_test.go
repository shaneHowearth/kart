//nolint:varnamelen // id and tc are clear enough.
package order_test

import (
	"fmt"
	"testing"

	"github.com/shanehowearth/kart/order"
	"github.com/shanehowearth/kart/order/datastore/inmemoryorderdatastore"
	"github.com/shanehowearth/kart/product"
	"github.com/stretchr/testify/assert"
)

type MockProductGetter struct {
	products map[string]product.Product
	err      error // To simulate errors.
}

func (m *MockProductGetter) GetProductById(id string) (product.Product, error) {
	if m.err != nil {
		return product.Product{}, m.err
	}

	if p, ok := m.products[id]; ok {
		return p, nil
	}

	return product.Product{}, product.ErrNotFound
}

type MockOrderStore struct {
	orders map[string]*order.Order
	err    error // Set this to force errors.
}

func (m *MockOrderStore) CreateOrder(o *order.Order) error {
	if m.err != nil {
		return m.err
	}

	m.orders[o.ID] = o

	return nil
}

func (m *MockOrderStore) GetByID(id string) (order.Order, error) {
	if o, ok := m.orders[id]; ok {
		return *o, nil
	}

	return order.Order{}, order.ErrNotFound
}

func TestNewOrderService(t *testing.T) {
	testcases := map[string]struct {
		orderStore    order.Store
		productGetter order.ProductGetter
		expectedError error
	}{
		"New Order service created": {
			orderStore: inmemoryorderdatastore.NewInMemoryOrderStore(),
			productGetter: &MockProductGetter{
				products: map[string]product.Product{
					"1": {ID: "1", Name: "Test", PriceCents: 100},
				},
			},
		},
		"No product getter supplied causes error": {
			orderStore:    inmemoryorderdatastore.NewInMemoryOrderStore(),
			expectedError: order.ErrCannotCreateOrderService,
		},
		"New order store supplied causes error": {
			productGetter: &MockProductGetter{
				products: map[string]product.Product{
					"1": {ID: "1", Name: "Test", PriceCents: 100},
				},
			},
			expectedError: order.ErrCannotCreateOrderService,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			_, actualError := order.NewOrderService(
				tc.orderStore,
				tc.productGetter,
			)

			if tc.expectedError != nil {
				assert.ErrorIsf(
					t,
					actualError,
					tc.expectedError,
					"expected error %q, but got %q",
					tc.expectedError.Error(),
					actualError.Error(),
				)
			} else {
				assert.Nilf(t, actualError, "unexpectedly got error %v", actualError)
			}
		})
	}
}

func TestNewOrder(t *testing.T) {
	testcases := map[string]struct {
		orderStore    order.Store
		productGetter order.ProductGetter
		items         []order.Item
		setupOrders   [][]order.Item
		expectedOrder order.Order
		expectedError error
	}{
		"Single Item Order": {
			orderStore: inmemoryorderdatastore.NewInMemoryOrderStore(),
			productGetter: &MockProductGetter{
				products: map[string]product.Product{
					"1": {ID: "1", Name: "Test", PriceCents: 100},
				},
			},
			items: []order.Item{{ProductID: "1", Quantity: 1}},
		},
		"Multi Item Order": {
			orderStore: inmemoryorderdatastore.NewInMemoryOrderStore(),
			productGetter: &MockProductGetter{
				products: map[string]product.Product{
					"1": {ID: "1", Name: "Test1", PriceCents: 100},
					"2": {ID: "2", Name: "Test2", PriceCents: 200},
					"3": {ID: "3", Name: "Test3", PriceCents: 300},
					"4": {ID: "4", Name: "Test4", PriceCents: 400},
					"5": {ID: "5", Name: "Test5", PriceCents: 500},
					"6": {ID: "6", Name: "Test6", PriceCents: 600},
				},
			},
			items: []order.Item{
				{ProductID: "1", Quantity: 2},
				{ProductID: "4", Quantity: 5},
				{ProductID: "2", Quantity: 3},
				{ProductID: "5", Quantity: 100},
			},
		},
		"Zero Item Order should fail": {
			orderStore: inmemoryorderdatastore.NewInMemoryOrderStore(),
			productGetter: &MockProductGetter{
				products: map[string]product.Product{
					"1": {ID: "1", Name: "Test", PriceCents: 100},
				},
			},
			expectedError: order.ErrCreateFailed,
		},
		"No such product ordered": {
			orderStore: inmemoryorderdatastore.NewInMemoryOrderStore(),
			productGetter: &MockProductGetter{
				products: map[string]product.Product{
					"1": {ID: "1", Name: "Test", PriceCents: 100},
				},
			},
			items:         []order.Item{{ProductID: "10", Quantity: 1}},
			expectedError: order.ErrCreateFailed,
		},
		"Respository error during save": {
			orderStore: &MockOrderStore{
				err: fmt.Errorf("Mocked error"),
			},
			productGetter: &MockProductGetter{
				products: map[string]product.Product{
					"1": {ID: "1", Name: "Test", PriceCents: 100},
				},
			},
			items: []order.Item{{ProductID: "1", Quantity: 1}},

			expectedError: order.ErrCreateFailed,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			nos, err := order.NewOrderService(
				tc.orderStore,
				tc.productGetter,
			)
			assert.Nil(t, err)

			_, actualError := nos.NewOrder(tc.items)

			if tc.expectedError != nil {
				assert.ErrorIsf(
					t,
					actualError,
					tc.expectedError,
					"expected error %q, but got %q",
					tc.expectedError.Error(),
					actualError.Error(),
				)
			} else {
				assert.Nilf(t, actualError, "unexpectedly got error %v", actualError)
			}
		})
	}
}
