package data

import (
	"github.com/pavlo/slack-time/models"
	"gopkg.in/mgo.v2"
	"time"
)

// TimerService - the structure of the service
type TimerService struct {
	session    *mgo.Session
	repository *TimerRepository
}

// NewTimerService constructs an instance of the service
func NewTimerService(session *mgo.Session) *TimerService {
	return &TimerService{
		session:    session,
		repository: NewTimerRepository(session),
	}
}

// GetActiveTimer returns a timer the user is currently working on
func (s *TimerService) GetActiveTimer(teamID, userID string) (*models.Timer, error) {
	timer, err := s.repository.findActiveByTeamAndUser(teamID, userID)
	return timer, err
}

// StopTimer stops the timer and updates its Minutes field
func (s *TimerService) StopTimer(timer *models.Timer) error {
	now := time.Now()
	timer.Minutes = s.calculateMinutesForTimer(timer)
	timer.FinishedAt = &now
	return s.repository.update(timer)
}

// StartTimer creates a new timer
func (s *TimerService) StartTimer(teamID, projectID, teamUserID, taskName string) (*models.Timer, error) {
	return s.repository.create(teamID, projectID, teamUserID, taskName)
}

// TotalMinutesForTaskToday calculates the total number of minutes the user was/is working on particular task today
func (s *TimerService) TotalMinutesForTaskToday(timer *models.Timer) int {
	endDate := time.Now()
	startDate := time.Now().Truncate(24 * time.Hour)

	result := s.repository.totalMinutesForTaskAndUser(
		timer.TaskHash, timer.TeamUserID, startDate, endDate)

	if timer.FinishedAt == nil {
		result += s.calculateMinutesForTimer(timer)
	}

	return result
}

// UserTotalMinutesForToday todo
func (s *TimerService) UserTotalMinutesForToday(userID string) int {
	return 0
}

func (s *TimerService) calculateMinutesForTimer(timer *models.Timer) int {
	duration := time.Since(timer.CreatedAt)
	return int(duration.Minutes())
}
