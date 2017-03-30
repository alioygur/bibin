package domain

import (
	"context"
	"fmt"
	"time"
)

type (
	// AccountType ...
	AccountType uint8
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
		ID          uint64      `bson:"id"`
		FacebookID  uint64      `bson:"facebook_id"`
		FirstName   string      `bson:"first_name"`
		LastName    string      `bson:"last_name"`
		Email       string      `bson:"email"`
		Gender      Gender      `bson:"gender"`
		Birthday    *time.Time  `bson:"birthday"`
		Status      UserStatus  `bson:"status"`
		IsAdmin     *bool       `bson:"is_admin"`
		AccountType AccountType `bson:"account_type"`
		CreatedAt   time.Time   `bson:"created_at"`
		UpdatedAt   time.Time   `bson:"updated_at"`

		Friends []*User  `json:",omitempty" bson:"-"`
		Images  []*Image `json:",omitempty" bson:"images"`
	}

	// Image entity
	Image struct {
		ID        uint64    `bson:"-"`    // todo: delete this field
		UserID    uint64    `bson:"-"`    // todo: delete this field
		Name      string    `bson:"name"` // uniq name
		IsProfile bool      `bson:"is_profile"`
		CreatedAt time.Time `bson:"created_at"`
	}

	// FBTestUser entity
	FBTestUser struct {
		ID          string `json:"id"`
		LoginURL    string `json:"login_url"`
		AccessToken string `json:"access_token"`
	}

	// Reaction entity
	Reaction struct {
		ID         uint64       `bson:"id"`
		FromUserID uint64       `bson:"from_user_id"`
		ToUserID   uint64       `bson:"to_user_id"`
		Type       ReactionType `bson:"type"`
		CreatedAt  time.Time    `bson:"created_at"`
	}

	// Credit ...
	Credit struct {
		ID        uint64                 `bson:"id"`
		UserID    uint64                 `bson:"user_id"`
		Amount    int                    `bson:"amount"`
		Type      CreaditTransactionType `bson:"type"`
		Desc      string                 `bson:"desc"`
		CreatedAt time.Time              `bson:"created_at"`
	}

	// Abuse ...
	Abuse struct {
		ID        uint64    `bson:"id"`
		UserID    uint64    `bson:"user_id"`    // reporter
		ToUserID  uint64    `bson:"to_user_id"` // reported
		Reason    string    `bson:"reason"`
		CreatedAt time.Time `bson:"created_at"`
	}
)

// Account Types
const (
	FreeAccount    AccountType = iota
	PremiumAccount AccountType = iota
)

// Genders
const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

// User statuses
const (
	StatusNewUser UserStatus = iota
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

// NewUser instances new User with default values
func NewUser() *User {
	var u User
	var f bool
	u.Gender = GenderUnknown
	u.AccountType = FreeAccount
	u.Status = StatusNewUser
	u.IsAdmin = &f
	u.Images = make([]*Image, 0)
	u.Friends = make([]*User, 0)
	return &u
}

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

// SetAge sets user's age.
func (u *User) SetAge(age int) {
	now := time.Now()
	b := time.Date(now.Year()-age, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	u.Birthday = &b
}
