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

// Field 结构体属性
type Field struct {
	Name string
	Type string
	Tags map[string]string
}

func (field *Field) IsID() bool {
	bsonTag, bsonDeclared := field.Tags["bson"]
	if !bsonDeclared {
		return false
	}
	return strings.Contains(bsonTag, "_id")
}

// ToCode 生成代码
func (field *Field) ToCode() string {
	codeTemplate := "%s\t%s\t"
	code := fmt.Sprintf(codeTemplate, field.Name, field.Type)
	tagStrings := []string{}
	for key, val := range field.Tags {
		if key == "" {
			continue
		}
		tagString := fmt.Sprintf("%s:\"%s\"", key, val)
		tagStrings = append(tagStrings, tagString)
	}
	sort.Sort(ByFirstLetter(tagStrings))
	tagsString := strings.Join(tagStrings, "\t")
	tagsString = fmt.Sprintf("`%s`", tagsString)
	code = code + tagsString
	return code
}

type ByFirstLetter []string

func (a ByFirstLetter) Len() int           { return len(a) }
func (a ByFirstLetter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFirstLetter) Less(i, j int) bool { return a[i][0] < a[j][0] }

// FieldFromString 根据属性字符串生成属性对象
func FieldFromString(str string) (Field, error) {
	field := Field{}
	field.Tags = map[string]string{}
	str = strings.TrimSpace(str)
	re := regexp.MustCompile("\\s+")

	spaceLocation := re.FindStringIndex(str)
	if spaceLocation == nil {
		return field, errors.New("syntax error around: " + str)
	}
	field.Name = str[:spaceLocation[0]]
	str = str[spaceLocation[1]:]

	spaceLocation = re.FindStringIndex(str)
	if spaceLocation == nil {
		field.Type = str
		str = ""
	} else {
		field.Type = str[:spaceLocation[0]]
		str = str[spaceLocation[1]:]
	}

	tagsString := str
	tagsString = strings.Trim(tagsString, "`")
	tagStrings := RegSplit(tagsString, "\\s+")
	for _, tagString := range tagStrings {
		tag := strings.Split(tagString, ":")
		if len(tag) < 2 {
			tag = append(tag, "")
		}
		tagName := tag[0]
		tagValue := strings.Trim(tag[1], "\"")
		field.Tags[tagName] = tagValue
	}
	defaultTagNames := []string{
		"json",
		"bson",
	}
	for _, defaultTagName := range defaultTagNames {
		_, declared := field.Tags[defaultTagName]
		if declared == false {
			field.Tags[defaultTagName] = field.Name
		}
	}
	return field, nil
}

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
		fmt.Println("linesLen < 2")
		return modelMeta, errors.New("syntax error around: " + str)
	}
	firstLine := lines[0]
	splitedFirstLine := RegSplit(firstLine, "\\s+")
	if len(splitedFirstLine) < 3 {
		fmt.Println("len(splitedFirstLine) < 3")
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
			fmt.Println("parse field error")
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

func (modelMeta *ModelMeta) GenCode() string {
	t := template.Must(template.New("model").Parse(templates.ModelTmpl))
	var doc bytes.Buffer
	err := t.Execute(&doc, modelMeta)
	if err != nil {
		return "error " + modelMeta.Name
	}
	return doc.String()
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
