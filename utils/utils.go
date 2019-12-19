package utils

import (
	"bytes"
	"regexp"
	"text/template"
)

// process applies the data structure 'vars' onto an already
// parsed template 't', and returns the resulting string.
func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

func ProcessString(str string, vars interface{}) string {
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}

func GetHashFromMagnet(magnetURI string) string {
	re := regexp.MustCompile(`xt=urn:btih:(?P<hash>[^&/]+)`)
	hash := re.FindStringSubmatch(magnetURI)[1]
	return hash
}
