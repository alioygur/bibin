package mysql

import (
	"time"

	"fmt"

	"database/sql"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/pkg/errors"
)

func (r *repository) PutCredit(c *domain.Credit) error {
	c.CreatedAt = time.Now()
	id, err := r.insert(creditsTbl, c)
	c.ID = id
	return err
}

func (r *repository) CalcUserCredits(id uint64) (int, error) {
	var total sql.NullInt64
	row := r.db.QueryRow(fmt.Sprintf("SELECT SUM(amount) FROM %s WHERE user_id=?", creditsTbl), id)

	err := row.Scan(&total)

	if total.Valid {
		return int(total.Int64), nil
	}
	return 0, errors.WithStack(err)
}
