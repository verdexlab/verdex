package products

import "fmt"

type ProductCpe struct {
	Vendor  string `yaml:"vendor"`
	Product string `yaml:"product"`
	Type    string `yaml:"type" validate:"required,oneof=a h o"`
}

// Generate CPE code from product version
func (productCpe *ProductCpe) Build(version string) string {
	return fmt.Sprintf("cpe:2.3:%s:%s:%s:%s:*:*:*:*:*:*:*", productCpe.Type, productCpe.Vendor, productCpe.Product, version)
}
