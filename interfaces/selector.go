package interfaces

import "gopkg.in/mgo.v2/bson"

type selector interface {
	GetBsonM() bson.M
}
