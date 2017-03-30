package mongo

import (
	"time"

	"strconv"

	mgo "gopkg.in/mgo.v2"
)

func randDBName() string {
	return "testdb_" + strconv.Itoa(time.Now().Nanosecond())
}

func newTestRepo(dropDB bool) (*repository, func(), error) {
	dbname := randDBName()
	sess, err := mgo.Dial("localhost/" + dbname)
	if err != nil {
		return nil, nil, nil
	}
	deferFnc := func() {
		if dropDB {
			sess.DB("").DropDatabase()
		}
		sess.Close()
	}

	r := repository{sess: sess}
	return &r, deferFnc, nil
}
