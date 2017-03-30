package mongo

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/pkg/errors"
)

func (r *repository) PutReaction(rr *domain.Reaction) error {
	rr.CreatedAt = time.Now()
	err := r.c(reactionsTbl).Insert(rr)
	return errors.WithStack(wrapErr(err))
}

func (r *repository) ReactionExistsBy(fromUserID uint64, toUserID uint64, typ domain.ReactionType) (bool, error) {
	return r.existsBy(r.c(reactionsTbl), bson.M{"from_user_id": fromUserID, "to_user_id": toUserID, "type": typ})
}
