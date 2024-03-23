package main

import (
	"html/template"
	"io"
	"io/fs"
)

type Templater struct {
	Templates fs.FS
}

func NewTemplater(templateFS fs.FS) Templater {
	return Templater{Templates: templateFS}
}

func (t Templater) renderTemplateErr(w io.Writer, name string, data any) error {
	var err error

	templateFile := "templates/" + name + ".gohtml"
	tmpl, err := template.New("base").ParseFS(t.Templates, "templates/base.gohtml", templateFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (t Templater) render(w io.Writer, name string, data any) {
	err := t.renderTemplateErr(w, name, data)
	if err != nil {
		t.renderError(w)
	}
}

func (t Templater) renderError(w io.Writer) {
	_ = t.renderTemplateErr(w, "errors/500", nil)
}
