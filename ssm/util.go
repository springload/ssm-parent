package ssm

import (
	"os"
)

var env Env

// difference returns the elements in a that aren't in b
// the second argument is slice of string pointers to suit AWS SDK
func stringSliceDifference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// Env just adapts os.LookupEnv to this interface
type Env struct{}

// Get gets env var by the provided key
func (e Env) Get(key string) (string, bool) {
	return os.LookupEnv(key)
}
