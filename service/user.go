package service

import (
	"fmt"
	"os"
	"time"

	"github.com/alioygur/fb-tinder-app/domain"
	"github.com/alioygur/goutil"
	"github.com/pkg/errors"
)

type (
	// TokenType ...
	TokenType uint8

	// RegisterRequest ...
	RegisterRequest struct {
		AccessToken string
	}
	// ReactRequest ...
	ReactRequest struct {
		FromUserID uint64
		ToUserID   uint64
		Type       domain.ReactionType
	}

	// ShowUserRequest ...
	ShowUserRequest struct {
		ID uint64
	}

	// UpdateUserRequest ...
	UpdateUserRequest struct {
		ID       uint64 `json:"-"`
		Birthday string `json:"birthday"`
	}

	// DiscoverPeopleRequest ...
	DiscoverPeopleRequest struct {
		UserID uint64 // discoverer
		Gender domain.Gender
		AgeMin int
		AgeMax int
		Limit  int
	}
)

// Token Types
const (
	AuthToken TokenType = iota
	ActivationToken
	PasswordResetToken
)

func (s *service) Register(r *RegisterRequest) (*domain.User, error) {
	// check permissions granted
	perms, err := s.fb.Permissions(r.AccessToken)
	if err != nil {
		return nil, err
	}

	requiredPerms := goutil.EnvMustSliceStr("FB_REQUIRED_PERMS", ",")

	for _, p := range requiredPerms {
		status, ok := perms[p]
		if ok && status == "granted" {
			continue
		}
		err := NewErr(PermissionNotGrantedErrCode, fmt.Errorf("%s permission not granted", p))
		return nil, errors.WithStack(err)
	}

	usr, err := s.fb.OneByAccessToken(r.AccessToken)
	if err != nil {
		return nil, err
	}

	// already registered?
	exists, err := s.storage.UserExistsByFacebookID(usr.FacebookID)
	if err != nil {
		return nil, err
	}
	if !exists {
		// also get profile pic from facebook.
		picURL, err := s.fb.ProfilePicture(r.AccessToken)
		if err != nil {
			return nil, err
		}

		img, err := s.imgCDN.UploadURL(picURL)
		if err != nil {
			return nil, err
		}
		img.IsProfile = true
		usr.Images = append(usr.Images, img)

		if usr.IsAdmin == nil {
			usr.IsAdmin = boolPtr(false)
		}

		if err := s.storage.AddUser(usr); err != nil {
			return nil, err
		}

		// // add credit.
		// var c domain.Credit
		// c.Amount = 100
		// c.UserID = usr.ID
		// c.Type = domain.CreditGift
		// if err := s.storage.PutCredit(&c); err != nil {
		// 	if err := s.storage.Rollback(); err != nil {
		// 		return nil, err
		// 	}
		// 	return nil, err
		// }
	}

	usr, err = s.storage.UserByFacebookID(usr.FacebookID)

	return usr, err
}

// React creates a reaction from a user to a user
// if there is a matches between users then return true, nil
func (s *service) React(r *ReactRequest) (bool, error) {
	react := domain.Reaction{
		FromUserID: r.FromUserID,
		ToUserID:   r.ToUserID,
		Type:       r.Type,
	}

	// from user and to user are same ?
	if r.FromUserID == r.ToUserID {
		err := NewErr(UnknownErrCode, errors.New("users are same"))
		return false, errors.WithStack(err)
	}

	// // from user and to user are exists ?
	// users := []uint64{r.FromUserID, r.ToUserID}
	// for _, id := range users {
	// 	exists, err := s.storage.UserExistsByID(id)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	if !exists {
	// 		return false, errors.WithStack(NewErr(NotFoundErrCode, nil))
	// 	}
	// }

	// // check credit
	// credit, err := s.storage.CalcUserCredits(r.FromUserID)
	// if err != nil {
	// 	return false, err
	// }
	// if credit < 1 {
	// 	return false, errors.WithStack(NewErr(NoMoreCreditErrCode, nil))
	// }

	// // add credit transaction
	// ct := domain.Credit{
	// 	UserID: r.FromUserID,
	// 	Type:   domain.CreditTransactionReact,
	// 	Amount: -1,
	// 	Desc:   fmt.Sprintf("react type: %v user: %v", r.Type, r.ToUserID),
	// }
	// if err := s.storage.PutCredit(&ct); err != nil {
	// 	if err := s.storage.Rollback(); err != nil {
	// 		return false, err
	// 	}
	// 	return false, err
	// }

	// already in matches table
	exists, err := s.storage.AreFriends(r.FromUserID, r.ToUserID)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	if err := s.storage.PutReaction(&react); err != nil {
		return false, err
	}

	if r.Type == domain.ReactUnlike {
		return false, nil
	}

	// to user also liked from user before?
	matched, err := s.storage.ReactionExistsBy(r.ToUserID, r.FromUserID, domain.ReactLike)
	if err != nil {
		return false, err
	}

	// yes? then put to matches table
	if matched {
		if err := s.storage.MakeFriend(r.FromUserID, r.ToUserID); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (s *service) Show(r *ShowUserRequest) (*domain.User, error) {
	return s.storage.UserByID(r.ID)
}

func (s *service) GetFromAuthToken(tokenStr string) (*domain.User, error) {
	email, err := s.getEmailFromToken(tokenStr, AuthToken)
	if err != nil {
		return nil, err
	}

	return s.storage.UserByEmail(email)
}

func (s *service) UpdateUser(r *UpdateUserRequest) error {
	u, err := s.storage.UserByID(r.ID)
	if err != nil {
		return err
	}

	if r.Birthday != "" {
		t, err := time.Parse("02.01.2006", r.Birthday)
		if err != nil {
			return errors.WithStack(NewErr(ValidationErrCode, err))
		}
		u.Birthday = &t
	}

	return s.storage.UpdateUser(u)
}

func (s *service) DiscoverPeople(r *DiscoverPeopleRequest) ([]*domain.User, error) {
	if r.UserID == 0 {
		return nil, errors.New("you must set UserID")
	}

	if r.AgeMin < 18 || r.AgeMin > 100 {
		err := NewErr(ValidationErrCode, errors.New("ageMin must between 18 and 100"))
		return nil, errors.WithStack(err)
	}

	if r.AgeMax < 18 || r.AgeMax > 100 {
		err := NewErr(ValidationErrCode, errors.New("ageMax must between 18 and 100"))
		return nil, errors.WithStack(err)
	}

	if r.AgeMin > r.AgeMax {
		err := NewErr(ValidationErrCode, errors.New("ageMax must bigger than ageMin"))
		return nil, errors.WithStack(err)
	}

	return s.storage.DiscoverPeople(r.UserID, r.Gender, r.AgeMin, r.AgeMax, r.Limit)
}

func (s *service) GenToken(usr *domain.User, t TokenType) (string, error) {
	claims := map[string]interface{}{
		"type":  t,
		"email": usr.Email,
	}
	switch t {
	case AuthToken:
		claims["exp"] = time.Now().Add(time.Hour * 6).Unix()
	case ActivationToken:
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	case PasswordResetToken:
		claims["exp"] = time.Now().Add(time.Hour * 3).Unix()
	default:
		return "", errors.Errorf("undefined token type %v", t)
	}
	return s.jwt.Sign(claims, os.Getenv("SECRET_KEY"))
}

func (s *service) getEmailFromToken(token string, t TokenType) (string, error) {
	claims, err := s.jwt.Parse(token, os.Getenv("SECRET_KEY"))
	if err != nil {
		return "", err
	}

	if ct, ok := claims["type"].(float64); ok != true || TokenType(ct) != t {
		return "", errors.Errorf("invalid token type %v", t)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.Errorf("email can't get from token claims: %v", claims)
	}

	return email, nil
}
