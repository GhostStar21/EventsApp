package events

import (
	"time"	
)

type Events struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	IsExclusive      bool      `json:"isExclusive"`
	Date             time.Time `json:"event_date"`
	Time             time.Time `json:"event_time"`
	Location         string    `json:"location"`
	Description      string    `json:"description"`
	IsRegistration   bool      `json:"isRegistration"`
	RegistrationLink string    `json:"registrationLink,omitempty"`
	OrganizerId      *int      `json:"organizer_id,omitempty"`
	OrganizerName    string    `json:"organizer_name,omitempty"`
}

type EventOrganizer struct {
    EventID     int
    OrganizerID int
}