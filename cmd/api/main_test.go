// +build ignore

// main integration tests
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"database/sql"

	"github.com/alioygur/fb-tinder-app/api"
	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/providers"
	fbrepo "github.com/alioygur/fb-tinder-app/providers/fb"
	mysqlrepo "github.com/alioygur/fb-tinder-app/providers/mysql"
	services "github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/fb-tinder-app/util"
)

type (
	afterFunc  func(*httptest.ResponseRecorder) error
	beforeFunc func(*http.Request) error

	testCase struct {
		name               string
		url                string
		method             string
		body               interface{}
		expectedStatusCode int
		afterFuncs         []afterFunc
	}

	testSuite struct {
		server http.Handler
		db     *sql.DB
		dbURL  string
	}
)

const (

	// api endpoints
	loginURL          = "/v1/auth/login"
	registerURL       = "/v1/auth/register"
	activateURL       = "/v1/auth/activate"
	forgotPasswordURL = "/v1/password/forgot"
	resetPasswordURL  = "/v1/password/reset"

	// http request methods
	mGet  = "GET"
	mPost = "POST"
	mPut  = "PUT"
	mDel  = "DELETE"

	// global user email and password
	email    = "user@example.com"
	password = "password"
)

var (
	ts          *testSuite
	authToken   string
	testUser    *domain.User
	fbTestUsers []*domain.FBTestUser
)

func TestMain(m *testing.M) {
	ts = newTestSuite()
	ts.setup()
	c := m.Run()
	ts.teardown()
	os.Exit(c)
}

func TestUser(t *testing.T) {
	type (
		rr services.RegisterRequest
	)

	// testCases := []testCase{
	// 	{"register with no params", registerURL, mPost, nil, 500, nil},
	// 	{"register with good params", registerURL, mPost, rr{fbTestUsers[0].AccessToken}, http.StatusOK, nil},
	// 	{"register with existsing user", registerURL, mPost, rr{fbTestUsers[0].AccessToken}, http.StatusOK, nil},
	// 	{"register with missing permissions", registerURL, mPost, rr{fbTestUsers[9].AccessToken}, 500, nil},
	// }

	var testCases []testCase
	for _, fbtu := range fbTestUsers {
		testCases = append(testCases, testCase{
			"register with good params",
			registerURL,
			mPost,
			rr{fbtu.AccessToken},
			http.StatusOK,
			nil,
		})
	}

	ts.runTestCases(testCases, t)
}

func newTestSuite() *testSuite {
	var ts testSuite

	if err := ensureEnv(); err != nil {
		log.Fatal(err)
	}

	db, err := mysqlrepo.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	ts.db = db

	// repos
	fb := fbrepo.New(uint64(util.EnvMustInt("FB_APP_ID")), util.EnvMustStr("FB_APP_SECRET"))
	sql := mysqlrepo.New(ts.db)

	// deps
	jwt := providers.NewJWT()

	// set facebook test users
	users, err := fb.FindTestUsers()
	if err != nil {
		log.Fatal(err)
	}
	fbTestUsers = users

	// services
	userService := services.New(fb, sql, jwt)

	ts.server = api.NewHandler(userService)

	return &ts
}

func (ts *testSuite) closeDB() {
	if err := ts.db.Close(); err != nil {
		log.Fatalf("close db connection failed: %v", err)
	}
}

func (ts *testSuite) setup() {
	ts.applyScheme()
}

func (ts *testSuite) applyScheme() {
	scheme := "../../providers/mysql/scheme.sql"
	if err := mysqlrepo.ApplyScheme(ts.db, scheme); err != nil {
		log.Fatalf("apply scheme: %v", err)
	}
}

func (ts *testSuite) teardown() {
	ts.closeDB()
}

func (ts *testSuite) runTestCases(tcs []testCase, t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ts.runTestCase(&tc, t)
		})
	}
}

func (ts *testSuite) runTestCase(tc *testCase, t *testing.T) {
	b, err := json.Marshal(tc.body)
	if err != nil {
		t.Fatalf("test case's body Marshal failed: %v", err)
	}
	reqBody := bytes.NewReader(b)
	r, err := http.NewRequest(tc.method, tc.url, reqBody)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}

	w := httptest.NewRecorder()

	ts.server.ServeHTTP(w, r)

	if w.Code != tc.expectedStatusCode {
		t.Errorf("%s %s (%s) status code want %v got %v", tc.method, tc.url, tc.name, tc.expectedStatusCode, w.Code)
		t.Logf("request body: %v", string(b))
		t.Logf("response body: %v", w.Body)
	}

	// run after funcs
	for _, cb := range tc.afterFuncs {
		if err := cb(w); err != nil {
			t.Error(err)
		}
	}
}
