// +build integration

package fb

import (
	"log"
	"strconv"
	"testing"

	"github.com/alioygur/fb-tinder-app/domain"
)

func Test_repository_PutTestUser(t *testing.T) {
	r := newRepository()

	type args struct {
		installed bool
		perms     []string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FBTestUser
		wantErr bool
	}{
		{"", args{true, perms()}, nil, false},
		{"", args{false, perms()}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := r.PutTestUser(tt.args.installed, tt.args.perms)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID == "" || got.AccessToken == "" || got.LoginURL == "" {
				t.Errorf("there are nil fields")
				return
			}

			// check perms
			perms, err := r.Permissions(got.AccessToken)
			if err != nil {
				t.Fatal(err)
			}
			for _, p := range tt.args.perms {
				_, ok := perms[p]
				if !ok {
					t.Errorf("want perms %v got %v", tt.args.perms, perms)
					return
				}
			}

			// everything is okey. let delete test user
			if err := r.DeleteTestUser(got.ID); err != nil {
				log.Fatalf("delete user(%s) failed: %v ", got.ID, err)
			}
		})
	}
}

func Test_repository_FindTestUsers(t *testing.T) {
	r := newRepository()
	var users []*domain.FBTestUser
	defer func() {
		// delete test users
		for _, u := range users {
			r.DeleteTestUser(u.ID)
		}
	}()

	count := 5
	for i := 0; i < count; i++ {
		u, err := r.PutTestUser(true, perms())
		if err != nil {
			t.Fatalf("put test user failed: %v", err)
		}
		users = append(users, u)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindTestUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.FindTestUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, u := range users {
				var exists bool
				for _, gu := range got {
					if u.ID == gu.ID {
						exists = true
					}
				}
				if !exists {
					t.Errorf("user(%s) not exists in result", u.ID)
					return
				}
			}
		})
	}
}

func Test_repository_MakeFriend(t *testing.T) {
	r := newRepository()

	u1, err1 := r.PutTestUser(true, perms())
	u2, err2 := r.PutTestUser(true, perms())

	if err1 != nil || err2 != nil {
		t.Fatalf("put test user failed: %v, %v", err1, err2)
	}
	defer func() {
		r.DeleteTestUser(u1.ID)
		r.DeleteTestUser(u2.ID)
	}()

	if err := r.MakeFriend(u1, u2); err != nil {
		t.Fatal(err)
		return
	}

	friends, err := r.Friends(u1.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(friends) != 1 && strconv.Itoa(int(friends[0])) != u2.ID {
		t.Errorf("user(%s) and user(%s) are not friends", u1.ID, u2.ID)
	}
}

func Test_repository_DeleteTestUser(t *testing.T) {
	r := newRepository()

	u, err := r.PutTestUser(true, perms())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		r.DeleteTestUser(u.ID)
	}()

	if err := r.DeleteTestUser(u.ID); err != nil {
		t.Fatal(err)
	}

	users, err := r.FindTestUsers()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range users {
		if u.ID == v.ID {
			t.Fail()
		}
	}
}
