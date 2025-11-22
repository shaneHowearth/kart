package datastore

import "github.com/shanehowearth/kart/product"

// SeedProducts is the default menu data - exported for use in tests.
//
//nolint:mnd // This is test data.
var SeedProducts = []product.Product{
	{
		ID:         "1",
		Name:       "Waffle with Berries",
		PriceCents: 650,
		Category:   "Waffle",
	},
	{
		ID:         "2",
		Name:       "Vanilla Bean Crème Brûlée",
		PriceCents: 700,
		Category:   "Crème Brûlée",
	},
	{
		ID:         "3",
		Name:       "Macaron Mix of Five",
		PriceCents: 800,
		Category:   "Macaron",
	},
	{
		ID:         "4",
		Name:       "Classic Tiramisu",
		Category:   "Tiramisu",
		PriceCents: 550,
	},
	{
		ID:         "5",
		Name:       "Pistachio Baklava",
		Category:   "Baklava",
		PriceCents: 400,
	},
	{
		ID:         "6",
		Name:       "Lemon Meringue Pie",
		Category:   "Pie",
		PriceCents: 500,
	},
	{
		ID:         "7",
		Name:       "Red Velvet Cake",
		Category:   "Cake",
		PriceCents: 450,
	},
	{
		ID:         "8",
		Name:       "Salted Caramel Brownie",
		Category:   "Brownie",
		PriceCents: 450,
	},
	{
		ID:         "9",
		Name:       "Vanilla Panna Cotta",
		Category:   "Panna Cotta",
		PriceCents: 650,
	},
	{
		ID:         "10",
		Name:       "Chicken Waffle",
		PriceCents: 100,
		Category:   "Waffle",
	},
}
