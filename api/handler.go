package api

import (
	"net/http"

	"fmt"

	"os"

	"github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/gores"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
)

type (
	errHandler func(http.ResponseWriter, *http.Request) error
	httpErr    struct {
		Code     service.ErrCode `json:"code"`
		HTTPCode int             `json:"httpCode"`
		Msg      string          `json:"error"`
	}
	errResponse struct {
		Code     service.ErrCode `json:"code"`
		HTTPCode int             `json:"httpCode"`
		Error    string          `json:"error"`
		Inner    string          `json:"inner,omitempty"`
	}
)

func (err *httpErr) Error() string {
	return ""
}

// NewHandler set routes and returns the http handler
func NewHandler(s service.Service) http.Handler {
	r := mux.NewRouter()
	// base handler
	base := alice.New(newSetUserMid(s))
	// handler with auth required
	authRequired := base.Append(newAuthRequiredMid)

	h := &handler{s}

	// r.PathPrefix("/images").Handler(httputil.NewSingleHostReverseProxy(proxyURL))
	r.Handle("/v1/login", base.Then(errHandler(h.register))).Methods(http.MethodPost)
	r.Handle("/v1/me", authRequired.Then(errHandler(h.me))).Methods(http.MethodGet)
	r.Handle("/v1/me", authRequired.Then(errHandler(h.update))).Methods(http.MethodPatch)
	r.Handle("/v1/me/reacts", authRequired.Then(errHandler(h.react))).Methods(http.MethodPost)
	r.Handle("/v1/me/abuses", authRequired.Then(errHandler(h.reportAbuse))).Methods(http.MethodPost)

	r.Handle("/v1/me/discover-people", authRequired.Then(errHandler(h.discoverPeople))).Methods(http.MethodGet)

	r.Handle("/v1/me/pictures", authRequired.Then(errHandler(h.uploadPicture))).Methods(http.MethodPost)
	r.Handle("/v1/me/pictures", authRequired.Then(errHandler(h.pictures))).Methods(http.MethodGet)
	r.Handle("/v1/me/pictures/{id}", authRequired.Then(errHandler(h.deletePicture))).Methods(http.MethodDelete)
	r.Handle("/v1/me/pictures/{id}/profile", authRequired.Then(errHandler(h.setProfilePicture))).Methods(http.MethodPut)

	return r
}

func (h errHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errRes errResponse
	if err := h(w, r); err != nil {
		// handle service error
		e, ok := errors.Cause(err).(*service.Error)
		if ok {
			errRes.HTTPCode = e.HTTPStatusCode()
			errRes.Code = e.Code
			errRes.Error = e.Error()
			if e.Err != nil {
				errRes.Inner = e.Err.Error()
			}
		} else {
			// handle unknown error. would you like to log error ?
			errRes.HTTPCode = http.StatusInternalServerError
			errRes.Code = service.UnknownErrCode
			errRes.Error = err.Error()
		}

		// error response
		if os.Getenv("APP_ENV") == "development" && errRes.HTTPCode >= 500 {
			fmt.Fprintf(w, "%+v", err)
			return
		}
		gores.JSON(w, errRes.HTTPCode, errRes)
	}
}
