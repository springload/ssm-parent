package transformations

import (
	"strings"
)

func TrimKeys(parameters map[string]string, trim string, startsWith string) {

	for key, value := range parameters {
		if strings.HasPrefix(key, startsWith) {
			parameters[strings.TrimPrefix(key, trim)] = value
			delete(parameters, key)
		}
	}

}
