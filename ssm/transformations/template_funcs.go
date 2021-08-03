package transformations

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"env":          GetEnv,
	"url_host":     URLHost,
	"url_port":     URLPort,
	"url_password": URLPassword,
	"url_path":     URLPath,
	"url_scheme":   URLScheme,
	"url_user":     URLUser,
	"trim_prefix":  strings.TrimPrefix,
	"replace":      strings.Replace,
}

// GetEnv gets the environment variable
func GetEnv(input string) (string, error) {
	val, ok := os.LookupEnv(input)
	if !ok {
		return "", fmt.Errorf("can't find %s in the environment variables", input)
	}
	return val, nil
}

// URLUser extracts user from the URL or returns "" if it's not set
func URLUser(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.User.Username(), nil
}

// URLPassword extracts password from the URL or returns "" if it's not set
func URLPassword(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	p, ok := u.User.Password()
	if !ok {
		return "", nil
	}
	return p, nil
}

// URLScheme extracts password from the URL or returns "" if it's not set
func URLScheme(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.Scheme, nil
}

// URLHost extracts host from the URL or returns "" if it's not set. It also strips any port if there is any
func URLHost(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// URLHost extracts host from the URL or returns "" if it's not set
func URLPort(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.Port(), nil
}

// URLPath extracts path from the URL or returns "" if it's not set
func URLPath(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}
