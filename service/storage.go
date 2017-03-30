package service

import "github.com/alioygur/fb-tinder-app/domain"

type (
	// FacebookRepository facebook repository interface
	FacebookRepository interface {
		OneByAccessToken(accessToken string) (*domain.User, error)
		Permissions(accessToken string) (map[string]string, error)
		Friends(accessToken string) ([]uint64, error)
		MakeFriend(*domain.FBTestUser, *domain.FBTestUser) error
		ProfilePicture(accessToken string) (string, error)

		PutTestUser(installed bool, perms []string) (*domain.FBTestUser, error)
		FindTestUsers() ([]*domain.FBTestUser, error)
		DeleteTestUser(id string) error
	}

	// Repository database repository interface
	Repository interface {
		AddUser(*domain.User) error
		UserExistsByID(id uint64) (bool, error)
		UserExistsByFacebookID(id uint64) (bool, error)
		UserByFacebookID(id uint64) (*domain.User, error)
		UserByEmail(email string) (*domain.User, error)
		UserByID(id uint64) (*domain.User, error)
		DiscoverPeople(user uint64, gender domain.Gender, ageMin int, ageMax int, limit int) ([]*domain.User, error)
		SyncUserFriendsByFacebookID(id uint64, friends []uint64) error
		UpdateUser(*domain.User) error

		PutReaction(*domain.Reaction) error
		ReactionExistsBy(fromUserID uint64, toUserID uint64, typ domain.ReactionType) (bool, error)

		MakeFriend(user1 uint64, user2 uint64) error
		AreFriends(user1 uint64, user2 uint64) (bool, error)

		PutCredit(*domain.Credit) error
		CalcUserCredits(id uint64) (int, error)

		AddAbuse(*domain.Abuse) error

		PutPicture(*domain.Image) error
		PicturesByUserID(uint64) ([]*domain.Image, error)
		ProfilePicture(uint64) (*domain.Image, error)
		PictureByID(uint64) (*domain.Image, error)
		PictureExistsByUserIDAndIsProfile(uint64, bool) (bool, error)
		UpdatePicture(*domain.Image) error
		DeletePicture(uint64) error
	}
)
