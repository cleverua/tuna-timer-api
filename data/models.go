package data

import (
	"time"
)

type Team struct {
	Id int64
	SlackTeamId string
	CreatedAt time.Time
}

// not going to call it User because we may want to have admin users
// to administer stuff via UI etc
type TeamUser struct {
	Id int64

	Email string
	Team Team `gorm:"ForeignKey:TeamId"`
	TeamId int64
	SlackUserId string

	Tasks []Task
}

type Project struct {

	Id int64

	Name string `sql:"size:64"`
	SlackChannelId string `sql:"size:32"`
	SlackChannelName string `sql:"size:64"`

	Team Team `gorm:"ForeignKey:TeamId"`
	TeamId int64

	Tasks []Task
}

type Task struct {
	Id int64

	Name string `sql:"size:128"`
	Hash string `sql:"size:12"`

	Team Team `gorm:"ForeignKey:TeamId"`
	TeamId int64

	Project Project `gorm:"ForeignKey:ProjectId"`
	ProjectId int64

	TotalMinutes int

	Timers []Timer
}

type Timer struct {
	Id int64

	//User TeamUser `gorm:"ForeignKey:TeamUserId"`
	TeamUserId int64

	//Task Task `gorm:"ForeignKey:TaskId"`
	TaskId int64

	StartedAt time.Time
	FinishedAt time.Time

	Minutes int

	DeletedAt time.Time
}
