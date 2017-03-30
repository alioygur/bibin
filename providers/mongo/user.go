package mongo

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/pkg/errors"
)

type (
	user struct {
		domain.User `bson:",inline"`
		FriendList  []uint64 `bson:"friend_list"`
	}
)

func (r *repository) AddUser(du *domain.User) error {
	now := time.Now().Round(time.Second)
	du.CreatedAt = now
	du.UpdatedAt = now
	du.ID = r.id(usersTbl)

	var u user
	u.User = *du

	return r.c(usersTbl).Insert(&u)
}

func (r *repository) UserByID(id uint64) (*domain.User, error) {
	var u domain.User
	err := r.c(usersTbl).Find(bson.M{"id": id}).One(&u)

	return &u, errors.WithStack(wrapErr(err))
}

func (r *repository) UserByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.c(usersTbl).Find(bson.M{"email": email}).One(&u)

	return &u, errors.WithStack(wrapErr(err))
}

func (r *repository) UserByFacebookID(id uint64) (*domain.User, error) {
	var u domain.User
	err := r.c(usersTbl).Find(bson.M{"facebook_id": id}).One(&u)

	return &u, errors.WithStack(wrapErr(err))
}

func (r *repository) ExistsByEmail(email string) (bool, error) {
	n, err := r.c(usersTbl).Find(bson.M{"email": email}).Count()

	return n > 0, errors.WithStack(err)
}

func (r *repository) UserExistsByID(id uint64) (bool, error) {
	n, err := r.c(usersTbl).Find(bson.M{"id": id}).Count()

	return n > 0, errors.WithStack(err)
}

func (r *repository) UserExistsByFacebookID(id uint64) (bool, error) {
	n, err := r.c(usersTbl).Find(bson.M{"facebook_id": id}).Count()

	return n > 0, errors.WithStack(err)
}

func (r *repository) UpdateUser(u *domain.User) error {
	u.UpdatedAt = time.Now()
	err := r.c(usersTbl).Update(bson.M{"id": u.ID}, u)

	return errors.WithStack(wrapErr(err))
}

func (r *repository) SyncUserFriendsByFacebookID(user uint64, friends []uint64) error {
	return nil
}

func (r *repository) BindFriends(users ...*domain.User) error {
	count := len(users)
	errc := make(chan error, count)
	for _, u := range users {
		go func(u *domain.User) {
			sess := r.sess.Copy()
			defer sess.Close()

			c := r.c(usersTbl).With(sess)
			var res struct {
				FriendList []uint64 `bson:"friend_list"`
			}
			if err := c.Find(bson.M{"id": u.ID}).Select(bson.M{"friend_list": 1}).One(&res); err != nil {
				errc <- errors.WithStack(err)
				return
			}

			err := c.Find(bson.M{"id": bson.M{"$in": res.FriendList}}).All(&u.Friends)
			errc <- errors.WithStack(err)
		}(u)
	}

	// check for errors
	for i := 0; i < count; i++ {
		if err := <-errc; err != nil {
			return err
		}
	}
	return nil
}
