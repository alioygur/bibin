package api

import (
	"net/http"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/gores"
	"github.com/pkg/errors"
)

type (
	handler struct {
		service.Service
	}
)

func (h *handler) register(w http.ResponseWriter, r *http.Request) error {
	req := new(service.RegisterRequest)
	if err := decodeReq(r, req); err != nil {
		return err
	}

	u, err := h.Register(req)
	if err != nil {
		return err
	}

	jwt, err := h.GenToken(u, service.AuthToken)
	if err != nil {
		return err
	}

	return gores.JSON(w, http.StatusOK, response{jwt})
}

func (h *handler) me(w http.ResponseWriter, r *http.Request) error {
	me := domain.UserMustFromContext(r.Context())
	req := service.ShowUserRequest{ID: me.ID}
	usr, err := h.Show(&req)
	if err != nil {
		return err
	}
	return gores.JSON(w, http.StatusOK, response{usr})
}

func (h *handler) react(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())

	req := new(service.ReactRequest)
	if err := decodeReq(r, req); err != nil {
		return err
	}
	req.FromUserID = u.ID

	if _, err := h.React(req); err != nil {
		return err
	}

	gores.NoContent(w)
	return nil
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())

	req := new(service.UpdateUserRequest)
	if err := decodeReq(r, req); err != nil {
		return err
	}
	req.ID = u.ID

	if err := h.UpdateUser(req); err != nil {
		return err
	}

	gores.NoContent(w)
	return nil
}

func (h *handler) discoverPeople(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())

	var rr service.DiscoverPeopleRequest
	rr.UserID = u.ID
	rr.Gender = domain.GenderFemale
	rr.AgeMin = 18
	rr.AgeMax = 100
	rr.Limit = 10

	g := queryValue("gender", r)
	if g != "m" && g != "f" {
		err := service.NewErr(service.ValidationErrCode, errors.New("gender must be m or f"))
		return errors.WithStack(err)
	}
	if g == "m" {
		rr.Gender = domain.GenderMale
	}

	aMin, err := queryValueInt("ageMin", r)
	if err != nil {
		err := service.NewErr(service.ValidationErrCode, err)
		return errors.WithStack(err)
	}
	rr.AgeMin = aMin

	aMax, err := queryValueInt("ageMax", r)
	if err != nil {
		err := service.NewErr(service.ValidationErrCode, err)
		return errors.WithStack(err)
	}
	rr.AgeMax = aMax

	people, err := h.DiscoverPeople(&rr)
	if err != nil {
		return err
	}
	return gores.JSON(w, http.StatusOK, response{people})
}
