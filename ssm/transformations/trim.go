package transformations

import (
	"strings"
)

func TrimKeys(parameters map[string]string, trim string, startswith string) error {

	for key, value := range parameters {
		if strings.HasPrefix(key, startswith) {
			parameters[strings.TrimPrefix(key, trim)] = value
			delete(parameters, key)
		}
	}

	return nil

}
