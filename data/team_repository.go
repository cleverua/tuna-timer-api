package data

import (
	"time"

	"github.com/tuna-timer/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const teamsCollectionName = "teams"

// TeamRepository todo
type TeamRepository struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// TeamRepositoryInterface describes the functionality this repo provides.
// It has two implementations `TeamRepository` and `testTeamRepository`.
// The latter used to mimic/test error cases
type TeamRepositoryInterface interface {
	FindByExternalID(externalTeamID string) (*models.Team, error)
	save(team *models.Team) error
	createTeam(externalID, externalName string) (*models.Team, error)
	addProject(team *models.Team, externalProjectID, externalProjectName string) error
	//addUser(team *models.Team, externalUserID, externalUserName string) error
}

// NewTeamRepository is a factory method
func NewTeamRepository(session *mgo.Session) *TeamRepository {
	return &TeamRepository{
		session:    session,
		collection: session.DB("").C(teamsCollectionName),
	}
}

func (r *TeamRepository) FindByExternalID(externalTeamID string) (*models.Team, error) {
	team := &models.Team{}
	err := r.collection.Find(bson.M{"ext_id": externalTeamID}).One(team)

	if err != nil && err == mgo.ErrNotFound {
		team = nil
		err = nil
	}
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

// CreateTeam creates a new team - this method used for tests only!
func (r *TeamRepository) createTeam(externalID, externalName string) (*models.Team, error) {

	team := &models.Team{
		ID:               bson.NewObjectId(),
		CreatedAt:        time.Now(),
		ExternalTeamID:   externalID,
		ExternalTeamName: externalName,
		Projects:         []*models.Project{},
		ModelVersion:     models.ModelVersionTeam,
	}

	err := r.collection.Insert(team)
	return team, err
}

func (r *TeamRepository) save(team *models.Team) error {
	if team.ID == "" {
		team.ID = bson.NewObjectId()
		team.CreatedAt = time.Now()
		team.ModelVersion = models.ModelVersionTeam
		return r.collection.Insert(team)
	}

	log.Printf("Updating: %+v", team)
	return r.collection.Update(bson.M{"_id": team.ID}, team)
}
