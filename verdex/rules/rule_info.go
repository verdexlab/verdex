package rules

type RuleInfo struct {
	Product string `yaml:"product" validate:"required"`
	Author  string `yaml:"author"`
}
