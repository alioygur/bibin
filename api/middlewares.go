package api

import (
	"net/http"
	"strings"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/pkg/errors"
)

type (
	userFinder interface {
		GetFromAuthToken(tokenStr string) (*domain.User, error)
	}
)

// getToken gets Authorization key from headers
func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		err := service.NewErr(service.UnknownErrCode, errors.New("invalid authorization token format"))
		return "", errors.WithStack(err)
	}

	return authHeaderParts[1], nil
}

func newSetUserMid(uf userFinder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return errHandler(func(w http.ResponseWriter, r *http.Request) error {
			tokenStr, err := getToken(r)
			if err != nil {
				return err
			}

			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return nil
			}

			u, err := uf.GetFromAuthToken(tokenStr)
			if err != nil {
				return err
			}

			ctx := u.NewContext(r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))

			return nil
		})
	}
}

func newAuthRequiredMid(next http.Handler) http.Handler {
	return errHandler(func(w http.ResponseWriter, r *http.Request) error {
		_, ok := domain.UserFromContext(r.Context())
		if !ok {
			err := service.NewErr(service.UnknownErrCode, errors.New("auth required")).SetHTTPStatusCode(http.StatusUnauthorized)
			return errors.WithStack(err)
		}

		next.ServeHTTP(w, r)
		return nil
	})
}

func newAdminOnlyMid(next http.Handler) http.Handler {
	return errHandler(func(w http.ResponseWriter, r *http.Request) error {
		usr, ok := domain.UserFromContext(r.Context())
		if !ok {
			return errors.New("user can't get from request's context")
		}

		if *usr.IsAdmin != true {
			err := service.NewErr(service.UnknownErrCode, errors.New("admin only")).SetHTTPStatusCode(http.StatusUnauthorized)
			return errors.WithStack(err)
		}

		next.ServeHTTP(w, r)
		return nil
	})
}
