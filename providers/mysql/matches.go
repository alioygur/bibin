package mysql

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// PutMatch insters to mathes table with cross records.
// insert userID, friendID
// insert friendID, userID
func (r *repository) PutMatch(userID, friendID uint64) error {
	now := time.Now()
	q := fmt.Sprintf(`insert into %s values(?, ?, ?), (?, ?, ?)`, matchesTbl)
	_, err := r.sess().Exec(q, userID, friendID, now, friendID, userID, now)
	return errors.WithStack(err)
}

func (r *repository) MatchExistsBy(userID, friendID uint64) (bool, error) {
	var matched bool
	q := fmt.Sprintf(`SELECT COUNT(f1.user_id) > 0 FROM %s AS f1 
	INNER JOIN %s AS f2 ON f1.user_id = f2.friend_id AND f1.friend_id = f2.user_id
	WHERE f1.user_id=? AND f1.friend_id=?`, matchesTbl, matchesTbl)

	err := r.sess().QueryRow(q, userID, friendID).Scan(&matched)

	return matched, errors.WithStack(err)
}
