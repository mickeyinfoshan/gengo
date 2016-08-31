package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mickeyinfoshan/gengo/meta"
)

var (
	structIndexItem = regexp.MustCompile(`(?s)type .+? struct(.*?)\}`)
)

func main() {
	var dbURL string
	var dbName string
	var path string
	var packageName string
	var outputPath string

	flag.StringVar(&dbURL, "url", "127.0.0.1:27017", "DB url")
	flag.StringVar(&dbName, "db", "", "DB Name")
	flag.StringVar(&packageName, "pname", "models", "package name")
	flag.StringVar(&outputPath, "opath", "", "out path")
	flag.StringVar(&path, "path", "", "file path")
	flag.Parse()

	if len(path) == 0 {
		path = CurPath()
	}
	if len(outputPath) == 0 {
		outputPath = CurPath()
	}

	modelMetas := readFiles(path)
	if len(modelMetas) > 0 {
		genCodeContext := GenCodeContext{
			DatabaseURL:  dbURL,
			DatabaseName: dbName,
			PackageName:  packageName,
			OutputPath:   outputPath,
			ModelMetas:   modelMetas,
		}

		genCodeContext.Execute()
	}
}

func readFiles(dir string) []meta.ModelMeta {
	metas := []meta.ModelMeta{}

	fi, err := os.Stat(dir)
	if err != nil {
		return metas
	} else if fi.IsDir() {
		files, e := ioutil.ReadDir(dir)
		if e != nil {
			return metas
		}
		for _, file := range files {
			//遍历文件夹内所有文件
			if file.IsDir() {
				mts := readFiles(dir + SystemSep() + file.Name())
				if len(mts) > 0 {
					metas = append(metas, mts...)
				}
			} else {
				mts := makeModelMeta(dir + SystemSep() + file.Name())
				if mts != nil {
					metas = append(metas, mts...)
				}
			}
		}
	} else {
		mts := makeModelMeta(dir)
		if mts != nil {
			metas = append(metas, mts...)
		}
	}

	return metas
}

func makeModelMeta(filename string) []meta.ModelMeta {
	metas := []meta.ModelMeta{}
	if !strings.HasSuffix(filename, ".go") {
		return metas
	}
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		return metas
	}

	matches := structIndexItem.FindAllStringSubmatch(string(data), -1)

	for _, item := range matches {
		if len(item) == 2 {
			mt, e := meta.ModelMetaFromString(item[0])
			if e != nil {
				continue
			}
			metas = append(metas, mt)
		}
	}

	return metas
}

// CurPath 获取当前运行目录
func CurPath() (path string) {
	file, _ := exec.LookPath(os.Args[0])
	pt, _ := filepath.Abs(file)

	return filepath.Dir(pt)
}

// SystemSep 获取系统分隔符
func SystemSep() (path string) {

	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	return path
}
