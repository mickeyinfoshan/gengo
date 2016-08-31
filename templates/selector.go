package templates

// SelectorTmpl 选择器的控制器
const SelectorTmpl = `
{{$modelMeta := .ModelMeta}}
// model {{$modelMeta.Name}}'s Selector
type {{$modelMeta.Name}}Selector struct {
    {{range $index, $field := .OmitEmptyFields}}{{$field.Name}} {{$field.Type}}{{end}}
    {{range $index, $field := .NoneOmitEmptyFields}}{{$field.Name}} {{$field.Type}}{{end}}
}

// MakeBsonM generate a bson.M object from Selector
func ({{$modelMeta.GetInstanceName}}Selector *{{$modelMeta.Name}}Selector) MakeBsonM() bson.M {
    bsonM := bson.M{}
    {{range $index, $field := .NoneOmitEmptyFields}}bsonM["{{$field.GetMongoAttrName}}"] = {{$modelMeta.GetInstanceName}}Selector.{{$field.Name}}{{end}}
    {{range $index, $field := .OmitEmptyFields}}
    var default{{$field.Name}} {{$field.Type}}
    if default{{$field.Name}} != {{$modelMeta.GetInstanceName}}Selector.{{$field.Name}} {
        bsonM["{{$field.GetMongoAttrName}}"] = {{$modelMeta.GetInstanceName}}Selector.{{$field.Name}}
    }
    {{end}}
    return bsonM
}

// MakeQuery generate a qurey from a Selector
func ({{$modelMeta.GetInstanceName}}Selector *{{$modelMeta.Name}}Selector) MakeQuery(sess *mgo.Session) *mgo.Query {
    bsonM := {{$modelMeta.GetInstanceName}}Selector.MakeBsonM()
    query := sess.DB(dbName).C(collectionNames["{{$modelMeta.Name}}"]).Find(bsonM)
    return query
}

// FindOne Get an instance from a {{$modelMeta.Name}}Selector
func ({{$modelMeta.GetInstanceName}}Selector *{{$modelMeta.Name}}Selector) FindOne() ({{$modelMeta.Name}}, error) {
    sess := NewDBSession()
    defer sess.Close()
    var {{$modelMeta.GetInstanceName}} {{$modelMeta.Name}}
    err := {{$modelMeta.GetInstanceName}}Selector.MakeQuery(sess).One(&{{$modelMeta.GetInstanceName}})
    return {{$modelMeta.GetInstanceName}}, err
}

// FindAll Get all instances from a {{$modelMeta.Name}}Selector
func ({{$modelMeta.GetInstanceName}}Selector *{{$modelMeta.Name}}Selector) FindAll() ([]{{$modelMeta.Name}}, error) {
    sess := NewDBSession()
    defer sess.Close()
    var {{$modelMeta.GetInstanceName}}s []{{$modelMeta.Name}}
    err := {{$modelMeta.GetInstanceName}}Selector.MakeQuery(sess).All(&{{$modelMeta.GetInstanceName}}s)
    return {{$modelMeta.GetInstanceName}}s, err
}
`
