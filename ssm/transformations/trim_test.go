package transformations

import (
	"testing"
)

func TestTrimKeys(t *testing.T) {
	t.Setenv("ENVIRONMENT", "teststaging")

	parameters := map[string]string{
		"_PREFIXED_PARAMETER": "prefixed_value",
		"_ANOTHER_PARAMETER":  "another_value",
	}

	expecting := map[string]string{
		"prefixed_value": "PREFIXED_PARAMETER",
		"another_value": "_ANOTHER_PARAMETER",
	}

	TrimKeys(parameters, "_", "_PREFIXED_")

	// Swap result keys and values to check against expectations
	result := make(map[string]string)
	for key, value := range parameters {
		result[value] = key
	}

	for val, expectedParam := range expecting {
		if result[val] != expectedParam {
			t.Errorf("'%s's key should be '%s', but got '%s'", val, expectedParam, result[val])
		}
	}
}
