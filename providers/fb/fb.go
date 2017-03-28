package fb

import (
	"fmt"
	"net/url"

	"github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/goutil"
	"github.com/facebookgo/fbapi"
	"github.com/pkg/errors"
)

type (
	repository struct {
		appID     uint64
		appSecret string
		client    *fbapi.Client
	}
)

const (
	apiVersion = "v2.8"
)

var requiredPerms []string

// New instances new facebook repository
func New(appID uint64, appSecret string) service.FacebookRepository {
	baseURL := &url.URL{
		Scheme: "https",
		Host:   "graph.facebook.com",
		Path:   "/" + apiVersion,
	}
	return &repository{
		appID:     appID,
		appSecret: appSecret,
		client:    &fbapi.Client{BaseURL: baseURL},
	}
}

func (r *repository) appAccessToken() string {
	return fmt.Sprintf("%d|%s", r.appID, r.appSecret)
}

func wrapErr(err error) error {
	if _, ok := err.(*fbapi.Error); ok {
		err = service.NewErr(service.FacebookAPIErrCode, err)
		return errors.WithStack(err)
	}
	return errors.WithStack(err)
}

// perms returns required facebook app permissions as slice of string
func perms() []string {
	if requiredPerms == nil {
		requiredPerms = goutil.EnvMustSliceStr("FB_REQUIRED_PERMS", ",")
	}
	return requiredPerms
}
