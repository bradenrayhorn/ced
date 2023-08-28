package responses

import "github.com/bradenrayhorn/ced/ced"

type Group struct {
	ID           ced.ID `json:"id"`
	Name         string `json:"name"`
	Attendees    uint8  `json:"attendees"`
	MaxAttendees uint8  `json:"max_attendees"`
	HasResponded bool   `json:"has_responded"`
}

func FromGroup(group ced.Group) Group {
	return Group{
		ID:           group.ID,
		Name:         string(group.Name),
		Attendees:    group.Attendees,
		MaxAttendees: group.MaxAttendees,
		HasResponded: group.HasResponded,
	}
}
