package api

import (
	"net/http"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/gores"
)

func (h *handler) reportAbuse(w http.ResponseWriter, r *http.Request) error {
	u := domain.UserMustFromContext(r.Context())

	req := new(service.ReportAbuseRequest)
	if err := decodeReq(r, req); err != nil {
		return err
	}
	req.UserID = u.ID

	// todo: handle errors like duplicate
	_, err := h.ReportAbuse(req)
	if err != nil {
		return err
	}

	gores.NoContent(w)
	return nil
}
