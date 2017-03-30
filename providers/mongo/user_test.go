package mongo

import (
	"strconv"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"reflect"

	"github.com/alioygur/fb-tinder-app/domain"
)

func Test_repository_PutUser(t *testing.T) {
	r, deferFnc, err := newTestRepo(true)
	if err != nil {
		t.Fatal(err)
	}
	defer deferFnc()

	want := domain.NewUser()
	want.FirstName = "Ali"
	want.LastName = "OYGUR"
	want.Email = "alioygur@gmail.com"

	if err := r.AddUser(want); err != nil {
		t.Error(err)
		return
	}

	var got domain.User

	if err := r.c(usersTbl).Find(bson.M{"id": want.ID}).One(&got); err != nil {
		t.Error(err)
		return
	}

	got.Friends = make([]*domain.User, 0)

	if !reflect.DeepEqual(got, *want) {
		t.Errorf("not equal, got %+v, want %+v", got, want)
	}
}

// GenUsers generates number of users. useful for tests
func GenUsers(count int) []*domain.User {
	var users []*domain.User
	now := time.Now().Round(time.Second)
	for i := 1; i <= count; i++ {
		idStr := strconv.Itoa(i)
		var u domain.User
		u.FacebookID = uint64(i)
		u.FirstName = "Jhon " + idStr
		u.LastName = "Doe"
		u.Email = "user" + idStr + "@example.com"
		u.Gender = domain.GenderFemale
		if (i % 2) == 1 {
			u.Gender = domain.GenderMale
		}
		// age = 17 + i
		b := time.Date(now.Year()-(17+i), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		u.Birthday = &b
		u.Status = domain.StatusNewUser
		fls := false
		u.IsAdmin = &fls
		u.CreatedAt = now
		u.UpdatedAt = now

		users = append(users, &u)
	}
	return users
}

func FindFriends(userID uint64, totalUsers int) []uint64 {
	var friends []uint64
	var i uint64
	for i = 1; i <= uint64(totalUsers); i++ {
		if i == userID {
			continue
		}
		if (i % 2) == 1 {
			friends = append(friends, i)
		}
	}
	return friends
}
