package datastore_test

import (
	"testing"

	"github.com/shanehowearth/kart/product"
	"github.com/shanehowearth/kart/product/datastore"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	testcases := map[string]struct {
		id              string
		expectedProduct product.Product
		expectedError   error
	}{
		"Successfully fetch product 2": {
			id:              "2",
			expectedProduct: datastore.SeedProducts[2],
		},
		"Fail to fetch non-existant product": {
			id:              "does-not-exist",
			expectedProduct: product.Product{},
			expectedError:   product.ErrNotFound,
		},
	}
	for name, tc := range testcases { //nolint:varnamelen // tc is fine in a test.
		t.Run(name, func(t *testing.T) {
			imps := datastore.NewSeededInMemoryProductStore()
			actualProduct, actualError := imps.GetByID(tc.id)

			if tc.expectedError != nil {
				assert.ErrorIsf(
					t,
					tc.expectedError,
					actualError,
					"expected error %v, but got %v",
					tc.expectedError, actualError,
				)

				assert.EqualExportedValuesf(
					t,
					tc.expectedProduct,
					actualProduct,
					"expected product to match %#v but got %#v",
					tc.expectedProduct, actualProduct,
				)
			}
		})
	}
}

func TestList(t *testing.T) {
	testcases := map[string]struct {
		imps             *datastore.InMemoryProductStore
		expectedProducts []product.Product
	}{
		"Seeded IMPS returns list of 10 products, in correct order": {
			imps:             datastore.NewSeededInMemoryProductStore(),
			expectedProducts: datastore.SeedProducts,
		},
		"Empty IMPS returns empty list": {
			imps: datastore.NewInMemoryProductStore(),
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actualProducts := tc.imps.List()
			assert.ElementsMatch(t, tc.expectedProducts, actualProducts)
		})
	}
}
