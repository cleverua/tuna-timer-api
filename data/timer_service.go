package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
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
	timer.Minutes = s.CalculateMinutesForActiveTimer(timer)
	timer.FinishedAt = &now
	return s.repository.update(timer)
}

// StartTimer creates a new timer
func (s *TimerService) StartTimer(teamID string, project *models.Project, teamUserID string, taskName string) (*models.Timer, error) {
	return s.repository.create(teamID, project, teamUserID, taskName)
}

// TotalMinutesForTaskToday calculates the total number of minutes the user was/is working on particular task today
func (s *TimerService) TotalMinutesForTaskToday(timer *models.Timer) int {
	endDate := time.Now()
	startDate := time.Now().Truncate(24 * time.Hour)

	result := s.repository.totalMinutesForTaskAndUser(
		timer.TaskHash, timer.TeamUserID, startDate, endDate)

	if timer.FinishedAt == nil {
		result += s.CalculateMinutesForActiveTimer(timer)
	}

	return result
}

// UserTotalMinutesForToday calculates the total number of minute this user contributed to any project today
func (s *TimerService) TotalUserMinutesForDay(userID string, day time.Time) int {
	startDate := day.Truncate(24 * time.Hour)

	result := s.repository.totalMinutesForUser(userID, startDate, day)

	activeTimer, _ := s.repository.findActiveByUser(userID)
	if activeTimer != nil {
		if activeTimer.CreatedAt.Unix() <= startDate.Unix() {
			activeTimer.CreatedAt = startDate
		}

		result += s.CalculateMinutesForActiveTimer(activeTimer)
	}

	return result
}

// GetCompletedTasksForDay
func (s *TimerService) GetCompletedTasksForDay(userID string, day time.Time) ([]*models.TaskAggregation, error) {

	startDate := day.Truncate(24 * time.Hour)
	endDate := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, time.UTC)

	tasks, err := s.repository.completedTasksForUser(userID, startDate, endDate)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TimerService) CalculateMinutesForActiveTimer(timer *models.Timer) int {
	duration := time.Since(timer.CreatedAt)
	return int(duration.Minutes())
}
