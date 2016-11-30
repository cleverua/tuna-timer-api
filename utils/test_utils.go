package utils

import (
	"gopkg.in/mgo.v2"
	"time"
	"github.com/nlopes/slack"
	"github.com/cleverua/tuna-timer-api/models"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

const testTimeParseLayout = "2006 Jan 02 15:04:05"

// TruncateTables - clears database tables, supposed to be run in test's setup method
func TruncateTables(session *mgo.Session) {
	tablesToTruncate := []string{
		MongoCollectionTeams,
		MongoCollectionTimers,
		MongoCollectionTeamUsers,
		MongoCollectionPasses,
	}

	for _, tableName := range tablesToTruncate {
		//log.Printf("Truncating table: %s", tableName)
		session.DB("").C(tableName).RemoveAll(nil)
	}
}

// stands for parse time
func PT(value string) time.Time {
	result, _ := time.Parse(testTimeParseLayout, value)
	return result
}

// Default models for tests
var (
	defaultUser = &models.TeamUser{
		ID:               bson.NewObjectId(),
		TeamID:           "team-id",
		ExternalUserID:   "ext-user-id",
		ExternalUserName: "user-name",
		SlackUserInfo:    &slack.User{
			IsAdmin: true,
		},
	}

	defaultPass = &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "pass-for-jwt-generation",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(5 * time.Minute),
		ClaimedAt:    nil,
		ModelVersion: models.ModelVersionPass,
		TeamUserID:   defaultUser.ID.Hex(),
	}
)

//Insert models for tests into database
func Create(object interface{}, session *mgo.Session) (interface{}, error) {
	var err error

	switch t := object.(type) {
	case *models.TeamUser:
		if object.(*models.TeamUser).ExternalUserID == "" {	object = *defaultUser }
		userCollection := session.DB("").C(MongoCollectionTeamUsers)
		err = userCollection.Insert(object)
	case *models.Pass:
		if object.(*models.Pass).Token == "" { object = *defaultPass }
		passCollection := session.DB("").C(MongoCollectionPasses)
		err = passCollection.Insert(object)
	default:
		_ = t
		err = errors.New("Can't recognize the type of model")
	}
	return object, err
}
