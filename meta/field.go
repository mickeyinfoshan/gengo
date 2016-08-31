package meta

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
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
	sort.Sort(ByFirstLetter(tagStrings))
	tagsString := strings.Join(tagStrings, "\t")
	tagsString = fmt.Sprintf("`%s`", tagsString)
	code = code + tagsString
	return code
}

// GenBsonM 生成形如 "xxx":xxx.xxx的代码
func (field *Field) GenBsonMPair(instanceName string) string {
	mongoAttrName := field.GetMongoAttrName()
	return fmt.Sprintf("\"%s\" : %s.%s", mongoAttrName, instanceName, field.Name)
}

// GetMongoAttrName 获取mongo属性名
func (field *Field) GetMongoAttrName() string {
	bsonTag := field.Tags["bson"]
	mongoAttrName := strings.Split(bsonTag, ",")[0]
	return mongoAttrName
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
