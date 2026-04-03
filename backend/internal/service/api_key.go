package service

import "time"

type APIKey struct {
	ID          int64
	UserID      int64
	Key         string
	Name        string
	GroupID     *int64
	OrgID          *int64
	OrgProjectID   *int64
	Status      string
	IPWhitelist []string
	IPBlacklist []string
	UsageLimit  *float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        *User
	Group       *Group
	Organization *Organization
	OrgProject   *OrgProject
}

func (k *APIKey) IsActive() bool {
	return k.Status == StatusActive
}
