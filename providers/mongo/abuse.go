package mongo

import (
	"time"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/pkg/errors"
)

func (r *repository) AddAbuse(a *domain.Abuse) error {
	a.CreatedAt = time.Now()
	a.ID = r.id(abusesTbl)

	err := r.c(abusesTbl).Insert(a)
	return errors.WithStack(wrapErr(err))
}
