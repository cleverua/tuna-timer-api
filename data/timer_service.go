package data

import (
	"github.com/pavlo/slack-time/models"
	mgo "gopkg.in/mgo.v2"
	"time"
)

// TimerService todo
type TimerService struct {
	session *mgo.Session
	repository *TimerRepository
}

// NewTimerService todo
func NewTimerService(session *mgo.Session) *TimerService {
	return &TimerService{
		session: session,
		repository: NewTimerRepository(session),
	}
}

// GetActiveTimer todo
func (s *TimerService) GetActiveTimer(teamID, userID string) (*models.Timer, error) {
	timer, err := s.repository.findActiveByTeamAndUser(teamID, userID)
	return timer, err
}

// StopTimer todo
func (s *TimerService) StopTimer(timer *models.Timer) error {
	//dao.updateTimer
	return nil
}

// StartTimer todo
func (s *TimerService) StartTimer(teamID, projectID, teamUserID, taskName string) (*models.Timer, error) {
	// dao.create
	return nil, nil
}

// TaskTotalMinutesForToday todo
func (s *TimerService) TotalMinutesForToday(timer *models.Timer) int {

	taskHash := "hash"
	//taskHash := timer.TaskHash
	userID := "userID"
	//userID := timer.TeamUserID
	startDate := time.Now().Truncate(24*time.Hour)
	endDate := time.Now()


	s.repository.totalMinutesForTaskAndUser(taskHash, userID, startDate, endDate)





	return 0
}

// UserTotalMinutesForToday todo
func (s *TimerService) UserTotalMinutesForToday(userID string) int {
	// dao.aggregation
	return 0
}
