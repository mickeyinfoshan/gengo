package meta

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/mickeyinfoshan/gengo/templates"
)

// 根据名字排序
type ByName []Field

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name[0] < a[j].Name[0] }

// ModelMeta 结构体解析结果
type ModelMeta struct {
	Name   string
	Type   string
	Fields []Field
}

// ModelMetaFromString 解析结构体
func ModelMetaFromString(str string) (ModelMeta, error) {
	modelMeta := ModelMeta{}
	modelMeta.Fields = []Field{}
	lines := strings.Split(str, "\n")
	linesLen := len(lines)
	if linesLen < 2 {
		return modelMeta, errors.New("syntax error around: " + str)
	}
	firstLine := lines[0]
	splitedFirstLine := RegSplit(firstLine, "\\s+")
	if len(splitedFirstLine) < 3 {
		return modelMeta, errors.New("syntax error around: " + firstLine)
	}
	modelMeta.Name = splitedFirstLine[1]
	modelMeta.Type = splitedFirstLine[2]

	fieldLines := lines[1 : linesLen-1]
	for _, fieldLine := range fieldLines {
		if fieldLine == "" {
			continue
		}
		field, err := FieldFromString(fieldLine)
		if err != nil {
			return modelMeta, err
		}
		modelMeta.Fields = append(modelMeta.Fields, field)
	}

	if !modelMeta.HasIDField() {
		IDField := Field{}
		IDField.Name = modelMeta.Name + "ID"
		IDField.Type = "bson.ObjectId"
		IDField.Tags = map[string]string{}
		IDField.Tags["bson"] = "_id,omitempty"
		IDField.Tags["json"] = modelMeta.Name + "ID"
		modelMeta.Fields = append(modelMeta.Fields, IDField)
	}

	sort.Sort(ByName(modelMeta.Fields))

	return modelMeta, nil
}

func (modelMeta *ModelMeta) GetIDField() (Field, error) {
	var IDField Field
	for _, field := range modelMeta.Fields {
		if field.IsID() {
			return field, nil
		}
	}
	return IDField, errors.New("ID Field not found")
}

func (modelMeta *ModelMeta) HasIDField() bool {
	for _, field := range modelMeta.Fields {
		if field.IsID() {
			return true
		}
	}
	return false
}

func (modelMeta *ModelMeta) GetInstanceName() string {
	typeName := modelMeta.Name
	a := []rune(typeName)
	a[0] = unicode.ToLower(a[0])
	instanceName := string(a)
	return instanceName
}

func (modelMeta *ModelMeta) GenModelCode() string {
	t := template.Must(template.New("model").Parse(templates.ModelTmpl))
	var doc bytes.Buffer
	err := t.Execute(&doc, modelMeta)
	if err != nil {
		return "error " + modelMeta.Name
	}
	return doc.String()
}

func (modelMeta *ModelMeta) GetUpsertFields() []Field {
	fields := FilterFieldsByTags(modelMeta.Fields, codeGenerationTag, upsertNote)
	if len(fields) == 0 {
		IDField, _ := modelMeta.GetIDField()
		fields = append(fields, IDField)
	}
	return fields
}

func (modelMeta *ModelMeta) GenUpsertSelector() string {
	fields := modelMeta.GetUpsertFields()
	instanceName := modelMeta.GetInstanceName()
	fieldStrings := []string{}
	for _, field := range fields {
		fieldStrings = append(fieldStrings, field.GenBsonMPair(instanceName))
	}
	return fmt.Sprintf("bson.M{%s}", strings.Join(fieldStrings, ","))
}

// MakeSelectorMeta Make a selectorMeta from a modelMeta
func (modelMeta *ModelMeta) MakeSelectorMeta() SelectorMeta {
	selectorMeta := SelectorMeta{
		NoneOmitEmptyFields: []Field{},
		OmitEmptyFields:     []Field{},
	}
	selectorMeta.ModelMeta = modelMeta
	filterFields := FilterFieldsByTags(modelMeta.Fields, codeGenerationTag, filterNote)
	for _, field := range filterFields {
		tag := field.Tags[codeGenerationTag]
		if strings.Contains(tag, "omitempty") {
			selectorMeta.OmitEmptyFields = append(selectorMeta.OmitEmptyFields, field)
		} else {
			selectorMeta.NoneOmitEmptyFields = append(selectorMeta.NoneOmitEmptyFields, field)
		}
	}
	return selectorMeta
}

func (modelMeta *ModelMeta) GenSeletorCode() string {
	selectorMeta := modelMeta.MakeSelectorMeta()
	codeGenerated := selectorMeta.genCode()
	return codeGenerated
}

func (modelMeta *ModelMeta) GenFileCode(packageName string) string {
	modelFile := ModelFile{
		ModelMeta:   modelMeta,
		PackageName: packageName,
	}
	codeGenerated := modelFile.GenCode()
	return codeGenerated
}

func FilterFieldsByTags(fieldsToFilter []Field, tagName, tagSubString string) []Field {
	var fields []Field
	for _, field := range fieldsToFilter {
		tag, hasTag := field.Tags[tagName]
		if !hasTag {
			continue
		}
		if strings.Contains(tag, tagSubString) {
			fields = append(fields, field)
		}
	}
	return fields
}

// RegSplit 使用正则表达式分割字符串
func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

type ModelFile struct {
	PackageName string
	ModelMeta   *ModelMeta
}

func (modelFile *ModelFile) GenCode() string {
	t := template.Must(template.New("modelFile").Parse(templates.ModelFileTmpl))
	var doc bytes.Buffer
	err := t.Execute(&doc, modelFile)
	if err != nil {
		return "error " + modelFile.ModelMeta.Name
	}
	return doc.String()
}
