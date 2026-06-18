package users

type User struct {
	Id                   int      `json:"id"`
	Name                 string   `json:"name"`
}

type InterestedOrganizers struct {
	UserID				int 	   	`json:"user_id"`
	OrganizerID			int 		`json:"organizer_id"`
}

