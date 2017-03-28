package providers

import (
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/gocart/engine"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func NewJWT() engine.JWTSignParser {
	return &jwt{}
}

type jwt struct{}

func (j *jwt) Sign(claims map[string]interface{}, secret string) (string, error) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, jwtgo.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return tokenString, nil
}

func (j *jwt) Parse(tokenStr string, secret string) (map[string]interface{}, error) {
	token, err := jwtgo.Parse(tokenStr, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		e, ok := err.(*jwtgo.ValidationError)
		if ok {
			return nil, newTokenErr(e)
		}
		return nil, errors.WithStack(err)
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return nil, errors.Errorf("can't get map claims from: %s", tokenStr)
	}
	return claims, nil
}

func newTokenErr(err *jwtgo.ValidationError) error {
	return errors.WithStack(service.NewErr(service.JWTTokenErrCode, err))
}
