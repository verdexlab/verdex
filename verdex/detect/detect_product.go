package detect

import (
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/products"
)

// Try to detect product of given target with smoke tests
func DetectProduct(execution *core.Execution, detection *core.Detection) (detectedProduct *products.Product) {
	for _, product := range products.ListProducts() {
		execution.Product = product.ID

		if product.SmokeTests.DetectProduct(execution, detection, product) {
			return product
		}
	}

	return nil
}
