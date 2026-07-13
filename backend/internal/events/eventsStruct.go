package events

import (
	"time"	
)

type Events struct {
	Id				int 		`json:"id"`
	Name        	string   	`json:"name"`
	IsExclusive 	bool     	`json:"isExclusive"`
	Date  			time.Time 	`json:"date"`
	Time 			time.Time 	`json:"time"`
	Location 		string 		`json:"location"`
	Description 	string 		`json:"description"`
	IsRegistration 	bool 		`json:"isRegistration"`	    
	// Keywords 	[]string 	`json:"keywords"` 	//TODO: do we have this or? Could be good for potential searches?
}

type EventOrganizer struct {
    EventID     int
    OrganizerID int
}