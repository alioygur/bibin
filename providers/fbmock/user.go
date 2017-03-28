package fbmock

import (
	"strconv"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"

	"github.com/alioygur/goutil"
	"github.com/pkg/errors"
)

type (
	// Repository ...
	Repository struct {
		Users []*domain.User

		OneByAccessTokenFunc func(string) (*domain.User, error)
		PermissionsFunc      func(string) (map[string]string, error)
		MakeFriendFunc       func(*domain.FBTestUser, *domain.FBTestUser) error
		FriendsFunc          func(string) ([]uint64, error)
		ProfilePictureFunc   func(string) (string, error)

		PutTestUserFunc    func(bool, []string) (*domain.FBTestUser, error)
		FindTestUsersFunc  func() ([]*domain.FBTestUser, error)
		DeleteTestUserFunc func(string) error
	}
)

// New instances new repository
func New(appID uint64, appSecret string) *Repository {
	return &Repository{}
}

// OneByAccessToken returns a user by accessToken
// if OneByAccessTokenFunc != nil returns it, else return default
// default
// if accessToken == "" returns error, else converts access token to int then
// returns a user at Users[index]
func (r *Repository) OneByAccessToken(accessToken string) (*domain.User, error) {
	if r.OneByAccessTokenFunc != nil {
		return r.OneByAccessTokenFunc(accessToken)
	}

	if accessToken == "" {
		err := service.NewErr(service.FacebookAPIErrCode, errors.New("access token not provided"))
		return nil, errors.WithStack(err)
	}
	id, err := strconv.Atoi(accessToken)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return r.Users[id-1], nil
}

// Permissions returns permissions of user
// if PermissionsFunc != nil returns it, else default
// if accessToken == "" returns error
// if accessToken == "99" returns missing permissions
func (r *Repository) Permissions(accessToken string) (map[string]string, error) {
	if r.PermissionsFunc != nil {
		return r.PermissionsFunc(accessToken)
	}

	envPerms := goutil.EnvMustSliceStr("FB_REQUIRED_PERMS", ",")

	if accessToken == "" {
		err := service.NewErr(service.FacebookAPIErrCode, errors.New("access token not provided"))
		return nil, errors.WithStack(err)
	}

	// return missing perms
	if accessToken == "99" {
		envPerms = envPerms[2:]
	}

	perms := make(map[string]string)
	for _, p := range envPerms {
		perms[p] = "granted"
	}
	return perms, nil
}

// MakeFriend does nothing
func (r *Repository) MakeFriend(u1 *domain.FBTestUser, u2 *domain.FBTestUser) error {
	if r.MakeFriendFunc != nil {
		return r.MakeFriendFunc(u1, u2)
	}
	return nil
}

// Friends does nothing
func (r *Repository) Friends(accessToken string) ([]uint64, error) {
	if r.FriendsFunc != nil {
		return r.FriendsFunc(accessToken)
	}
	return nil, nil
}

// ProfilePicture always returns same picture url
func (r *Repository) ProfilePicture(accessToken string) (string, error) {
	if r.ProfilePictureFunc != nil {
		return r.ProfilePictureFunc(accessToken)
	}
	if accessToken == "" {
		err := service.NewErr(service.FacebookAPIErrCode, errors.New("access token not provided"))
		return "", errors.WithStack(err)
	}
	return "https://x1.xingassets.com/assets/frontend_minified/img/users/nobody_m.original.jpg", nil
}
