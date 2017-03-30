package mongo

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"
)

// MakeFriend makes given two users friend
func (r *repository) MakeFriend(user1, user2 uint64) error {
	if err := r.c(usersTbl).Update(bson.M{"id": user1}, bson.M{"$addToSet": bson.M{"friend_list": user2}}); err != nil {
		return errors.WithStack(err)
	}
	err := r.c(usersTbl).Update(bson.M{"id": user2}, bson.M{"$addToSet": bson.M{"friend_list": user1}})
	return errors.WithStack(err)
}

// AreFriends checks given users are friends or not
func (r *repository) AreFriends(user1, user2 uint64) (bool, error) {
	user, err := r.existsBy(r.c(usersTbl), bson.M{"id": user1, "friend_list": bson.M{"$in": []uint64{user2}}})
	if err != nil {
		return false, err
	}
	friend, err := r.existsBy(r.c(usersTbl), bson.M{"id": user2, "friend_list": bson.M{"$in": []uint64{user1}}})

	return user && friend, err
}
