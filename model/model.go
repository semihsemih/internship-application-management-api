package model

import "time"

type Candidate struct {
	ID              string      `json:"id" bson:"_id,omitempty"`
	FirstName       string      `json:"first_name" bson:"first_name"`
	LastName        string      `json:"last_name" bson:"last_name"`
	Email           string      `json:"email" bson:"email"`
	Department      string      `json:"department" bson:"department"`
	University      string      `json:"university" bson:"university"`
	Experience      bool        `json:"experience" bson:"experience"`
	ApplicationDate time.Time   `json:"application_date,omitempty" bson:"application_date"`
	Status          string      `json:"status,omitempty" bson:"status"`
	MeetingCount    int         `json:"meeting_count,omitempty" bson:"meeting_count"`
	NextMeeting     *time.Time  `json:"next_meeting,omitempty" bson:"next_meeting"`
	Assignee        interface{} `json:"assignee,omitempty" bson:"assignee"`
}

type Assignee struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Name       string `json:"name" bson:"name"`
	Department string `json:"department" bson:"department"`
}
