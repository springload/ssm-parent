package transformations

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/apex/log"
)

type Transformation interface {
	Transform(map[string]string) (map[string]string, error)
}

// DeleteTransformation takes an input map and removes the specified keys
type DeleteTransformation struct {
	Action string
	Rule   []string
}

func (t *DeleteTransformation) Transform(source map[string]string) (map[string]string, error) {
	for _, key := range t.Rule {
		delete(source, key)
	}

	return source, nil
}

// RenameTransformation renames the keys in the input map
type RenameTransformation struct {
	Action string
	Rule   map[string]string
}

func (t *RenameTransformation) Transform(source map[string]string) (map[string]string, error) {
	for key, newkey := range t.Rule {
		if _, ok := source[key]; ok {
			log.Debugf("Renaming '%s' to '%s'", key, newkey)
			source[newkey] = source[key]
			delete(source, key)
		} else {
			log.Warnf("Not renaming '%s' to '%s' as it isn't set", key, newkey)
		}
	}

	return source, nil
}

// TemlateTransformation renames the keys in the input map
type TemplateTransformation struct {
	Action string
	Rule   map[string]string
}

func (t *TemplateTransformation) Transform(source map[string]string) (map[string]string, error) {
	for key, value := range t.Rule {
		tmpl, err := template.New("value").Funcs(funcMap).Parse(value)
		if err != nil {
			return source, fmt.Errorf("can't template '%s': %s", key, err)
		}
		var renderedValue strings.Builder
		err = tmpl.Execute(&renderedValue, source)
		if err != nil {
			return source, fmt.Errorf("can't execute template '%s': %s", key, err)
		}
		// if all good, then assign the rendered value
		source[key] = renderedValue.String()
	}

	return source, nil
}

// TrimTransformation modifies the keys in the input map
type TrimTransformation struct {
	Action string
	Rule   map[string]string
}

func (t *TrimTransformation) Transform(source map[string]string) (map[string]string, error) {
	if _, found := t.Rule["trim"]; !found {
		return source, fmt.Errorf("\"trim\" rule not set")
	}
	if _, found := t.Rule["starts_with"]; !found {
		return source, fmt.Errorf("\"starts_with\" rule not set")
	}
	TrimKeys(source, t.Rule["trim"], t.Rule["starts_with"])

	return source, nil
}
