package products

type Product struct {
	ID          string            `yaml:"-"`
	Name        string            `yaml:"name" validate:"required"`
	Description string            `yaml:"description" validate:"required"`
	Website     string            `yaml:"website" validate:"required"`
	Versions    ProductVersions   `yaml:"versions" validate:"required"`
	Cpe         ProductCpe        `yaml:"cpe"`
	SmokeTests  ProductSmokeTests `yaml:"smoke-tests"`
}
