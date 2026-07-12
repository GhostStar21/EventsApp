package users

import (
	"EventsApp/internal/consts"
)

type User struct {
	Id           int         `json:"id"`
	Name         string      `json:"name"`
	Email        string      `json:"email,omitempty"`
	Role         consts.Role `json:"role"`
	PasswordHash string      `json:"-"`
}

type InterestedOrganizers struct {
	UserID      int `json:"user_id"`
	OrganizerID int `json:"organizer_id"`
}
