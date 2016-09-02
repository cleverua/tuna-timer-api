package data

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Team struct {
	gorm.Model

	SlackTeamId string
	CreatedAt time.Time
}

// not going to call it User because we may want to have admin users
// to administer stuff via UI etc
type TeamUser struct {
	gorm.Model

	Email string `sql:"size:128"`
	Team Team
	SlackUserId string `sql:"size:32"`

	Tasks []Task
}

type Project struct {
	gorm.Model

	Name string `sql:"size:32"`
	SlackChannelId string `sql:"size:32"`
	SlackChannelName string `sql:"size:64"`

	Team Team
	Tasks []Task
}

type Task struct {
	gorm.Model

	Name string `sql:"size:128"`
	Hash string `sql:"size:12"`

	Team Team
	User TeamUser
	Project Project

	TotalMinutes int

	Timers []Timer
}

type Timer struct {
	gorm.Model

	Id int64

	Task Task `gorm:"ForeignKey:TaskId"`
	TaskId int64

	StartedAt time.Time
	FinishedAt time.Time

	TotalMinutes int

	DeletedAt time.Time
}
