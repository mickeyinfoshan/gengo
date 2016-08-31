package meta

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSelectorGenCode(t *testing.T) {
	pwd, _ := os.Getwd()
	Convey("选择器生成代码", t, func() {
		testfileBytes, err := ioutil.ReadFile(pwd + "/../tests/case2.idl")
		if err != nil {
			panic(err)
		}
		resultBytes, err := ioutil.ReadFile(pwd + "/../tests/case2Selector.go.generated")
		if err != nil {
			panic(err)
		}
		testFileStr := string(testfileBytes)
		resultCode := string(resultBytes)
		resultCode = strings.TrimSpace(resultCode)
		resultCode = strings.Trim(resultCode, "\n")
		modelMeta, e := ModelMetaFromString(testFileStr)
		if e != nil {
			panic(e)
		}
		selectorMeta := modelMeta.MakeSelectorMeta()
		generatedCode := selectorMeta.genCode()
		generatedCode = strings.Trim(generatedCode, "\n")
		generatedCode = strings.TrimSpace(generatedCode)
		So(generatedCode, ShouldEqual, resultCode)
	})
}
