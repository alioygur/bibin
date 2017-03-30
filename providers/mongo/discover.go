package mongo

import (
	"time"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// DiscoverUsers discover users for discoverer user.
// Rules:
// gender == given gender
// age between ageMin and ageMax
// Not me
// Not my friends
func (r *repository) DiscoverPeople(userID uint64, gender domain.Gender, ageMin int, ageMax int, limit int) ([]*domain.User, error) {
	now := time.Now()
	min := time.Date(now.Year()-ageMax, 1, 1, 0, 0, 0, 0, now.Location())
	max := time.Date(now.Year()-ageMin, 12, 31, 0, 0, 0, 0, now.Location())

	sess := r.sess.Copy()
	defer sess.Close()
	usersCol := r.c(usersTbl).With(sess)

	var me user
	if err := usersCol.Find(bson.M{"id": userID}).One(&me); err != nil {
		return nil, errors.WithStack(wrapErr(err))
	}

	// just hack, so i don't need extra where condition
	me.FriendList = append(me.FriendList, me.ID)

	var users []*domain.User

	q := bson.M{
		"id":       bson.M{"$nin": me.FriendList},
		"gender":   gender,
		"birthday": bson.M{"$gte": min, "$lte": max},
	}
	err := r.c(usersTbl).Find(q).All(&users)

	return users, errors.WithStack(err)
}
