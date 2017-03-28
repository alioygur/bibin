package domain

import (
	"context"
	"fmt"
	"time"
)

type (
	// Gender type
	Gender uint8
	// UserStatus user status
	UserStatus uint8
	// ReactionType ...
	ReactionType uint8
	// CreaditTransactionType ...
	CreaditTransactionType uint8

	contextKey string

	// User entity
	User struct {
		ID         uint64     `db:"id"`
		FacebookID uint64     `db:"facebook_id"`
		FirstName  string     `db:"first_name"`
		LastName   string     `db:"last_name"`
		Email      string     `db:"email"`
		Gender     Gender     `db:"gender"`
		Birthday   *time.Time `db:"birthday"`
		Status     UserStatus `db:"status"`
		IsAdmin    *bool      `db:"is_admin"`
		CreatedAt  time.Time  `db:"created_at"`
		UpdatedAt  time.Time  `db:"updated_at"`

		Images []*Image
	}

	// Image entity
	Image struct {
		ID        uint64    `db:"id"`
		UserID    uint64    `db:"user_id"`
		Name      string    `db:"name"` // uniq name
		IsProfile *bool     `db:"is_profile"`
		CreatedAt time.Time `db:"created_at"`
	}

	// FBTestUser entity
	FBTestUser struct {
		ID          string `json:"id"`
		LoginURL    string `json:"login_url"`
		AccessToken string `json:"access_token"`
	}

	// Reaction entity
	Reaction struct {
		ID         uint64       `db:"id"`
		FromUserID uint64       `db:"from_user_id"`
		ToUserID   uint64       `db:"to_user_id"`
		Type       ReactionType `db:"type"`
		CreatedAt  time.Time    `db:"created_at"`
	}

	// Credit ...
	Credit struct {
		ID        uint64                 `db:"id"`
		UserID    uint64                 `db:"user_id"`
		Amount    int                    `db:"amount"`
		Type      CreaditTransactionType `db:"type"`
		Desc      string                 `db:"desc"`
		CreatedAt time.Time              `db:"created_at"`
	}

	// Abuse ...
	Abuse struct {
		ID        uint64    `db:"id"`
		UserID    uint64    `db:"user_id"`    // reporter
		ToUserID  uint64    `db:"to_user_id"` // reported
		Reason    string    `db:"reason"`
		CreatedAt time.Time `db:"created_at"`
	}
)

// Genders
const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

// User statuses
const (
	NewUser UserStatus = iota
)

// React types
const (
	ReactUnlike ReactionType = iota
	ReactLike
)

// credit transaction types
const (
	_ CreaditTransactionType = iota
	CreditGift
	CreditTransactionReact
)

const (
	userContextKey contextKey = "user"
)

func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// UserFromContext gets user from context
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

// UserMustFromContext gets user from context. if can't make panic
func UserMustFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		panic("user can't get from request's context")
	}
	return u
}

// ProfilePicture returns profile picture url
func (u *User) ProfilePicture() string {
	return fmt.Sprintf("https://graph.facebook.com/%d/picture?type=%s", u.ID, "large")
}
