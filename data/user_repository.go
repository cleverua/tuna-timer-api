package data

import (
	"github.com/cleverua/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"errors"
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

func (r *UserRepository) FindByID(userID string) (*models.TeamUser, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("id is not valid")
	}

	teamUser := &models.TeamUser{}
	err := r.collection.FindId(bson.ObjectIdHex(userID)).One(&teamUser)

	if err != nil && err == mgo.ErrNotFound {
		teamUser = nil
		err = nil
	}
	return teamUser, err
}

func (r *UserRepository) Save(user *models.TeamUser) (*models.TeamUser, error) {
	if user.ID == "" {
		user.ID = bson.NewObjectId()
		user.CreatedAt = time.Now()
		user.ModelVersion = models.ModelVersionTeamUser

		return user, r.collection.Insert(user)
	}
	return user, r.collection.Update(bson.M{"_id": user.ID}, user)
}
