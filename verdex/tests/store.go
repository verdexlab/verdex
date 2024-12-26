package tests

var testCases = make(map[string][]*TestCase, 0)

func GetProductTestCases(productID string) []*TestCase {
	if _, hasProduct := testCases[productID]; hasProduct {
		return testCases[productID]
	}

	return []*TestCase{}
}

func ClearTestCases() {
	testCases = make(map[string][]*TestCase, 0)
}
