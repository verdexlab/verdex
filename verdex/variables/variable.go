package variables

type Variable struct {
	Info    VariableInfo    `yaml:"info"`
	Resolve VariableResolve `yaml:"resolve"`
	Key     string          `yaml:"-"`
}

type VariableInfo struct {
	Product string `yaml:"product" validate:"required"`
	Author  string `yaml:"author"`
}

type VariableResolve struct {
	Type   VariableResolveType `yaml:"type" validate:"required,oneof=regex"`
	Method VariableResolveType `yaml:"method" validate:"required,oneof=GET"`
	Path   string              `yaml:"path" validate:"required"`
	Part   VariableResolvePart `yaml:"part" validate:"required,oneof=body"`
	Regex  string              `yaml:"regex" validate:"required"`
	Group  int                 `yaml:"group" validate:"required,gte=1"`
}

type VariableResolveType string

const (
	VariableResolveTypeRegex VariableResolveType = "regex"
)

type VariableResolveMethod string

const (
	VariableResolveMethodGet VariableResolveMethod = "GET"
)

type VariableResolvePart string

const (
	VariableResolvePartBody VariableResolvePart = "body"
)
