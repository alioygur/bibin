package mongo

import (
	"reflect"
	"testing"

	"github.com/alioygur/fb-tinder-app/domain"
)

func Test_repository_DiscoverPeople(t *testing.T) {
	r, deferFnc, err := newTestRepo(true)
	if err != nil {
		t.Fatal(err)
	}
	defer deferFnc()

	// create users
	var count = 11
	var users []interface{}

	for i := 1; i < count; i++ {
		var u user
		u.ID = uint64(i)
		u.Gender = domain.GenderFemale
		u.SetAge(20 + i)

		if (i % 2) == 1 { // odd numbers are male
			u.Gender = domain.GenderMale
		}

		if i == 1 {
			u.FriendList = []uint64{5, 10}
		}
		users = append(users, &u)
	}

	if err := r.c(usersTbl).Insert(users...); err != nil {
		t.Error(err)
		return
	}

	type args struct {
		userID uint64
		gender domain.Gender
		ageMin int
		ageMax int
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{"for:1, females between 18 and 40 years old", args{1, domain.GenderFemale, 18, 40, 10}, []uint64{2, 4, 6, 8}, false},
		{"for:1, females between 18 and 25 years old", args{1, domain.GenderFemale, 18, 25, 10}, []uint64{2, 4}, false},
		{"for:1, males between 18 and 40 years old", args{1, domain.GenderMale, 18, 40, 10}, []uint64{3, 7, 9}, false},
		{"for:1, males older than max age", args{1, domain.GenderMale, 31, 40, 10}, nil, false},
		{"for:2, females between 18 and 40 years old", args{2, domain.GenderFemale, 18, 40, 10}, []uint64{4, 6, 8, 10}, false},
		{"for:2, females between 28 and 30 years old", args{2, domain.GenderFemale, 28, 30, 10}, []uint64{8, 10}, false},
		{"for:2, females between 30 and 40 years old", args{2, domain.GenderFemale, 30, 40, 10}, []uint64{10}, false},
		{"for:2, males between 18 and 40 years old", args{2, domain.GenderMale, 18, 40, 10}, []uint64{1, 3, 5, 7, 9}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.DiscoverPeople(tt.args.userID, tt.args.gender, tt.args.ageMin, tt.args.ageMax, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.DiscoverPeople() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotIDs := userIDs(got)
			if !reflect.DeepEqual(gotIDs, tt.want) {
				t.Errorf("repository.DiscoverPeople() = %v, want %v", gotIDs, tt.want)
			}
		})
	}
}

func userIDs(users []*domain.User) []uint64 {
	var ids []uint64
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	return ids
}
