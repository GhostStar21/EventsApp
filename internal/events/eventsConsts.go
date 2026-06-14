package events

import (
	"time"	
)

type Events struct {
	Name        string   	`json:"name"`
	Organizer   []string 	`json:"organizer"`
	IsExclusive bool     	`json:"isExclusive"`
	Date  		time.Time 	`json:"date"`
	Time 		time.Time 	`json:"time"`
	Location 	string 		`json:"location"`
	Description string 		`json:"description"`	    
	// Keywords 	[]string 	`json:"keywords"` 	//TODO: do we have this or? Could be good for potential searches?
}

