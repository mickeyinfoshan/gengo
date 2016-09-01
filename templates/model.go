package templates

// ModelTmpl 模型结构体
const ModelTmpl = `
// model {{.Name}}
type {{.Name}} {{.Type}} {
    {{range $key,$field := .Fields }}
        {{$field.ToCode}}
    {{end}}
}

{{$instance := .GetInstanceName}}
{{$IDField := .GetIDField}}

// Delete Delete {{.Name}} from database
func ({{$instance}} *{{.Name}}) Delete() error {
    sess := NewDBSession()
    defer sess.Close()
    err := sess.DB(dbName).C(collectionNames["{{.Name}}"]).RemoveId({{$instance}}.{{$IDField.Name}})
    return err
}

// Save Save {{.Name}} to database
func ({{$instance}} *{{.Name}}) Save() error {
    sess := NewDBSession()
    defer sess.Close()
    collection := sess.DB(dbName).C(collectionNames["{{.Name}}"])
	selectorBsonM := {{.GenUpsertSelector}}
	info, err := collection.Upsert(selectorBsonM, *{{$instance}})
	if err == nil {
		if info.UpsertedId != nil {
			{{$instance}}.{{$IDField.Name}} = info.UpsertedId.(bson.ObjectId)
		} else if {{$instance}}.{{$IDField.Name}}.Hex() == "" {
			err = collection.Find(selectorBsonM).One({{$instance}})
		}
	}
    return err
}

// Get{{.Name}}ByID Get an instance from ID
func Get{{.Name}}ByID(id bson.ObjectId) ({{.Name}}, error) {
    sess := NewDBSession()
    defer sess.Close()
    var {{$instance}} {{$.Name}}
    err := sess.DB(dbName).C(collectionNames["{{$.Name}}"]).Find(bson.M{"_id" : id}).One(&{{$instance}})
    return {{$instance}}, err
} 
`
const ModelFileTmpl = `
package {{.PackageName}}

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

{{.ModelMeta.GenModelCode}}

{{.ModelMeta.GenSeletorCode}}
`
