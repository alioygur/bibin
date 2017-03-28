package mysql

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/goutil"
	"github.com/pkg/errors"
)

func (r *repository) PutPicture(p *domain.Image) error {
	p.CreatedAt = time.Now()
	id, err := r.insert(imagesTbl, p)
	p.ID = id
	return err
}

func (r *repository) ProfilePicture(userID uint64) (*domain.Image, error) {
	var p domain.Image
	return &p, r.oneBy(&p, imagesTbl, "user_id=? AND is_profile=?", userID, true)
}

func (r *repository) PicturesByUserID(userID uint64) ([]*domain.Image, error) {
	ss, err := goutil.NewSQLStruct(&domain.Image{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rows, err := squirrel.Select(ss.Columns()...).From(imagesTbl).Where("user_id=?", userID).RunWith(r.sess()).Query()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var pictures []*domain.Image
	for rows.Next() {
		var p domain.Image
		ss, err := goutil.NewSQLStruct(&p)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := rows.Scan(ss.Ptrs()...); err != nil {
			return nil, errors.WithStack(err)
		}
		pictures = append(pictures, &p)
	}
	return pictures, nil
}

func (r *repository) PictureByID(id uint64) (*domain.Image, error) {
	var p domain.Image
	return &p, r.oneBy(&p, imagesTbl, "id=?", id)
}

func (r *repository) PictureExistsByUserIDAndIsProfile(id uint64, isProfile bool) (bool, error) {
	return r.existsBy(imagesTbl, "user_id=? AND is_profile=?", id, isProfile)
}

func (r *repository) UpdatePicture(p *domain.Image) error {
	return r.update(p, imagesTbl, "id=?", p.ID)
}

func (r *repository) DeletePicture(id uint64) error {
	_, err := squirrel.Delete(imagesTbl).Where("id=?", id).RunWith(r.sess()).Exec()
	return errors.WithStack(err)
}
