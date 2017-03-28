package fb

import (
	"net/http"

	"time"

	"strconv"

	"fmt"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/facebookgo/fbapi"
	"github.com/pkg/errors"
)

type (
	fbUser struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Gender    string `json:"gender"`
		Birthday  string `json:"birthday"`
	}
)

// OneByAccessToken returns user by its access token
func (r *repository) OneByAccessToken(accessToken string) (*domain.User, error) {
	var fbu fbUser
	var params []fbapi.Param

	params = append(params, fbapi.ParamAccessToken(accessToken))
	params = append(params, fbapi.ParamFields("first_name", "last_name", "email", "gender", "birthday"))
	urlValues, err := fbapi.ParamValues(params...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.URL.RawQuery = urlValues.Encode()

	if _, err := r.client.Do(req, &fbu); err != nil {
		return nil, wrapErr(err)
	}

	return mapToUser(&fbu)
}

// Permissions returns user's perms that has been given as map
func (r *repository) Permissions(accessToken string) (map[string]string, error) {
	type perm struct {
		Permission string `json:"permission"`
		Status     string `json:"status"`
	}
	var perms struct {
		Data []perm `json:"data"`
	}

	urlStr := fmt.Sprintf("/me/permissions?access_token=%s", accessToken)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if _, err := r.client.Do(req, &perms); err != nil {
		return nil, wrapErr(err)
	}

	permissions := make(map[string]string)
	for _, p := range perms.Data {
		permissions[p.Permission] = p.Status
	}

	return permissions, nil
}

// Friends gets users' friends who installed this app
func (r *repository) Friends(accessToken string) ([]uint64, error) {
	urlStr := fmt.Sprintf("/me/friends?access_token=%s", accessToken)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// paging
	var ids []uint64
	paging := r.paginate(req)
	for paging.Next() {
		var users []fbUser
		if err := paging.Data(&users); err != nil {
			return nil, err
		}

		for _, f := range users {
			id, err := strconv.Atoi(f.ID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			ids = append(ids, uint64(id))
		}
	}

	return ids, nil
}

// ProfilePicture returns user's profile picture url
func (r *repository) ProfilePicture(accessToken string) (string, error) {
	return fmt.Sprintf("https://graph.facebook.com/me/picture?type=large&access_token=%s", accessToken), nil
}

// mapToUser map user from fbUser to domain user
func mapToUser(fbu *fbUser) (*domain.User, error) {
	var u domain.User

	id, err := strconv.Atoi(fbu.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	u.FacebookID = uint64(id)
	u.Email = fbu.Email
	u.FirstName = fbu.FirstName
	u.LastName = fbu.LastName

	if fbu.Gender == "male" {
		u.Gender = domain.GenderMale
	} else if fbu.Gender == "female" {
		u.Gender = domain.GenderFemale
	} else {
		u.Gender = domain.GenderUnknown
	}

	if fbu.Birthday != "" {
		t, err := time.Parse("01/02/2006", fbu.Birthday)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		u.Birthday = &t
	}

	return &u, nil
}
