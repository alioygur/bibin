package mysql

import (
	"log"
	"os"
	"testing"

	"database/sql"

	"github.com/joho/godotenv"
)

var (
	dbSess *sql.DB
)

func TestMain(m *testing.M) {
	// set env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("error when loading .env file: %v", err)
	}

	// connect to mysql
	db, err := ConnectToDB()
	if err != nil {
		log.Fatalf("connect db err: %v", err)
	}
	dbSess = db

	// load db scheme
	if err := ApplyScheme(db, "./scheme.sql"); err != nil {
		log.Fatalf("reset db err: %v", err)
	}

	c := m.Run()

	dbSess.Close()

	os.Exit(c)
}

func newRepository() *repository {
	return New(dbSess).(*repository)
}

func Test_repository_insert(t *testing.T) {
	r := newRepository()

	TruncateTables(dbSess)
	users := GenUsers(5)

	type args struct {
		tbl string
		v   interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"with good params", args{usersTbl, users[0]}, 1, false},
		{"with good params", args{usersTbl, users[1]}, 2, false},
		{"with good params", args{usersTbl, users[2]}, 3, false},
		{"with good params", args{usersTbl, users[3]}, 4, false},
		{"with good params", args{usersTbl, users[4]}, 5, false},
		{"to unknown table", args{"unknown table", nil}, 0, true},
		{"with nil", args{usersTbl, nil}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.insert(tt.args.tbl, tt.args.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("repository.insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("repository.insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_oneBy(t *testing.T) {
	r := newRepository()

	TruncateTables(dbSess)

	users := GenUsers(3)

	for _, u := range users {
		id, err := r.insert(usersTbl, u)
		if err != nil {
			t.Fatalf("insert failed: %v", err)
		}
		u.ID = id
	}

	type args struct {
		v    interface{}
		tbl  string
		w    string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"one by id", args{users[0], usersTbl, "id=?", []interface{}{1}}, false},
		{"one by email", args{users[1], usersTbl, "email=?", []interface{}{"user2@example.com"}}, false},
		{"one by id and email", args{users[2], usersTbl, "id=? AND email=?", []interface{}{3, "user3@example.com"}}, false},
		{"one by unknown column", args{users[0], usersTbl, "unknown=?", []interface{}{1}}, true},
		{"one by empty", args{users[0], usersTbl, "", []interface{}{}}, true},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.oneBy(tt.args.v, tt.args.tbl, tt.args.w, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("repository.oneBy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if i < len(users) {
				if users[i].ID != uint64(i+1) {
					t.Errorf("repository.oneBy() row id: %d, want: %d", users[i].ID, uint64(i+1))
				}
			}
		})
	}
}

func Test_repository_existsBy(t *testing.T) {
	r := newRepository()

	TruncateTables(dbSess)

	users := GenUsers(3)

	for _, u := range users {
		id, err := r.insert(usersTbl, u)
		if err != nil {
			t.Fatalf("insert failed: %v", err)
		}
		u.ID = id
	}
	type args struct {
		tbl  string
		w    string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"exists by id", args{usersTbl, "id=?", []interface{}{1}}, true, false},
		{"exists by email", args{usersTbl, "email=?", []interface{}{"user2@example.com"}}, true, false},
		{"exists by id and email", args{usersTbl, "id=? AND email=?", []interface{}{3, "user3@example.com"}}, true, false},
		{"exists by unknown column", args{usersTbl, "unknown=?", []interface{}{1}}, false, true},
		{"exists by empty", args{usersTbl, "", []interface{}{}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.existsBy(tt.args.tbl, tt.args.w, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.existsBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repository.existsBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
