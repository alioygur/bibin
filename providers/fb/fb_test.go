// +build integration

package fb

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/facebookgo/fbapi"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("error when loading .env file: %v", err)
	}
	os.Exit(m.Run())
}

// newRepository instances new repository.
// uses by tests
func newRepository() *repository {
	appID, err := strconv.Atoi(os.Getenv("FB_APP_ID"))
	if err != nil {
		panic(err)
	}
	appSecret := os.Getenv("FB_APP_SECRET")

	return New(uint64(appID), appSecret).(*repository)
}

func Test_repository_appAccessToken(t *testing.T) {
	type fields struct {
		appID     uint64
		appSecret string
		client    *fbapi.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{12345, "secret12345", nil}, "12345|secret12345"},
		{"", fields{123456789, "secret123456789", nil}, "123456789|secret123456789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				appID:     tt.fields.appID,
				appSecret: tt.fields.appSecret,
				client:    tt.fields.client,
			}
			if got := r.appAccessToken(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
