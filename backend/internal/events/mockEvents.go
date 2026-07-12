package events

import (
	"time"
)

var events = []Events{
	{
		Id:          1,
		Name:        "All dress party",
		IsExclusive: false,
		Date:        time.Date(2026, time.June, 20, 0, 0, 0, 0, time.UTC),
		Time:        time.Date(0, 1, 1, 18, 30, 0, 0, time.UTC),
		Location:    "Gjovik, Norway",
		Description: "An evening meetup for Go developers to network and discuss the latest trends in Go programming.",
	},
	{
		Id:          2,
		Name:        "Halloween party",
		IsExclusive: false,
		Date:        time.Date(2026, time.June, 21, 0, 0, 0, 0, time.UTC),
		Time:        time.Date(0, 1, 1, 18, 30, 0, 0, time.UTC),
		Location:    "Gjovik, Norway",
		Description: "You better dress up",
	},
}
