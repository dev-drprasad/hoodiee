package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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

func Respond(w http.ResponseWriter, statusCode int, data interface{}, err error) error {
	w.WriteHeader(statusCode)

	var r map[string]interface{}
	if err != nil {
		r = map[string]interface{}{"data": nil, "error": err}
	} else {
		r = map[string]interface{}{"data": data, "error": nil}
	}

	json.NewEncoder(w).Encode(r)

	log.Printf("Response: StatusCode=%d Error=%s", statusCode, err)
	return err
}
