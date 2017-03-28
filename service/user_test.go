// +build ignore

package service_test

import (
	"reflect"
	"testing"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/providers"
	"github.com/alioygur/fb-tinder-app/providers/fbmock"
	"github.com/alioygur/fb-tinder-app/providers/mysql"
	. "github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/goutil"
)

func Test_user_Register(t *testing.T) {
	dbSess, err := mysql.ConnectToDB()
	if err != nil {
		t.Fatal(err)
	}
	mysql.TruncateTables(dbSess)

	sqlrepo := mysql.New(dbSess)
	fbrepo := fbmock.New(uint64(goutil.EnvMustInt("FB_APP_ID")), goutil.EnvMustGet("FB_APP_SECRET"))
	jwt := providers.NewJWT()
	user := New(fbrepo, sqlrepo, jwt)

	users := mysql.GenUsers(3)

	fbrepo.Users = users

	type afterFunc func(*testing.T, *domain.User, error)

	assertCredit := func(c int) afterFunc {
		return afterFunc(func(t *testing.T, u *domain.User, _ error) {
			// check credits
			if uc, _ := sqlrepo.CalcUserCredits(u.ID); uc != c {
				t.Errorf("user credit %d, want %d", uc, c)
			}
		})
	}

	assertErrCode := func(c ErrCode) afterFunc {
		return afterFunc(func(t *testing.T, u *domain.User, err error) {
			if !errCodeIs(err, c) {
				t.Errorf("error = %v, expected err code %d", err, c)
			}
		})
	}

	assertDeepEqual := func(v interface{}) afterFunc {
		return afterFunc(func(t *testing.T, u *domain.User, err error) {
			if !reflect.DeepEqual(u, v) {
				t.Errorf("= %v, want %v", u, v)
				return
			}
		})
	}

	tests := []struct {
		name        string
		accessToken string
		afterFunc   []afterFunc
	}{
		{"with bad access token", "", []afterFunc{assertErrCode(FacebookAPIErrCode)}},
		{"with missing perms", "99", []afterFunc{assertErrCode(PermissionNotGrantedErrCode)}},
		{"with good params", "1", []afterFunc{assertDeepEqual(users[0]), assertCredit(100)}},
		{"with good params", "2", []afterFunc{assertDeepEqual(users[1]), assertCredit(100)}},
		{"with good params", "3", []afterFunc{assertDeepEqual(users[2]), assertCredit(100)}},
		{"already exists user", "1", []afterFunc{assertDeepEqual(users[0]), assertCredit(100)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := RegisterRequest{AccessToken: tt.accessToken}
			got, err := user.Register(&rr)
			for _, f := range tt.afterFunc {
				f(t, got, err)
				if t.Failed() {
					break
				}
			}
		})
	}
}

func Test_user_React(t *testing.T) {
	mysql.TruncateTables(dbSess)

	sqlrepo := mysql.New(dbSess)
	fbrepo := fbmock.New(uint64(goutil.EnvMustInt("FB_APP_ID")), goutil.EnvMustGet("FB_APP_SECRET"))
	jwt := providers.NewJWT()
	user := New(fbrepo, sqlrepo, jwt)

	users := mysql.GenUsers(3)
	for _, u := range users {
		if err := sqlrepo.PutUser(u); err != nil {
			t.Fatal(err)
		}
	}

	// insert credit for user 1 and 2
	var c domain.Credit
	c.UserID = 1
	c.Amount = 100
	c.Type = domain.CreditGift
	if err := sqlrepo.PutCredit(&c); err != nil {
		t.Fatal(err)
	}
	c.ID = 0
	c.UserID = 2
	if err := sqlrepo.PutCredit(&c); err != nil {
		t.Fatal(err)
	}

	type afterFunc func(*testing.T, bool, error)

	assertErrNil := func(t *testing.T, _ bool, err error) {
		if err != nil {
			t.Errorf("error = %v, expected %v", err, nil)
		}
	}

	assertErrCode := func(c ErrCode) afterFunc {
		return afterFunc(func(t *testing.T, _ bool, err error) {
			if !errCodeIs(err, c) {
				t.Errorf("error = %v, expected err code %d", err, c)
			}
		})
	}

	assertCredit := func(userID uint64, c int) afterFunc {
		return afterFunc(func(t *testing.T, _ bool, err error) {
			// check credits
			if uc, _ := sqlrepo.CalcUserCredits(userID); uc != c {
				t.Errorf("user credit %d, want %d", uc, c)
			}
		})
	}

	assertMatches := func(expect bool) afterFunc {
		return afterFunc(func(t *testing.T, matched bool, _ error) {
			if expect != matched {
				t.Errorf("matched = %v, expected %v", matched, expect)
			}
		})
	}

	type args struct {
		r *ReactRequest
	}
	tests := []struct {
		name       string
		args       args
		afterFuncs []afterFunc
	}{
		{
			"from invalid user to valid user",
			args{&ReactRequest{FromUserID: 99, ToUserID: 1, Type: domain.ReactLike}},
			[]afterFunc{assertErrCode(NotFoundErrCode)},
		},
		{
			"from valid user to invalid user",
			args{&ReactRequest{FromUserID: 1, ToUserID: 99, Type: domain.ReactLike}},
			[]afterFunc{assertErrCode(NotFoundErrCode)},
		},
		{
			"from invalid user to invalid user",
			args{&ReactRequest{FromUserID: 98, ToUserID: 99, Type: domain.ReactLike}},
			[]afterFunc{assertErrCode(NotFoundErrCode)},
		},
		{
			"same users",
			args{&ReactRequest{FromUserID: 1, ToUserID: 1, Type: domain.ReactLike}},
			[]afterFunc{assertErrCode(UnknownErrCode)},
		},
		{
			"no more credits",
			args{&ReactRequest{FromUserID: 3, ToUserID: 1, Type: domain.ReactLike}},
			[]afterFunc{assertErrCode(NoMoreCreditErrCode)},
		},
		{
			"credit should be 99 and matches should be false",
			args{&ReactRequest{FromUserID: 1, ToUserID: 2, Type: domain.ReactLike}},
			[]afterFunc{assertErrNil, assertCredit(1, 99), assertMatches(false)},
		},
		{
			"duplicate entry",
			args{&ReactRequest{FromUserID: 1, ToUserID: 2, Type: domain.ReactLike}},
			[]afterFunc{assertErrCode(AlreadyExistsErrCode)},
		},
		{
			"credit should be 99 and matches should be true",
			args{&ReactRequest{FromUserID: 2, ToUserID: 1, Type: domain.ReactLike}},
			[]afterFunc{assertErrNil, assertCredit(2, 99), assertMatches(true)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, err := user.React(tt.args.r)
			for _, f := range tt.afterFuncs {
				f(t, matched, err)
				if t.Failed() {
					break
				}
			}
		})
	}
}

func errCodeIs(err error, code ErrCode) bool {
	if e, ok := err.(*Error); ok && e.Code == code {
		return true
	}
	return false
}
