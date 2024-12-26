package variables

var variables = make(map[string]map[string]*Variable, 0)
var variablesKeys = make(map[string][]string, 0)

func GetAllProductVariables(product string) []string {
	if _, hasProduct := variables[product]; hasProduct {
		return variablesKeys[product]
	}

	return nil
}

func GetProductVariable(product string, variableKey string) *Variable {
	if _, hasProduct := variables[product]; hasProduct {
		if _, hasKey := variables[product][variableKey]; hasKey {
			return variables[product][variableKey]
		}
	}

	return nil
}

func ClearVariables() {
	variables = make(map[string]map[string]*Variable, 0)
	variablesKeys = make(map[string][]string, 0)
}
