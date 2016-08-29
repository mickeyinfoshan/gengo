package meta

import (
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
		So(field.ToCode(), ShouldEqual, "Name\tstring\t`json:\"Name\"\tbson:\"Name\"`")
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
		So(field.ToCode(), ShouldEqual, "Age\tint\t`json:\"-\"\tbson:\"createtime,omitempty\"`")
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
