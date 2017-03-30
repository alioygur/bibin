package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	nextID struct {
		Next uint64 `bson:"n"`
	}
)

const idTbl = `ids`

// id returns next id.
// if sess nil then it uses default session
func (r *repository) id(c string) uint64 {
	ids := r.c(idTbl)
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"n": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	id := new(nextID)
	ids.Find(bson.M{"_id": c}).Apply(change, id)
	return id.Next
}
