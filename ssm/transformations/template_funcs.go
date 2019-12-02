package transformations

import (
	"net/url"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"url_host":     URLHost,
	"url_password": URLPassword,
	"url_path":     URLPath,
	"url_scheme":   URLScheme,
	"url_user":     URLUser,
	"trim_prefix":  strings.TrimPrefix,
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

// URLHost extracts host from the URL or returns "" if it's not set
func URLHost(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

// URLPath extracts path from the URL or returns "" if it's not set
func URLPath(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}
