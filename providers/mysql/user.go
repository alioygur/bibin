package mysql

import (
	"time"

	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/goutil"
	"github.com/pkg/errors"
)

func (r *repository) PutUser(u *domain.User) error {
	now := time.Now().Round(time.Second)
	u.CreatedAt = now
	u.UpdatedAt = now

	id, err := r.insert(usersTbl, u)
	if err != nil {
		return err
	}
	u.ID = id

	// insert images if existsBy
	if len(u.Images) > 0 {
		for _, img := range u.Images {
			img.CreatedAt = now
			img.UserID = id
			_, err := r.insert(imagesTbl, img)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *repository) One(id uint64) (*domain.User, error) {
	var u domain.User
	return &u, r.oneBy(&u, usersTbl, "id=?", id)
}

func (r *repository) UserByEmail(email string) (*domain.User, error) {
	var u domain.User
	return &u, r.oneBy(&u, usersTbl, "email=?", email)
}

func (r *repository) UserByID(id uint64) (*domain.User, error) {
	var u domain.User
	return &u, r.oneBy(&u, usersTbl, "id=?", id)
}

func (r *repository) UserByFacebookID(id uint64) (*domain.User, error) {
	var u domain.User
	return &u, r.oneBy(&u, usersTbl, "facebook_id=?", id)
}

func (r *repository) ExistsByEmail(email string) (bool, error) {
	return r.existsBy(usersTbl, "email=?", email)
}

func (r *repository) UserExistsByID(id uint64) (bool, error) {
	return r.existsBy(usersTbl, "id=?", id)
}

func (r *repository) UserExistsByFacebookID(id uint64) (bool, error) {
	return r.existsBy(usersTbl, "facebook_id=?", id)
}

func (r *repository) UpdateUser(u *domain.User) error {
	u.UpdatedAt = time.Now()
	return r.update(u, usersTbl, "id=?", u.ID)
}

// TODO: handle errors. transactions are required?
func (r *repository) SyncUserFriendsByFacebookID(user uint64, friends []uint64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// delete old friendships
	_, err = tx.Exec(fmt.Sprintf(`delete from %s where user_id=?`, friendshipsTbl), user)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := squirrel.Select("id").
		From(usersTbl).
		Where(squirrel.Eq{"facebook_id": friends}).
		RunWith(r.db).
		Query()
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var friend uint64
		if err := rows.Scan(&friend); err != nil {
			tx.Rollback()
			return err
		}

		_, err := tx.Exec(fmt.Sprintf(`insert into %s values(?, ?)`, friendshipsTbl), user, friend)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *repository) DiscoverUsers(user uint64, gender domain.Gender, ageMin int, ageMax int, limit int) ([]*domain.User, error) {
	now := time.Now()
	min := now.Year() - ageMax
	max := now.Year() - ageMin

	rows, err := squirrel.Select("users.*").
		From(usersTbl).
		Where("id != ?", user).
		Where("gender = ?", gender).
		Where("(YEAR(birthday) BETWEEN ? AND ?)", min, max).
		// Where("id NOT IN (SELECT friend_id FROM friendships WHERE user_id = ?)", user).
		Where("id NOT IN (SELECT to_user_id FROM reactions WHERE from_user_id = ?)", user).
		Limit(uint64(limit)).
		RunWith(r.sess()).
		Query()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		ss, err := goutil.NewSQLStruct(&u)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if err := rows.Scan(ss.Ptrs()...); err != nil {
			return nil, errors.WithStack(err)
		}
		users = append(users, &u)
	}
	return users, nil
}
