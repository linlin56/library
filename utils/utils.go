package utils

import (
	"bytes"
	"html/template"
	"io/ioutil"
)

// helper pour parser/exécuter un fragment template (dans content/)
func RenderContentTemplate(filename string, ctx interface{}) (string, error) {
	b, err := ioutil.ReadFile("content/" + filename)
	if err != nil {
		return "", err
	}
	tpl := template.New("fragment").Funcs(template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) },
	})
	tpl, err = tpl.Parse(string(b))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}
