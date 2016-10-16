package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/tylerb/is.v1"
	"github.com/pavlo/gosuite"
)

func TestTimerService(t *testing.T) {
	gosuite.Run(t, &TimerServiceTestSuite{})
}

func (s *TimerServiceTestSuite) GSTgetActiveTimer(t *testing.T) {

	now := time.Now()

	// completed
	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  now,
		FinishedAt: &now,
		Minutes:    10,
	})

	// not completed
	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  now,
		Minutes:    20,
	})

	timer, err := s.service.GetActiveTimer("team", "user")
	s.Nil(err)
	s.NotNil(timer)

	s.Equal(timer.Minutes, 20)
}

func (s *TimerServiceTestSuite) GSTstopTimer(t *testing.T) {
	now := time.Now()

	offsetDuration, _ := time.ParseDuration("20m")
	timerStartedAt := now.Add(offsetDuration * -1) // 20 minutes ago

	id := bson.NewObjectId()
	timer, err := s.repo.createTimer(&models.Timer{
		ID:         id,
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  timerStartedAt,
		Minutes:    0,
	})

	s.Nil(err)
	s.NotNil(timer)

	s.service.StopTimer(timer)

	loadedTimer, err := s.repo.findByID(id.Hex())
	s.Nil(err)

	s.Equal(loadedTimer.Minutes, 20)
	s.NotNil(loadedTimer.FinishedAt)
}

func (s *TimerServiceTestSuite) GSTstartTimer(t *testing.T) {

	projectID := bson.NewObjectId()
	project := &models.Project{
		ID:                  projectID,
		ExternalProjectName: "project",
		ExternalProjectID:   "0987654321",
	}

	userID := bson.NewObjectId()
	user := &models.TeamUser{
		ID:             userID,
		ExternalUserID: "user",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	timer, err := s.service.StartTimer("team", project, user, "task")
	s.Nil(err)
	s.NotNil(timer)

	loadedTimer, err := s.repo.findByID(timer.ID.Hex())
	s.Nil(err)
	s.NotNil(loadedTimer)

	s.Equal(loadedTimer.TeamID, "team")
	s.Equal(loadedTimer.ProjectID, projectID.Hex())
	s.Equal(loadedTimer.TeamUserID, userID.Hex())
	s.Equal(loadedTimer.TaskName, "task")
	s.NotNil(loadedTimer.TaskHash)
	s.NotNil(loadedTimer.CreatedAt)
	s.Nil(loadedTimer.FinishedAt)
	s.Nil(loadedTimer.DeletedAt)
	s.Equal(loadedTimer.Minutes, 0)
}

func (s *TimerServiceTestSuite) GSTtotalMinutesForTodayAddsTimeForUnfinishedTask(t *testing.T) {
	now := time.Now()

	offsetDuration1, _ := time.ParseDuration("20m")
	firstTimerStartedAt := now.Add(offsetDuration1 * -1) // 20 minutes ago

	offsetDuration2, _ := time.ParseDuration("5m")
	secondTimerStartedAt := now.Add(offsetDuration2 * -1) // 5 minutes ago

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  now.Add(offsetDuration1 * -1),
		FinishedAt: &firstTimerStartedAt,
		Minutes:    10,
	})

	timer, _ := s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  secondTimerStartedAt,
		FinishedAt: nil,
		Minutes:    0,
	})

	//c.Assert(s.service.TotalMinutesForTaskToday(timer), Equals, 15)
	s.Equal(s.service.TotalMinutesForTaskToday(timer), 15)
}

func (s *TimerServiceTestSuite) GSTtotalCompletedMinutesForDay(t *testing.T) {
	now := time.Now()

	user := &models.TeamUser{
		ID: bson.NewObjectId(),
		SlackUserInfo: &slack.User{
			TZOffset: 10800, // UTC+3 Kiev
		},
	}

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 8:00:00"), // which is 11:00 in Kiev
		FinishedAt: &now,
		Minutes:    2,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 19:30:00"), // which is 22:30 in Kiev
		FinishedAt: &now,
		Minutes:    3,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 22:00:00"), // which is 1am of the next day in Kiev
		FinishedAt: &now,
		Minutes:    7,
	})

	targetDate := utils.PT("2016 Sep 12 00:00:00")
	s.Equal(s.service.TotalCompletedMinutesForDay(targetDate.Year(), targetDate.Month(), targetDate.Day(), user), 5)

	targetDate = utils.PT("2016 Sep 13 00:00:00")
	s.Equal(s.service.TotalCompletedMinutesForDay(targetDate.Year(), targetDate.Month(), targetDate.Day(), user), 7)
}

func (s *TimerServiceTestSuite) GSTgetCompletedTasksForDayPositiveTZOffset(t *testing.T) {

	now := time.Now()

	user := &models.TeamUser{
		ID: bson.NewObjectId(),
		SlackUserInfo: &slack.User{
			TZOffset: 10800, // UTC+3 Kiev
		},
	}

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 8:00:00"), // which is 11:00 in Kiev
		FinishedAt: &now,
		Minutes:    2,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 19:30:00"), // which is 22:30 in Kiev
		FinishedAt: &now,
		Minutes:    3,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 22:00:00"), // which is 1am of the next day in Kiev
		FinishedAt: &now,
		Minutes:    7,
	})

	targetDate := utils.PT("2016 Sep 12 00:00:00")
	v, err := s.service.GetCompletedTasksForDay(targetDate.Year(), targetDate.Month(), targetDate.Day(), user)
	s.Nil(err)
	s.Equal(len(v), 1)
	s.Equal(v[0].Minutes, 5)

	targetDate = utils.PT("2016 Sep 13 00:00:00")
	v, err = s.service.GetCompletedTasksForDay(targetDate.Year(), targetDate.Month(), targetDate.Day(), user)
	s.Nil(err)
	s.Equal(len(v), 1)
	s.Equal(v[0].Minutes, 7)
}

func (s *TimerServiceTestSuite) GSTgetCompletedTasksForDayNegativeTZOffset(t *testing.T) {

	now := time.Now()

	user := &models.TeamUser{
		ID: bson.NewObjectId(),
		SlackUserInfo: &slack.User{
			TZOffset: -18000, // UTC-5 Nashville
		},
	}

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 15:00:00"), // which is 10:00 in Nashville
		FinishedAt: &now,
		Minutes:    2,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 13 03:30:00"), // which is 22:30 in Nashville
		FinishedAt: &now,
		Minutes:    3,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 13 06:00:00"), // which is 1am of the next day in Nashville
		FinishedAt: &now,
		Minutes:    7,
	})

	targetDate := utils.PT("2016 Sep 12 00:00:00")
	v, err := s.service.GetCompletedTasksForDay(targetDate.Year(), targetDate.Month(), targetDate.Day(), user)
	s.Nil(err)
	s.Equal(len(v), 1)
	s.Equal(v[0].Minutes, 5)

	targetDate = utils.PT("2016 Sep 13 00:00:00")

	v, err = s.service.GetCompletedTasksForDay(targetDate.Year(), targetDate.Month(), targetDate.Day(), user)
	s.Nil(err)
	s.Equal(len(v), 1)
	s.Equal(v[0].Minutes, 7)
}

// CompleteActiveTimersAtMidnight
func (s *TimerServiceTestSuite) GSTcompleteActiveTimersAtMidnight(t *testing.T) {

	t1ID := bson.NewObjectId()
	s.repo.createTimer(&models.Timer{
		ID:               t1ID,
		TeamID:           "team",
		ProjectID:        "project",
		TeamUserID:       "user",
		TaskHash:         "task1",
		CreatedAt:        utils.PT("2016 Sep 12 20:40:00"),
		FinishedAt:       nil,
		TeamUserTZOffset: 10800, // +3 Kiev
	})

	t2ID := bson.NewObjectId()
	s.repo.createTimer(&models.Timer{
		ID:               t2ID,
		TeamID:           "team",
		ProjectID:        "project",
		TeamUserID:       "user",
		TaskHash:         "task2",
		CreatedAt:        utils.PT("2016 Sep 12 10:35:00"),
		FinishedAt:       nil,
		TeamUserTZOffset: -10800, // Rio
	})

	// let now be 21:00 UTC which is midnight in Kiev (+3 UTC)
	now := utils.PT("2016 Sep 12 21:00:00")

	// a side check just to make sure we're on the right track
	//c.Assert(utils.WhichTimezoneIsMidnightAt(now.Hour(), now.Minute()), Equals, 10800)
	s.Equal(utils.WhichTimezoneIsMidnightAt(now.Hour(), now.Minute()), 10800)

	err := s.service.CompleteActiveTimersAtMidnight(&now)
	s.Nil(err)

	timer, err := s.repo.findByID(t1ID.Hex())
	s.Nil(err)
	s.NotNil(timer.FinishedAt)
	s.Equal(timer.Minutes, 20-1)

	timer, err = s.repo.findByID(t2ID.Hex())
	s.Nil(err)
	s.Nil(timer.FinishedAt)
	s.Equal(timer.Minutes, 0)
}

func (s *TimerServiceTestSuite) SetUpSuite(t *testing.T) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.service = NewTimerService(s.session)
	s.repo = NewTimerRepository(s.session)
	s.Is = is.New(t)
}

func (s *TimerServiceTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *TimerServiceTestSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *TimerServiceTestSuite) TearDown() {
}

type TimerServiceTestSuite struct {
	*is.Is
	env     *utils.Environment
	session *mgo.Session
	repo    *TimerRepository
	service *TimerService
}
