package mysql

import (
	"time"

	"github.com/alioygur/fb-tinder-app/domain"
)

func (r *repository) AddAbuse(a *domain.Abuse) error {
	a.CreatedAt = time.Now()
	id, err := r.insert(abusesTbl, a)
	a.ID = id
	return err
}
