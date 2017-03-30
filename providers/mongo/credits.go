package mongo

import (
	"time"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/pkg/errors"
)

func (r *repository) PutCredit(c *domain.Credit) error {
	c.CreatedAt = time.Now()
	c.ID = r.id(creditsTbl)
	err := r.c(creditsTbl).Insert(c)

	return errors.WithStack(wrapErr(err))
}

// todo: implement it
func (r *repository) CalcUserCredits(id uint64) (int, error) {
	return 0, nil
}
