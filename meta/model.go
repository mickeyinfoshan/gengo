package meta

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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
	tagsString := strings.Join(tagStrings, "\t")
	tagsString = fmt.Sprintf("`%s`", tagsString)
	code = code + tagsString
	return code
}

// FieldFromString 根据属性字符串生成属性对象
func FieldFromString(str string) (Field, error) {
	var field Field
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

// ModelMeta 结构体解析结果
type ModelMeta struct {
	Name   string
	Fields []Field
}

// ModelMetaFromString 解析结构体
// func ModelMetaFromString(str string) (ModelMeta, error) {

// }

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
