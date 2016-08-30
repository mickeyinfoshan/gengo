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

// Delete Delete {{.Name}} from database
func ({{$instance}} *{{.Name}}) Delete() error {
    sess := NewDBSession()
    defer sess.Close()
    err := db(dbName).C(collectionNames["{{.Name}}"]).RemoveId({{$instance}}.{{IDField.Name}})
    return err
}`
