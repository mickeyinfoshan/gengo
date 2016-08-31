package templates

const InitDBTempl = `
    package {{.PackageName}}

    import (
	    "gopkg.in/mgo.v2"
    )

    var (
        dbName string
        dbURL string
        collectionNames map[string]string
        seedSession *mgo.Session
    )

    func InitDB() error {
        dbName = "{{.DatabaseName}}"
        dbURL = "{{.DatabaseURL}}"
        {{range $index, $modelMeta := .ModelMetas}}
        collectionNames["{{$modelMeta.Name}}"] = "{{$modelMeta.Name}}"
        {{end}}
        var err error
        seedSession, err = mgo.Dial(dbURL)
        if err == nil {
            seedSession.SetMode(mgo.Monotonic, true)
        }
        return err
    }

    func NewDBSession() *mgo.Session {
        sess := seedSession.Clone()
        return sess
    }
`
