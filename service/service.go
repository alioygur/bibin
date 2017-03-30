package service

import (
	"io"

	"github.com/alioygur/fb-tinder-app/domain"
)

type (
	// MailSender interface
	MailSender interface {
		Send(to []string, subject string, body []byte) error
	}

	// Validator interface
	Validator interface {
		CheckEmail(string) error
		CheckRequired(string, string) error
		CheckStringLen(s string, min int, max int, field string) error
	}

	// JWTSignParser interface
	JWTSignParser interface {
		Sign(claims map[string]interface{}, secret string) (string, error)
		Parse(tokenStr string, secret string) (map[string]interface{}, error)
	}

	// ImageCDN image cdn interface, uploader, getter e.g.
	ImageCDN interface {
		Upload(io.Reader) (*domain.Image, error)
		UploadURL(string) (*domain.Image, error)
	}

	// Service interface
	Service interface {
		Register(*RegisterRequest) (*domain.User, error)
		React(*ReactRequest) (bool, error)
		Show(*ShowUserRequest) (*domain.User, error)
		UpdateUser(*UpdateUserRequest) error
		GetFromAuthToken(tokenStr string) (*domain.User, error)
		GenToken(*domain.User, TokenType) (string, error)

		DiscoverPeople(*DiscoverPeopleRequest) ([]*domain.User, error)

		ReportAbuse(*ReportAbuseRequest) (*domain.Abuse, error)

		Pictures(*PicturesRequest) ([]*domain.Image, error)
		UploadPicture(*UploadPictureRequest) (*domain.Image, error)
		SetProfilePicture(userID uint64, picID uint64) error
		DeletePicture(*DeletePictureRequest) error
	}

	// user implementation of User
	service struct {
		fb      FacebookRepository
		storage Repository
		jwt     JWTSignParser
		imgCDN  ImageCDN
	}
)

// New instances new service
func New(fb FacebookRepository, storage Repository, jwt JWTSignParser, imgCDN ImageCDN) Service {
	return &service{fb: fb, storage: storage, jwt: jwt, imgCDN: imgCDN}
}

func boolPtr(v bool) *bool {
	return &v
}
