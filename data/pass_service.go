package data

import (
	"github.com/satori/go.uuid"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PassService struct {
	repository *PassRepository
}

func NewPassService(session *mgo.Session) *PassService {
	return &PassService{
		repository: NewPassRepository(session),
	}
}

func (s *PassService) EnsurePass(team *models.Team, user *models.TeamUser, project *models.Project) (*models.Pass, error) {
	pass, _ := s.repository.FindActiveByUserID(user.ID.Hex())

	if pass == nil {
		return s.createPass(team, user, project.ID.Hex())
	}

	pass.ExpiresAt = time.Now().Add(utils.PassExpiresInMinutes * time.Minute)
	err := s.repository.update(pass)

	return pass, err
}

func (s *PassService) RemoveStalePasses() error {
	err := s.repository.removeExpiredPasses()
	if err != nil {
		return err
	}

	return s.repository.removePassesClaimedBefore(time.Now().Add(
		-utils.ClaimedPassesToPurgeAfterDays * 24 * 60 * time.Minute))
}

func (s *PassService) createPass(team *models.Team, user *models.TeamUser, projectID string) (*models.Pass, error) {
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
