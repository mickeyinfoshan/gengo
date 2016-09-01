package meta

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldFromString(t *testing.T) {

	// testCase2 := "Age int `json:\"-\" bson:\"createtime,omitempty\"`"
	// testCase3 := "      Aid        bson.ObjectId `json:\"aid\" bson:\"_id,omitempty\"`"
	Convey("使用默认tags", t, func() {
		testCase := "Name string"
		field, err := FieldFromString(testCase)
		So(err, ShouldBeNil)
		So(field.Name, ShouldEqual, "Name")
		So(field.Type, ShouldEqual, "string")
		So(field.IsID(), ShouldBeFalse)
		So(field.Tags["json"], ShouldEqual, "Name")
		So(field.Tags["bson"], ShouldEqual, "Name")
		So(field.ToCode(), ShouldEqual, "Name\tstring\t`bson:\"Name\"\tjson:\"Name\"`")
	})

	Convey("不使用默认tags", t, func() {
		testCase := "Age int `json:\"-\" bson:\"createtime,omitempty\"`"
		field, err := FieldFromString(testCase)
		So(err, ShouldBeNil)
		So(field.Name, ShouldEqual, "Age")
		So(field.Type, ShouldEqual, "int")
		So(field.IsID(), ShouldBeFalse)
		So(field.Tags["json"], ShouldEqual, "-")
		So(field.Tags["bson"], ShouldEqual, "createtime,omitempty")
		So(field.ToCode(), ShouldEqual, "Age\tint\t`bson:\"createtime,omitempty\"\tjson:\"-\"`")
	})

	Convey("不使用默认tags，ID", t, func() {
		testCase := "      Aid        bson.ObjectId  bson:\"_id,omitempty\"`"
		field, err := FieldFromString(testCase)
		So(err, ShouldBeNil)
		So(field.Name, ShouldEqual, "Aid")
		So(field.Type, ShouldEqual, "bson.ObjectId")
		So(field.IsID(), ShouldBeTrue)
		So(field.Tags["json"], ShouldEqual, "Aid")
		So(field.Tags["bson"], ShouldEqual, "_id,omitempty")
	})
}

func TestModelMetaFromString(t *testing.T) {
	pwd, _ := os.Getwd()

	Convey("测试解析器，有ID属性", t, func() {
		testfileBytes, err := ioutil.ReadFile(pwd + "/../tests/case1.idl")
		if err != nil {
			panic(err)
		}
		testFileStr := string(testfileBytes)
		modelMeta, e := ModelMetaFromString(testFileStr)
		So(e, ShouldBeNil)
		So(modelMeta.Name, ShouldEqual, "AuthTest")
		So(modelMeta.Type, ShouldEqual, "struct")
		So(len(modelMeta.Fields), ShouldEqual, 11)
		So(modelMeta.GetInstanceName(), ShouldEqual, "authTest")
	})

	Convey("测试解析器，没有ID属性", t, func() {
		testfileBytes, err := ioutil.ReadFile(pwd + "/../tests/case2.idl")
		if err != nil {
			panic(err)
		}
		resultBytes, err := ioutil.ReadFile(pwd + "/../tests/case2Model.go.generated")
		if err != nil {
			panic(err)
		}
		testFileStr := string(testfileBytes)
		resultCode := string(resultBytes)
		resultCode = strings.TrimSpace(resultCode)
		resultCode = strings.Trim(resultCode, "\n")
		modelMeta, e := ModelMetaFromString(testFileStr)
		So(e, ShouldBeNil)
		So(modelMeta.Name, ShouldEqual, "Auth")
		So(modelMeta.Type, ShouldEqual, "struct")
		So(len(modelMeta.Fields), ShouldEqual, 11)
		So(modelMeta.HasIDField(), ShouldBeTrue)
		So(modelMeta.GetInstanceName(), ShouldEqual, "auth")

		generatedCode := modelMeta.GenModelCode()
		generatedCode = strings.Trim(generatedCode, "\n")
		generatedCode = strings.TrimSpace(generatedCode)

		So(generatedCode, ShouldEqual, resultCode)
	})
}
