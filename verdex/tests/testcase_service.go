package tests

type TestCaseService struct {
	Name string `yaml:"name" validate:"required"`
	Port int    `yaml:"port" validate:"required,gte=0,lte=65535"`
}
