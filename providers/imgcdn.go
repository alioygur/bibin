package providers

import (
	"io"

	"github.com/alioygur/cloudinary-go"
	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"
	"github.com/pkg/errors"
)

type (
	img struct {
		client *cloudinary.Client
	}
)

// NewImageCDN instances a impl of service.ImageCDN
func NewImageCDN(c *cloudinary.Client) service.ImageCDN {
	return &img{c}
}

func (i *img) Upload(r io.Reader) (*domain.Image, error) {
	res, err := i.client.Upload(r, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return mapToDomainImage(res), nil
}

func (i *img) UploadURL(imgURL string) (*domain.Image, error) {
	res, err := i.client.Fetch(imgURL, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return mapToDomainImage(res), nil
}

func mapToDomainImage(img *cloudinary.UploadResponse) *domain.Image {
	var dimg domain.Image
	dimg.Name = img.PublicID
	return &dimg
}
