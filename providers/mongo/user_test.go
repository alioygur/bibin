package mongo

import (
	"testing"

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
