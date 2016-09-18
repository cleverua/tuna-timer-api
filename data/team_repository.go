package data

import (
	"time"

	"github.com/pavlo/slack-time/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const collectionName = "teams"

// TeamRepository todo
type TeamRepository struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// TeamRepositoryInterface describes the functionality this repo provides.
// It has two implementations `TeamRepository` and `testTeamRepository`.
// The latter used to mimic/test error cases
type TeamRepositoryInterface interface {
	findByExternalID(externalTeamID string) (*models.Team, error)
	createTeam(externalID, externalName string) (*models.Team, error)
	addProject(team *models.Team, externalProjectID, externalProjectName string) error
	addUser(team *models.Team, externalUserID, externalUserName string) error
}

// NewTeamRepository is a factory method
func NewTeamRepository(session *mgo.Session) *TeamRepository {
	return &TeamRepository{
		session:    session,
		collection: session.DB("").C(collectionName),
	}
}

func (r *TeamRepository) findByExternalID(externalTeamID string) (*models.Team, error) {
	team := &models.Team{}
	err := r.collection.Find(bson.M{"ext_id": externalTeamID}).One(team)

	if err != nil && err == mgo.ErrNotFound {
		team = nil
		err = nil
	}
	return team, err
}

// CreateTeam creates a new team
func (r *TeamRepository) createTeam(externalID, externalName string) (*models.Team, error) {

	team := &models.Team{
		ID:               bson.NewObjectId(),
		CreatedAt:        time.Now(),
		ExternalTeamID:   externalID,
		ExternalTeamName: externalName,
		Projects:         []*models.Project{},
		Users:            []*models.TeamUser{},
	}

	err := r.collection.Insert(team)
	return team, err
}

func (r *TeamRepository) addProject(team *models.Team, externalProjectID, externalProjectName string) error {
	testTeam := &models.Team{}
	err := r.collection.Find(bson.M{"projects.ext_id": externalProjectID}).One(testTeam)
	if err != nil && err == mgo.ErrNotFound {

		project := &models.Project{
			ID:                  bson.NewObjectId(),
			ExternalProjectID:   externalProjectID,
			ExternalProjectName: externalProjectName,
			CreatedAt:           time.Now(),
		}

		err = r.collection.Update(bson.M{"_id": team.ID}, bson.M{"$push": bson.M{"projects": project}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepository) addUser(team *models.Team, externalUserID, externalUserName string) error {
	testTeam := &models.Team{}
	err := r.collection.Find(bson.M{"users.ext_id": externalUserID}).One(testTeam)
	if err != nil && err == mgo.ErrNotFound {
		user := models.TeamUser{
			ID:               bson.NewObjectId(),
			ExternalUserID:   externalUserID,
			ExternalUserName: externalUserName,
			CreatedAt:        time.Now(),
		}

		err = r.collection.Update(bson.M{"_id": team.ID}, bson.M{"$push": bson.M{"users": user}})
		if err != nil {
			return err
		}
	}
	return nil
}
