package meta

import (
	"bytes"
	"text/template"

	"github.com/mickeyinfoshan/gengo/templates"
)

// SelectorMeta meta of selector
type SelectorMeta struct {
	ModelMeta           *ModelMeta
	OmitEmptyFields     []Field
	NoneOmitEmptyFields []Field
}

func (selectorMeta *SelectorMeta) genCode() string {
	t := template.Must(template.New("selector").Parse(templates.SelectorTmpl))
	var doc bytes.Buffer
	err := t.Execute(&doc, selectorMeta)
	if err != nil {
		return "error " + selectorMeta.ModelMeta.Name
	}
	return doc.String()
}
