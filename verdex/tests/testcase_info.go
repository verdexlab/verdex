package tests

type TestCaseInfo struct {
	Product string `yaml:"product" validate:"required"`
}
