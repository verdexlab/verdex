package rules

var Rules = make(map[string][]*Rule, 0)

func GetProductRules(product string) []*Rule {
	if _, hasProduct := Rules[product]; hasProduct {
		return Rules[product]
	}

	return nil
}

func ClearRules() {
	Rules = make(map[string][]*Rule, 0)
}
