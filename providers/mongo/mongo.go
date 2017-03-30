package mongo

import (
	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/fb-tinder-app/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	fakePicture interface {
		PutPicture(*domain.Image) error
		PicturesByUserID(uint64) ([]*domain.Image, error)
		ProfilePicture(uint64) (*domain.Image, error)
		PictureByID(uint64) (*domain.Image, error)
		PictureExistsByUserIDAndIsProfile(uint64, bool) (bool, error)
		UpdatePicture(*domain.Image) error
		DeletePicture(uint64) error
	}
	repository struct {
		fakePicture
		sess *mgo.Session
	}
)

const (
	usersTbl       = `users`
	reactionsTbl   = `reactions`
	creditsTbl     = `credits`
	friendshipsTbl = `friendships`
	matchesTbl     = `matches`
	abusesTbl      = `abuses`
	imagesTbl      = `images`
)

// New instances new repository
func New(sess *mgo.Session) service.Repository {
	return &repository{sess: sess}
}

// c returns mgo collection.
// if sess nil then it uses default session
func (r *repository) c(c string) *mgo.Collection {
	return r.sess.DB("").C(c)
}

func (r *repository) oneBy(c *mgo.Collection, q interface{}, result interface{}) error {
	return c.Find(q).One(result)
}

func (r *repository) oneByID(c *mgo.Collection, id uint64, result interface{}) error {
	return r.oneBy(c, bson.M{"id": id}, result)
}

func (r *repository) existsBy(c *mgo.Collection, q interface{}) (bool, error) {
	n, err := c.Find(q).Limit(1).Count()
	return n > 0, err
}

func wrapErr(err error) error {
	return err
}
