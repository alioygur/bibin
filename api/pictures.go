package api

import (
	"net/http"

	"strconv"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/gores"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (h *handler) uploadPicture(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())

	rr := new(service.UploadPictureRequest)
	rr.Picture = r.Body
	rr.UserID = u.ID

	img, err := h.UploadPicture(rr)
	if err != nil {
		return err
	}
	return gores.JSON(w, http.StatusCreated, response{img})
}

func (h *handler) setProfilePicture(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())
	pID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		err = service.NewErr(service.UnknownErrCode, err)
		return errors.WithStack(err)
	}
	if err := h.SetProfilePicture(u.ID, uint64(pID)); err != nil {
		return err
	}

	gores.NoContent(w)
	return nil
}

func (h *handler) pictures(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())

	var rr service.PicturesRequest
	rr.UserID = u.ID

	pics, err := h.Pictures(&rr)
	if err != nil {
		return err
	}

	return gores.JSON(w, http.StatusOK, response{pics})
}

func (h *handler) deletePicture(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())

	pID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		err = service.NewErr(service.UnknownErrCode, err)
		return errors.WithStack(err)
	}

	rr := &service.DeletePictureRequest{ID: uint64(pID), UserID: u.ID}
	if err := h.DeletePicture(rr); err != nil {
		return err
	}

	gores.NoContent(w)
	return nil
}
