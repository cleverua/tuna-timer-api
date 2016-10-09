package data

import (
	"github.com/satori/go.uuid"
	"github.com/tuna-timer/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/tuna-timer/tuna-timer-api/utils"
)

type PassService struct {
	repository *PassRepository
}

func NewPassService(session *mgo.Session) *PassService {
	return &PassService{
		repository: NewPassRepository(session),
	}
}

func (s *PassService) CreatePass(team *models.Team, user *models.TeamUser, projectID string) (*models.Pass, error) {
	now := time.Now()

	pass := &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        uuid.NewV4().String(),
		ProjectID:    projectID,
		TeamID:       team.ID.Hex(),
		TeamUserID:   user.ID.Hex(),
		CreatedAt:    now,
		ExpiresAt:    now.Add(utils.PassExpiresInMinutes * time.Minute),
		ClaimedAt:    nil,
		ModelVersion: models.ModelVersionPass,
	}

	return pass, s.repository.insert(pass)
}
