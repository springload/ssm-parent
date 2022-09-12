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

	if err := TrimKeys(parameters, "_", "_PREFIXED_"); err != nil {
		t.Errorf("expected no error, but got: %s", err)
	}

	// Swap result keys and values to check against expectations
	result := make(map[string]string)
	for key, value := range parameters {
		result[value] = key
	}

	for val, expected_param := range expecting {
		if result[val] != expected_param {
			t.Errorf("'%s's key should be '%s', but got '%s'", val, expected_param, result[val])
		}
	}
}
