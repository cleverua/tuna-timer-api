package data

import (
	"github.com/pavlo/slack-time/models"
	mgo "gopkg.in/mgo.v2"
)

// TimerService todo
type TimerService struct {
	session *mgo.Session
}

// NewTimerService todo
func NewTimerService(session *mgo.Session) *TimerService {
	return &TimerService{
		session: session,
	}
}

// GetActiveTimer todo
func (s *TimerService) GetActiveTimer(teamID, userID string) (*models.Timer, error) {
	//dao.find
	return nil, nil
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
func (s *TimerService) TaskTotalMinutesForToday(*models.Timer) int {
	// dao.aggregation
	return 0
}

// UserTotalMinutesForToday todo
func (s *TimerService) UserTotalMinutesForToday(userID string) int {
	// dao.aggregation
	return 0
}
