package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

// TimerService - the structure of the service
type TimerService struct {
	repository *TimerRepository
}

// NewTimerService constructs an instance of the service
func NewTimerService(session *mgo.Session) *TimerService {
	return &TimerService{
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
func (s *TimerService) StartTimer(teamID string, project *models.Project, user *models.TeamUser, taskName string) (*models.Timer, error) {
	return s.repository.create(teamID, project, user, taskName)
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
func (s *TimerService) TotalCompletedMinutesForDay(year int, month time.Month, day int, user *models.TeamUser) int {

	log.Printf("TotalUserMinutesForDay, Year: %d, Month: %d, Day: %d", year, month, day)

	tzOffset := user.SlackUserInfo.TZOffset
	log.Printf("TotalUserMinutesForDay, tzOffset: %d", tzOffset)

	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(time.Duration(tzOffset) * time.Second * -1)
	endDate := time.Date(year, month, day, 23, 59, 59, 0, time.UTC).Add(time.Duration(tzOffset) * time.Second * -1)

	return s.repository.totalMinutesForUser(user.ID.Hex(), startDate, endDate)
}

// GetCompletedTasksForDay - returns the list of tasks the user had completed during given work day by his/her timezone
// - year, month, day - is the day to get the list of completed tasks for
// - user - whose tasks the viewer is interested in
func (s *TimerService) GetCompletedTasksForDay(year int, month time.Month, day int, user *models.TeamUser) ([]*models.TaskAggregation, error) {
	log.Printf("GetCompletedTasksForDay, Year: %d, Month: %d, Day: %d", year, month, day)

	tzOffset := user.SlackUserInfo.TZOffset
	log.Printf("GetCompletedTasksForDay, tzOffset: %d", tzOffset)

	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(time.Duration(tzOffset) * time.Second * -1)
	endDate := time.Date(year, month, day, 23, 59, 59, 0, time.UTC).Add(time.Duration(tzOffset) * time.Second * -1)

	log.Printf("GetCompletedTasksForDay, startDate: %+v", startDate)
	log.Printf("GetCompletedTasksForDay, endDate: %+v", endDate)

	tasks, err := s.repository.completedTasksForUser(user.ID.Hex(), startDate, endDate)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TimerService) CompleteActiveTimersAtMidnight(utcNow *time.Time) error {
	timezoneOffset := utils.WhichTimezoneIsMidnightAt(utcNow.Hour(), utcNow.Minute())
	timers, err := s.repository.findActiveByTimezoneOffset(timezoneOffset)
	if err != nil {
		return err
	}

	log.Printf("Found %d timer(s) to complete", len(timers))

	for _, timer := range timers {
		log.Printf("Completing %s timer", timer.TaskName)

		endDate := time.Date(timer.CreatedAt.Year(), timer.CreatedAt.Month(), timer.CreatedAt.Day(), utcNow.Hour()-1, 59, 59, 0, time.UTC)
		timer.FinishedAt = &endDate

		timer.Minutes = int(endDate.Sub(timer.CreatedAt).Minutes())
		err = s.repository.update(timer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TimerService) CalculateMinutesForActiveTimer(timer *models.Timer) int {
	duration := time.Since(timer.CreatedAt)
	return int(duration.Minutes())
}
