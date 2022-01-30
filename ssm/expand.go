package ssm

import (
	"fmt"
	"strings"

	"github.com/buildkite/interpolate"
)

// expandArgs expands arguments using env vars
func ExpandArgs(args []string) []string {
	var expanded []string
	for _, arg := range args {
		arg = expandValue(arg)
		expanded = append(expanded, arg)
	}
	return expanded
}

// expandValue interpolates values using env vars
func expandValue(val string) string {
	e, err := interpolate.Interpolate(env, val)
	if err == nil {
		return strings.TrimSpace(string(e))
	}
	return val
}

// expandParameters expands values using shell-like syntax
func expandParameters(parameters map[string]string, expand bool, expandValues []string) error {

	// if global expand is true then just it for all
	if expand {
		for key, value := range parameters {
			parameters[key] = expandValue(value)
		}
		// can return early as we've done the job
		return nil
	}
	// check if all values that we ask to expand present in the parameters
	// otherwise, it's a configuration error
	for _, val := range expandValues {
		if _, ok := parameters[val]; !ok {
			return fmt.Errorf("env var %s is present in the expand-values but doesn't exist in the environment", val)
		} else {
			// if the var is present we expand it
			parameters[val] = expandValue(parameters[val])
		}
	}

	return nil
}
