package organizers

type Organizer struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	OrgNumber int    `json:"orgnumber"`
	Type 	  string `json:"type"`
}
