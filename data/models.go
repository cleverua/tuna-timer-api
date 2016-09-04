package data

import (
	"time"
)

// Team represents a Slack team
type Team struct {
	ID          int64
	SlackTeamID string //todo: should be unique
	TeamUsers   []TeamUser
	Projects    []Project
	CreatedAt   time.Time
}

// TeamUser represents a Slack user that belongs to a team.
// We not going to call it `User` because we may want to have admin users to administer stuff via UI etc
type TeamUser struct {
	ID          int64
	Name        string
	Team        Team `gorm:"ForeignKey:TeamID"`
	TeamID      int64
	SlackUserID string
}

// Project - is a project you can associate tasks with and tracks their time
type Project struct {
	ID int64

	Name             string `sql:"size:64"`
	SlackChannelID   string `sql:"size:32"`
	SlackChannelName string `sql:"size:64"`

	Team   Team `gorm:"ForeignKey:TeamID"`
	TeamID int64

	Tasks []Task
}

// Task - a task that belongs to a project and a user, contains a collection of timers
type Task struct {
	ID int64

	Name string `sql:"size:128"`
	Hash string `sql:"size:12"`

	Team   Team `gorm:"ForeignKey:TeamID"`
	TeamID int64

	Project   Project `gorm:"ForeignKey:ProjectID"`
	ProjectID int64

	TotalMinutes int

	Timers []Timer
}

// Timer - a time record that has start and finish dates. Belongs to a slack user and a task
type Timer struct {
	ID int64

	TeamUser   TeamUser `gorm:"ForeignKey:TeamUserID"`
	TeamUserID int64

	Task   Task `gorm:"ForeignKey:TaskID"`
	TaskID int64

	StartedAt  time.Time
	FinishedAt time.Time

	Minutes int

	DeletedAt time.Time
}
