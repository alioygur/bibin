package mongo

import (
	"reflect"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func Test_repository_MakeFriend(t *testing.T) {
	r, deferFnc, err := newTestRepo(true)
	if err != nil {
		t.Fatal(err)
	}
	defer deferFnc()

	var u1, u2 user
	u1.ID, u2.ID = 1, 2

	if err := r.c(usersTbl).Insert(u1, u2); err != nil {
		t.Error(err)
		return
	}

	if err := r.MakeFriend(1, 2); err != nil {
		t.Error(err)
		return
	}

	// try add duplicate records
	if err := r.MakeFriend(2, 1); err != nil {
		t.Error(err)
		return
	}

	if err := r.c(usersTbl).Find(bson.M{"id": 1}).One(&u1); err != nil {
		t.Error(err)
		return
	}
	if err := r.c(usersTbl).Find(bson.M{"id": 2}).One(&u2); err != nil {
		t.Error(err)
		return
	}

	// ensure 1, and 2 users are friend, and there isn't duplicate rows
	if !reflect.DeepEqual(u1.FriendList, []uint64{2}) || !reflect.DeepEqual(u2.FriendList, []uint64{1}) {
		t.Errorf("want u1 friendlist %v, got %v and want u2 friendlist %v, got %v", []uint64{2}, u1.FriendList, []uint64{1}, u2.FriendList)
	}
}

func Test_repository_AreFriends(t *testing.T) {
	r, deferFnc, err := newTestRepo(true)
	if err != nil {
		t.Fatal(err)
	}
	defer deferFnc()

	var u1, u2 user
	u1.ID = 1
	u1.FirstName = "Ali"
	u1.FriendList = []uint64{2, 3, 4}

	u2.ID = 2
	u2.FirstName = "Huseyin"
	u2.FriendList = []uint64{1, 3, 4}

	if err := r.c(usersTbl).Insert(u1, u2); err != nil {
		t.Error(err)
		return
	}

	var tests = []struct {
		name string
		args []uint64
		want bool
	}{
		{"", []uint64{1, 2}, true},
		{"", []uint64{2, 1}, true},
		{"", []uint64{1, 3}, false},
		{"", []uint64{1, 4}, false},
		{"", []uint64{2, 3}, false},
		{"", []uint64{2, 4}, false},
		{"", []uint64{3, 4}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yes, err := r.AreFriends(tt.args[0], tt.args[1])
			if err != nil {
				t.Error(err)
				return
			}
			if yes != tt.want {
				t.Errorf("AreFriends(%q) = got %v, want %v", tt.args, yes, tt.want)
			}
		})
	}
}
