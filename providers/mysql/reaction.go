package mysql

import (
	"time"

	"github.com/alioygur/fb-tinder-app/domain"
)

func (r *repository) PutReaction(rr *domain.Reaction) error {
	rr.CreatedAt = time.Now()
	_, err := r.insert(reactionsTbl, rr)
	return err
}

func (r *repository) ReactionExistsBy(fromUserID uint64, toUserID uint64, typ domain.ReactionType) (bool, error) {
	return r.existsBy(reactionsTbl, "from_user_id=? AND to_user_id=? AND type=?", fromUserID, toUserID, typ)
}
