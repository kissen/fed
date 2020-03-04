package main

import (
	"bytes"
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
	"gitlab.cs.fau.de/kissen/fed/fedutil"
	"html/template"
)

type Card struct {
	Payload vocab.Type
}

func (c Card) HTML() (string, error) {
	switch v := c.Payload.(type) {
	case vocab.ActivityStreamsPerson:
		return c.renderPerson()
	default:
		return "", fmt.Errorf("unsupported type %v", v.GetTypeName())
	}
}

func (c Card) renderPerson() (string, error) {
	data := map[string]interface{}{
		"Id": fedutil.Id(c.Payload).String(),
	}

	return c.render("person.fragment.tmpl", data)
}

func (c Card) render(cardTemplate string, data map[string]interface{}) (string, error) {
	// compile template

	templates := []string{
		cardTemplate, "card.layout.tmpl",
	}

	ts, err := template.ParseFiles(templates...)

	if err != nil {
		return "", errors.Wrap(err, "parsing templates failed")
	}

	// write to buffer

	var buf bytes.Buffer

	if ts.Execute(&buf, data); err != nil {
		return "", errors.Wrap(err, "executing template failed")
	}

	// return as string

	return buf.String(), nil
}
