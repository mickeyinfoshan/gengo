package interfaces

import "gopkg.in/mgo.v2/bson"
import "gopkg.in/mgo.v2"

type selector interface {
	MakeBsonM() bson.M
	MakeQuery(*mgo.Session) *mgo.Query
}
