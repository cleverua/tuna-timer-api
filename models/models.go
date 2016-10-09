package models

import (
	"time"

	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2/bson"
)

const (
	ModelVersionTeam     = 1
	ModelVersionTeamUser = 1
	ModelVersionTimer    = 1
)

// Team represents a Slack team
type Team struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	// slack channel ID or
	// skype group ID or
	// hipchat channel ID
	ExternalSystem   string               `json:"system" bson:"system"`
	ExternalTeamID   string               `json:"ext_id" bson:"ext_id"`
	ExternalTeamName string               `json:"ext_name" bson:"ext_name"`
	Projects         []*Project           `json:"projects" bson:"projects"`
	CreatedAt        time.Time            `json:"created_at" bson:"created_at"`
	SlackOAuth       *slack.OAuthResponse `json:"slack_oauth" bson:"slack_oauth"`
	ModelVersion     int                  `json:"ver" bson:"ver"`
}

// Project - is a project you can associate tasks with and tracks their time. It is embedded in Team
type Project struct {
	ID                  bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ExternalProjectID   string        `json:"ext_id" bson:"ext_id"`
	ExternalProjectName string        `json:"ext_name" bson:"ext_name"`
	CreatedAt           time.Time     `json:"created_at" bson:"created_at"`
}

// TeamUser represents a Slack user that belongs to a team.
// We not going to call it `User` because we may want to have admin users to administer stuff via UI etc
type TeamUser struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TeamID           string        `json:"team_id" bson:"team_id"`
	SlackUserInfo    *slack.User   `json:"slack_user_info" bson:"slack_user_info"`
	ExternalUserID   string        `json:"ext_id" bson:"ext_id"`
	ExternalUserName string        `json:"ext_name" bson:"ext_name"`
	CreatedAt        time.Time     `json:"created_at" bson:"created_at"`
	ModelVersion     int           `json:"ver" bson:"ver"`
}

// Timer - a time record that has start and finish dates. Belongs to a slack user and a task
type Timer struct {
	ID                  bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TeamID              string        `json:"team_id" bson:"team_id"`
	ProjectID           string        `json:"project_id" bson:"project_id"`
	ProjectExternalName string        `json:"project_ext_name" bson:"project_ext_name"`
	ProjectExternalID   string        `json:"project_ext_id" bson:"project_ext_id"`
	TeamUserID          string        `json:"team_user_id" bson:"team_user_id"`
	TeamUserTZOffset    int           `json:"tz_offset" bson:"tz_offset"`
	TaskName            string        `json:"task_name" bson:"task_name"`
	TaskHash            string        `json:"task_hash" bson:"task_hash"`
	CreatedAt           time.Time     `json:"created_at" bson:"created_at"`
	FinishedAt          *time.Time    `json:"finished_at" bson:"finished_at"`
	Minutes             int           `json:"minutes" bson:"minutes"`
	DeletedAt           *time.Time    `json:"deleted_at" bson:"deleted_at"`
	ModelVersion        int           `json:"ver" bson:"ver"`
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
	SubCommand  string
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
	CreatedAt   time.Time
}
