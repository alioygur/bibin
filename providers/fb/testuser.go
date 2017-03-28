package fb

import (
	"fmt"
	"net/http"

	"strings"

	"net/url"

	"strconv"

	"sync"

	"log"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/facebookgo/fbapi"
	"github.com/pkg/errors"
)

// SetFriendships sets friendships between test users.
// Formula: (FacebookID[:3] % 2) == 1
func SetFriendships(r service.FacebookRepository, users []*domain.FBTestUser) {
	friends, err := guesFriends(users)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for i, u := range friends {
		if i+1 == len(friends) {
			break
		}
		for _, f := range friends[i+1:] {
			wg.Add(1)
			go func(u, f *domain.FBTestUser) {
				defer wg.Done()
				if err := r.MakeFriend(u, f); err != nil {
					log.Println(err)
				}
			}(u, f)
		}
	}
	wg.Wait()
}

func guesFriends(users []*domain.FBTestUser) ([]*domain.FBTestUser, error) {
	var friends []*domain.FBTestUser
	for _, u := range users {
		id, err := strconv.Atoi(u.ID[:3])
		if err != nil {
			return nil, err
		}

		if (id % 2) == 1 {
			friends = append(friends, u)
		}
	}
	return friends, nil
}

func (r *repository) PutTestUser(installed bool, perms []string) (*domain.FBTestUser, error) {
	var fbtu domain.FBTestUser
	urlStr := fmt.Sprintf("/%d/accounts/test-users?access_token=%s", r.appID, r.appAccessToken())

	v := make(url.Values)
	v["permissions"] = []string{strings.Join(perms, ",")}
	if installed {
		v["installed"] = []string{"true"}
	}

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	_, err = r.client.Do(req, &fbtu)

	return &fbtu, errors.WithStack(err)
}

func (r *repository) FindTestUsers() ([]*domain.FBTestUser, error) {
	urlStr := fmt.Sprintf("/%d/accounts/test-users?access_token=%s", r.appID, r.appAccessToken())
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// paging
	var users []*domain.FBTestUser

	paging := r.paginate(req)

	for paging.Next() {
		var _users []*domain.FBTestUser
		if err := paging.Data(&_users); err != nil {
			return nil, err
		}
		users = append(users, _users...)
	}

	return users, nil
}

func (r *repository) MakeFriend(u1 *domain.FBTestUser, u2 *domain.FBTestUser) error {
	urlStr := fmt.Sprintf("/%s/friends/%s?access_token=%s", u1.ID, u2.ID, u1.AccessToken)

	// send request from user 1
	req, err := http.NewRequest(http.MethodPost, urlStr, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = r.client.Do(req, nil)
	if err != nil {
		// if err is already pending request from user 2, then accept the request from user 2
		e, ok := err.(*fbapi.Error)
		if ok && e.Code == 520 {
			return r.MakeFriend(u2, u1)
		}
		// already friend?
		if ok && e.Code == 522 {
			return nil
		}
		return err
	}

	// accept request from user 1
	return r.MakeFriend(u2, u1)
}

func (r *repository) DeleteTestUser(id string) error {
	urlStr := fmt.Sprintf("/%s?access_token=%s", id, r.appAccessToken())

	req, err := http.NewRequest(http.MethodDelete, urlStr, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = r.client.Do(req, nil)
	return errors.WithStack(err)
}
