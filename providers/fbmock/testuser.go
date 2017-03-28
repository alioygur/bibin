package fbmock

import (
	"github.com/alioygur/fb-tinder-app/domain"
)

func (r *Repository) PutTestUser(installed bool, perms []string) (*domain.FBTestUser, error) {
	if r.PutTestUserFunc != nil {
		return r.PutTestUserFunc(installed, perms)

	}
	return nil, nil
}

func (r *Repository) FindTestUsers() ([]*domain.FBTestUser, error) {
	if r.FindTestUsersFunc != nil {
		return r.FindTestUsersFunc()
	}
	return nil, nil
}

func (r *Repository) DeleteTestUser(id string) error {
	if r.DeleteTestUserFunc != nil {
		r.DeleteTestUser(id)
	}
	return nil
}
