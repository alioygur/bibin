package mysql

import (
	"strconv"

	"time"

	"github.com/alioygur/fb-tinder-app/domain"
)

// GenUsers generates number of users
func GenUsers(count int) []*domain.User {
	var users []*domain.User
	now := time.Now().Round(time.Second)
	for i := 1; i <= count; i++ {
		idStr := strconv.Itoa(i)
		var u domain.User
		u.FacebookID = uint64(i)
		u.FirstName = "Jhon " + idStr
		u.LastName = "Doe"
		u.Email = "user" + idStr + "@example.com"
		u.Gender = domain.GenderFemale
		if (i % 2) == 1 {
			u.Gender = domain.GenderMale
		}
		// age = 17 + i
		b := time.Date(now.Year()-(17+i), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		u.Birthday = &b
		u.Status = domain.NewUser
		u.IsAdmin = boolPtr(false)
		u.CreatedAt = now
		u.UpdatedAt = now

		users = append(users, &u)
	}
	return users
}
