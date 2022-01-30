package ssm

import "testing"

func TestExpandPresentNotPresent(t *testing.T) {
	parameters := make(map[string]string)

	if err := expandParameters(parameters, false, []string{"test"}); err == nil {
		t.Errorf("expected error when supplying var which is not present in params, but got nil")
	}

	// set a value
	parameters["test"] = "test value"

	if err := expandParameters(parameters, false, []string{"test"}); err != nil {
		t.Errorf("expected no error, but got: %s", err)
	}
}

func TestExpandSelectiveExpansions(t *testing.T) {
	t.Setenv("ENVIRONMENT", "teststaging")

	parameters := map[string]string{
		"DATABASE_NAME": "DB_$ENVIRONMENT", // should be expanded
		"SOME_SECRET":   "abc$abc",         // should not be expanded
	}

	// don't want to expand all here, just specific vars
	if err := expandParameters(parameters, false, []string{"DATABASE_NAME"}); err != nil {
		t.Errorf("expected no error, but got: %s", err)
	}

	if parameters["DATABASE_NAME"] != "DB_teststaging" {
		t.Errorf("DATABASE_NAME should be expanded to 'DB_teststaging', but got '%s'", parameters["DATABASE_NAME"])
	}
	if parameters["SOME_SECRET"] != "abc$abc" {
		t.Errorf("SOME_SECRET should not be expanded and be 'abc$abc', but got '%s'", parameters["SOME_SECRET"])
	}
}
func TestExpandExpansions(t *testing.T) {
	t.Setenv("ENVIRONMENT", "teststaging")
	t.Setenv("abc", "def")

	parameters := map[string]string{
		"DATABASE_NAME": "DB_$ENVIRONMENT", // should be expanded
		"SOME_SECRET":   "abc$abc",         // should not be expanded
	}
	want := map[string]string{
		"DATABASE_NAME": "DB_teststaging",
		"SOME_SECRET":   "abcdef",
	}

	// want to expand all
	if err := expandParameters(parameters, true, []string{}); err != nil {
		t.Errorf("expected no error, but got: %s", err)
	}

	for key := range want {
		if parameters[key] != want[key] {
			t.Errorf("%s should be expanded to '%s', but got '%s'", key, want[key], parameters[key])
		}
	}

}
