package wocab

import (
	"bytes"
	"github.com/pkg/errors"
	"html/template"
)

func renderFragement(page string, data interface{}) (template.HTML, error) {
	// load template files

	templates := []string{
		page,
	}

	// compile template

	ts, err := template.ParseFiles(templates...)

	if err != nil {
		return "", errors.Wrapf(err, "parsing templates page=%v failed", page)
	}

	// render html

	buf := bytes.Buffer{}

	if ts.Execute(&buf, data); err != nil {
		return "", errors.Wrapf(err, "executing template page=%v failed", page)
	}

	// convert to string

	return template.HTML(buf.Bytes()), nil
}
