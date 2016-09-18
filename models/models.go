package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Team represents a Slack team
type Team struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	// slack channel ID or
	// skype group ID or
	// hipchat channel ID
	ExternalTeamID   string      `json:"ext_id" bson:"ext_id"`
	ExternalTeamName string      `json:"ext_name" bson:"ext_name"`
	Users            []*TeamUser `json:"users" bson:"users"`
	Projects         []*Project  `json:"projects" bson:"projects"`
	CreatedAt        time.Time   `json:"created_at" bson:"created_at"`
}

// TeamUser represents a Slack user that belongs to a team.
// We not going to call it `User` because we may want to have admin users to administer stuff via UI etc
type TeamUser struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ExternalUserID   string        `json:"ext_id" bson:"ext_id"`
	ExternalUserName string        `json:"ext_name" bson:"ext_name"`
	CreatedAt        time.Time     `json:"created_at" bson:"created_at"`
}

// Project - is a project you can associate tasks with and tracks their time
type Project struct {
	ID                  bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ExternalProjectID   string        `json:"ext_id" bson:"ext_id"`
	ExternalProjectName string        `json:"ext_name" bson:"ext_name"`
	CreatedAt           time.Time     `json:"created_at" bson:"created_at"`
}

// Task - a task that belongs to a project and a user, contains a collection of timers
type Task struct {
	ID           int64
	Name         string  `sql:"size:128"`
	Hash         *string `sql:"size:12"`
	Team         Team    `gorm:"ForeignKey:TeamID"`
	TeamID       int64
	Project      Project `gorm:"ForeignKey:ProjectID"`
	ProjectID    int64
	TotalMinutes int
	Timers       []Timer
}

// Timer - a time record that has start and finish dates. Belongs to a slack user and a task
type Timer struct {
	ID         int64
	TeamUser   TeamUser `gorm:"ForeignKey:TeamUserID"`
	TeamUserID int64
	Task       Task `gorm:"ForeignKey:TaskID"`
	TaskID     int64
	StartedAt  time.Time
	FinishedAt *time.Time
	Minutes    int
	DeletedAt  *time.Time
}

// SlackCustomCommand todo
type SlackCustomCommand struct {
	ID          int64
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
	CreatedAt   time.Time
}
