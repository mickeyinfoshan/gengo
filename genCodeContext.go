package gengo

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"path"

	"github.com/mickeyinfoshan/gengo/meta"
	"github.com/mickeyinfoshan/gengo/templates"
)

// GenCodeContext 生成代码的上下文
type GenCodeContext struct {
	DatabaseURL  string
	DatabaseName string
	PackageName  string
	OutputPath   string
	ModelMetas   []meta.ModelMeta
}

func (genCodeContext *GenCodeContext) GenInitDBCode() string {
	t := template.Must(template.New("initDB").Parse(templates.InitDBTempl))
	var doc bytes.Buffer
	err := t.Execute(&doc, genCodeContext)
	if err != nil {
		return "error "
	}
	return doc.String()
}

func (genCodeContext *GenCodeContext) WriteCodeToFile(code, fileName string) error {
	filePath := path.Join(genCodeContext.OutputPath, genCodeContext.PackageName, fileName)
	err := ioutil.WriteFile(filePath, []byte(code), 0666)
	return err
}

func (genCodeContext *GenCodeContext) Execute() error {
	initDBCode := genCodeContext.GenInitDBCode()
	var err error
	err = genCodeContext.WriteCodeToFile(initDBCode, "initDB.go")
	if err != nil {
		return err
	}
	for _, modelMeta := range genCodeContext.ModelMetas {
		modelFileCode := modelMeta.GenFileCode(genCodeContext.PackageName)
		err = genCodeContext.WriteCodeToFile(modelFileCode, modelMeta.Name+".go")
		if err != nil {
			return err
		}
	}

	return nil
}
