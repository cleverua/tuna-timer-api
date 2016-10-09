package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const usersCollectionName = "team_users"

type UserRepository struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func NewUserRepository(session *mgo.Session) *UserRepository {
	return &UserRepository{
		session:    session,
		collection: session.DB("").C(usersCollectionName),
	}
}

func (r *UserRepository) FindByExternalID(externalUserID string) (*models.TeamUser, error) {
	teamUser := &models.TeamUser{}
	err := r.collection.Find(bson.M{"ext_id": externalUserID}).One(teamUser)

	if err != nil && err == mgo.ErrNotFound {
		teamUser = nil
		err = nil
	}
	return teamUser, err
}

func (r *UserRepository) save(user *models.TeamUser) (*models.TeamUser, error) {
	if user.ID == "" {
		user.ID = bson.NewObjectId()
		user.CreatedAt = time.Now()
		user.ModelVersion = models.ModelVersionTeamUser

		return user, r.collection.Insert(user)
	}
	return user, r.collection.Update(bson.M{"_id": user.ID}, user)
}
