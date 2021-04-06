package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/ent/schema"
)

type User struct {
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
	DeletedAt         *time.Time            `json:"deletedAt"`
	UUID              uuid.UUID             `json:"uuid"`
	Email             string                `json:"email"`
	NotificationEmail string                `json:"notificationEmail"`
	EmailVerified     bool                  `json:"emailVerified"`
	ProviderID        string                `json:"providerId"`
	ProviderType      *schema.OAuthProvider `json:"providerType"`
	FirstName         string                `json:"firstName"`
	LastName          string                `json:"lastName"`
}

func NewFromUser(u *ent.User) *User {
	return &User{
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		DeletedAt:     u.DeletedAt,
		UUID:          u.UUID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		ProviderID:    u.ProviderID,
		ProviderType:  u.ProviderType,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
	}
}

func NewFromUsers(users []*ent.User) []*User {
	dtoUsers := make([]*User, len(users))
	for i := range dtoUsers {
		dtoUsers[i] = NewFromUser(users[i])
	}
	return dtoUsers
}
