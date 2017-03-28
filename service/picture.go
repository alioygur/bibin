package service

import (
	"io"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/pkg/errors"
)

type (
	// PicturesRequest ....
	PicturesRequest struct {
		UserID uint64
	}

	// UploadPictureRequest ...
	UploadPictureRequest struct {
		UserID  uint64
		Picture io.Reader
	}

	// DeletePictureRequest ...
	DeletePictureRequest struct {
		ID     uint64
		UserID uint64 // condition
	}
)

// Pictures gets pictures
func (s *service) Pictures(r *PicturesRequest) ([]*domain.Image, error) {
	if r.UserID == 0 {
		return nil, errors.New("you must set UserID")
	}
	return s.storage.PicturesByUserID(r.UserID)
}

// UploadPicture uploads picture
func (s *service) UploadPicture(r *UploadPictureRequest) (*domain.Image, error) {
	img, err := s.imgCDN.Upload(r.Picture)
	if err != nil {
		return nil, err
	}

	img.UserID = r.UserID
	img.IsProfile = boolPtr(false)

	// if it first, then set as profile picture
	exists, err := s.storage.PictureExistsByUserIDAndIsProfile(r.UserID, true)
	if err != nil {
		return nil, err
	}

	if !exists {
		img.IsProfile = boolPtr(true)
	}

	return img, s.storage.PutPicture(img)
}

// SetProfilePicture sets user's profile picture
func (s *service) SetProfilePicture(uID, pID uint64) error {
	p, err := s.storage.PictureByID(pID)
	if err != nil {
		return err
	}

	if p.UserID != uID || *p.IsProfile == true {
		return nil
	}

	// unset old profile picture
	pp, err := s.storage.ProfilePicture(uID)
	if err != nil && !isNotFoundErr(err) {
		return err
	}

	pp.IsProfile = boolPtr(false)
	if err := s.storage.UpdatePicture(pp); err != nil {
		return err
	}

	// set new profile picture
	p.IsProfile = boolPtr(true)

	return s.storage.UpdatePicture(p)
}

// DeletePicture deletes picture
func (s *service) DeletePicture(r *DeletePictureRequest) error {
	if r.ID == 0 {
		return errors.New("DeletePictureRequest.ID required")
	}

	p, err := s.storage.PictureByID(r.ID)
	if err != nil {
		return err
	}

	if r.UserID == 0 || r.UserID == p.UserID {
		return s.storage.DeletePicture(r.ID)
	}

	return nil
}
