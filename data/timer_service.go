package data

import (
	"github.com/pavlo/slack-time/models"
	"gopkg.in/mgo.v2"
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

// TotalMinutesForTaskToday calculates the total number of minutes the user was/is working on particular task today
func (s *TimerService) TotalMinutesForTaskToday(timer *models.Timer) int {
	endDate := time.Now()
	startDate := time.Now().Truncate(24*time.Hour)

	result := s.repository.totalMinutesForTaskAndUser(
		timer.TaskHash, timer.TeamUserID, startDate, endDate)

	if timer.FinishedAt == nil {
		duration := time.Since(timer.CreatedAt)
		result += int(duration.Minutes())
	}

	return result
}

// UserTotalMinutesForToday todo
func (s *TimerService) UserTotalMinutesForToday(userID string) int {
	return 0
}
