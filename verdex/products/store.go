package products

var products = make(map[string]*Product, 0)

func ListProducts() []*Product {
	allProducts := make([]*Product, len(products))

	i := 0
	for _, p := range products {
		allProducts[i] = p
		i++
	}

	return allProducts
}

func ListProductIDs() []string {
	productIDs := make([]string, len(products))

	i := 0
	for productID := range products {
		productIDs[i] = productID
		i++
	}

	return productIDs
}

func GetProduct(productID string) *Product {
	if _, hasProduct := products[productID]; hasProduct {
		return products[productID]
	}

	return nil
}

func ClearProducts() {
	products = make(map[string]*Product, 0)
}
