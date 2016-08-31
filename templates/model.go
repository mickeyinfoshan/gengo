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
    info, err := sess.DB(dbName).C(collectionNames["{{.Name}}"]).Upsert({{.GenUpsertSelector}}, *{{$instance}})
    if info.UpsertedId != nil {
        {{$instance}}.{{$IDField.Name}} = info.UpsertedId.(bson.ObjectId)
    }
    return err
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
