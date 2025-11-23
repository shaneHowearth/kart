//nolint:varnamelen // id and tc are clear enough.
package inmemoryorderdatastore_test

import (
	"testing"

	"github.com/shanehowearth/kart/order"
	"github.com/shanehowearth/kart/order/datastore/inmemoryorderdatastore"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	testcases := map[string]struct {
		setupOrders   []order.Order
		newOrder      order.Order
		expectedError error
	}{
		"Order saves properly": {
			newOrder: order.Order{
				ID: "6cb6e494-30fe-4a7a-9e82-3acb8e28e0de",
				Items: []order.Item{
					{ProductID: "1", Quantity: 5},
					{ProductID: "7", Quantity: 1},
				},
			},
		},
		"Duplicate Order throws error": {
			newOrder: order.Order{
				ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
				Items: []order.Item{
					{ProductID: "3", Quantity: 4},
					{ProductID: "5", Quantity: 3},
				},
			},
			setupOrders: []order.Order{
				{
					ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
				{
					ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
			},
			expectedError: order.ErrCreateFailed,
		},
		"Multiple different orders save properly": {
			newOrder: order.Order{
				ID: "a22fa8a0-772e-415f-a5fb-0016455f159e",
				Items: []order.Item{
					{ProductID: "3", Quantity: 4},
					{ProductID: "5", Quantity: 3},
				},
			},
			setupOrders: []order.Order{
				{
					ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
				{
					ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			imods := inmemoryorderdatastore.NewInMemoryOrderStore()

			// Setup: insert any existing orders.
			for _, setupOrder := range tc.setupOrders {
				_ = imods.CreateOrder(&setupOrder)
			}

			actualError := imods.CreateOrder(&tc.newOrder)

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

func TestGetByID(t *testing.T) {
	testcases := map[string]struct {
		ID            string
		setupOrders   []order.Order
		expectedOrder order.Order
		expectedError error
		callCount     int
	}{
		"Get only order in orders": {
			ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
			setupOrders: []order.Order{
				{
					ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
			},
			expectedOrder: order.Order{
				ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
				Items: []order.Item{
					{ProductID: "3", Quantity: 4},
					{ProductID: "5", Quantity: 3},
				},
			},
		},
		"Attempt to get non-existant order when no orders exist": {
			ID:            "non-existant",
			expectedError: order.ErrNotFound,
		},
		"Attempt to get non-existant order": {
			ID: "non-existant",
			setupOrders: []order.Order{
				{
					ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
				{
					ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
			},
			expectedError: order.ErrNotFound,
		},
		"Get one order of many in orders": {
			ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
			setupOrders: []order.Order{
				{
					ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
				{
					ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
			},
			expectedOrder: order.Order{
				ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
				Items: []order.Item{
					{ProductID: "3", Quantity: 4},
					{ProductID: "5", Quantity: 3},
				},
			},
		},
		"Get one order multiple times": {
			ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
			setupOrders: []order.Order{
				{
					ID: "5b68308d-81ba-4ff4-aea9-0c486f9de220",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
				{
					ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
					Items: []order.Item{
						{ProductID: "3", Quantity: 4},
						{ProductID: "5", Quantity: 3},
					},
				},
			},
			expectedOrder: order.Order{
				ID: "d367dd04-f183-40a9-b325-b55faf9cb18e",
				Items: []order.Item{
					{ProductID: "3", Quantity: 4},
					{ProductID: "5", Quantity: 3},
				},
			},
			callCount: 3,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			imods := inmemoryorderdatastore.NewInMemoryOrderStore()

			// Setup: insert any existing orders.
			for _, setupOrder := range tc.setupOrders {
				_ = imods.CreateOrder(&setupOrder)
			}

			for i := 0; i < tc.callCount+1; i++ {
				actualOrder, actualError := imods.GetByID(tc.ID)
				if tc.expectedError != nil {
					assert.ErrorIsf(
						t,
						actualError,
						tc.expectedError,
						"expected error %q, but got %q",
						tc.expectedError.Error(),
						actualError.Error(),
					)
					assert.Empty(t, actualOrder)
				} else {
					assert.Nilf(t, actualError, "unexpectedly got error %v", actualError)
					assert.EqualExportedValues(t, tc.expectedOrder, actualOrder)
				}
			}
		})
	}
}
