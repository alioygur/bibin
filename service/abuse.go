package service

import "github.com/alioygur/fb-tinder-app/domain"

type (
	// ReportAbuseRequest ...
	ReportAbuseRequest struct {
		UserID   uint64 `json:"-"`        // reporter
		ToUserID uint64 `json:"toUserId"` // reported
		Reason   string `json:"reason"`   // why ?
	}
)

func (s *service) ReportAbuse(r *ReportAbuseRequest) (*domain.Abuse, error) {
	var a domain.Abuse
	a.UserID = r.UserID
	a.ToUserID = r.ToUserID
	a.Reason = r.Reason

	return &a, s.storage.AddAbuse(&a)
}
