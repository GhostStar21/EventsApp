package users

type User struct {
	Id                   int      `json:"id"`
	Name                 string   `json:"name"`
	InterestedOrganizers []string `json:"interestedorganizers"`
}
