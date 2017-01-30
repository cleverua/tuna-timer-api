package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/nlopes/slack"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/tylerb/is.v1"
	"github.com/pavlo/gosuite"
)

func TestTimerService(t *testing.T) {
	gosuite.Run(t, &TimerServiceTestSuite{Is: is.New(t)})
}

func (s *TimerServiceTestSuite) TestgetActiveTimer(t *testing.T) {

	now := time.Now()

	// completed
	s.repo.CreateTimer(&models.Timer{
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
	s.repo.CreateTimer(&models.Timer{
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

func (s *TimerServiceTestSuite) TeststopTimer(t *testing.T) {
	now := time.Now()

	offsetDuration, _ := time.ParseDuration("20m")
	timerStartedAt := now.Add(offsetDuration * -1) // 20 minutes ago

	id := bson.NewObjectId()
	timer, err := s.repo.CreateTimer(&models.Timer{
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
	s.Equal(loadedTimer.ActualMinutes, 20)
	s.NotNil(loadedTimer.FinishedAt)
}

func (s *TimerServiceTestSuite) TeststartTimer(t *testing.T) {

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

func (s *TimerServiceTestSuite) TesttotalMinutesForTodayAddsTimeForUnfinishedTask(t *testing.T) {
	now := time.Now()

	offsetDuration1, _ := time.ParseDuration("20m")
	firstTimerStartedAt := now.Add(offsetDuration1 * -1) // 20 minutes ago

	offsetDuration2, _ := time.ParseDuration("5m")
	secondTimerStartedAt := now.Add(offsetDuration2 * -1) // 5 minutes ago

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  now.Add(offsetDuration1 * -1),
		FinishedAt: &firstTimerStartedAt,
		Minutes:    10,
	})

	timer, _ := s.repo.CreateTimer(&models.Timer{
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

func (s *TimerServiceTestSuite) TesttotalCompletedMinutesForDay(t *testing.T) {
	now := time.Now()

	user := &models.TeamUser{
		ID: bson.NewObjectId(),
		SlackUserInfo: &slack.User{
			TZOffset: 10800, // UTC+3 Kiev
		},
	}

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 8:00:00"), // which is 11:00 in Kiev
		FinishedAt: &now,
		Minutes:    2,
	})

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 19:30:00"), // which is 22:30 in Kiev
		FinishedAt: &now,
		Minutes:    3,
	})

	s.repo.CreateTimer(&models.Timer{
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

func (s *TimerServiceTestSuite) TestgetCompletedTasksForDayPositiveTZOffset(t *testing.T) {

	now := time.Now()

	user := &models.TeamUser{
		ID: bson.NewObjectId(),
		SlackUserInfo: &slack.User{
			TZOffset: 10800, // UTC+3 Kiev
		},
	}

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 8:00:00"), // which is 11:00 in Kiev
		FinishedAt: &now,
		Minutes:    2,
	})

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 19:30:00"), // which is 22:30 in Kiev
		FinishedAt: &now,
		Minutes:    3,
	})

	s.repo.CreateTimer(&models.Timer{
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

func (s *TimerServiceTestSuite) TestgetCompletedTasksForDayNegativeTZOffset(t *testing.T) {

	now := time.Now()

	user := &models.TeamUser{
		ID: bson.NewObjectId(),
		SlackUserInfo: &slack.User{
			TZOffset: -18000, // UTC-5 Nashville
		},
	}

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 15:00:00"), // which is 10:00 in Nashville
		FinishedAt: &now,
		Minutes:    2,
	})

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 13 03:30:00"), // which is 22:30 in Nashville
		FinishedAt: &now,
		Minutes:    3,
	})

	s.repo.CreateTimer(&models.Timer{
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
func (s *TimerServiceTestSuite) TestcompleteActiveTimersAtMidnight(t *testing.T) {

	t1ID := bson.NewObjectId()
	s.repo.CreateTimer(&models.Timer{
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
	s.repo.CreateTimer(&models.Timer{
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

func (s *TimerServiceTestSuite) TestGetUserTasksByRange(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		CreatedAt:  utils.PT("2016 Dec 20 10:35:00"),
		Minutes:    20,
	})
	s.repo.CreateTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: user.ID.Hex(),
		CreatedAt:  utils.PT("2016 Dec 21 10:35:00"),
		Minutes:    20,
	})

	timers, err := s.service.GetUserTimersByRange("2016-12-20", "2016-12-21", user)
	s.Nil(err)
	s.Len(timers, 2)
	for _, timer := range timers {
		s.Equal(timer.TeamUserID, user.ID.Hex())
	}
}

func (s *TimerServiceTestSuite) TestGetUserTasksByRangeFail(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	//Should return an error, with wrong startDate format
	startDate := "2016/12/20"
	endDate := "2016-12-21"

	timers, err := s.service.GetUserTimersByRange(startDate, endDate, user)
	s.Err(err)
	s.Len(timers, 0)

	//Should return an error, with wrong endDate format
	startDate = "2016-12-20"
	endDate = "2016:12:21"

	timers, err = s.service.GetUserTimersByRange(startDate, endDate, user)
	s.Err(err)
	s.Len(timers, 0)

	//Should return an error, with range more than 31 day
	startDate = "2016-11-30"
	endDate = "2016-12-31"

	timers, err = s.service.GetUserTimersByRange(startDate, endDate, user)
	s.Err(err)
	s.Equal(err.Error(), "Too much days in range")
	s.Len(timers, 0)
}

func (s *TimerServiceTestSuite) TestUpdateUserTimer(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		TeamID:		"team",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	timer := &models.Timer{
		ID:			bson.NewObjectId(),
		TaskName:		"task-name",
		TeamID:			"team",
		ProjectID:		"project-id",
		ProjectExternalID:	"project-external-id",
		ProjectExternalName:	"project-external-name",
		TeamUserID:		user.ID.Hex(),
		CreatedAt:		utils.PT("2016 Dec 20 10:35:00"),
		Minutes:		20,
		ActualMinutes:		20,
	}
	s.repo.CreateTimer(timer)

	newTimerData := &models.Timer{
		ID:			bson.NewObjectId(),
		TaskName:		"new-task-name",
		TeamID:			"new-team-id",
		ProjectID:		"new-project-id",
		ProjectExternalID:	"new-project-external-id",
		ProjectExternalName:	"new-project-external-name",
		TeamUserID:		bson.NewObjectId().Hex(),
		CreatedAt:		utils.PT("2016 Dec 20 12:50:00"),
		Minutes:		40,
		ActualMinutes:		40,
		Edits:	[]*models.TimeEdit{
			{TeamUserID: user.ID.Hex(), CreatedAt: time.Now(), Minutes: 10},
		},
	}

	err := s.service.UpdateUserTimer(user, timer, newTimerData)
	s.Nil(err)

	// Check permit parameters: Edits, TaskName, ProjectID, ProjectExternalID, ProjectExternalName
	s.Equal(timer.Edits, newTimerData.Edits)
	s.Equal(timer.TaskName, newTimerData.TaskName)
	s.Equal(timer.ProjectID, newTimerData.ProjectID)
	s.Equal(timer.ProjectExternalID, newTimerData.ProjectExternalID)
	s.Equal(timer.ProjectExternalName, newTimerData.ProjectExternalName)
	// Check calculated params
	s.Equal(timer.Minutes, 30)
	s.Equal(timer.ActualMinutes, 20)
	// Check for other parameters didn't change
	s.NotEqual(timer.ID, newTimerData.ID)
	s.NotEqual(timer.TeamID, newTimerData.TeamID)
	s.NotEqual(timer.TeamUserID, newTimerData.TeamUserID)
	s.NotEqual(timer.Minutes, newTimerData.Minutes)
	s.NotEqual(timer.Minutes, newTimerData.Minutes)
	s.NotEqual(timer.ActualMinutes, newTimerData.ActualMinutes)
	s.NotEqual(timer.CreatedAt, newTimerData.CreatedAt)
}

func (s *TimerServiceTestSuite) TestUpdateUserTimerWithWrongUserID(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	timer := &models.Timer{
		TaskName:		"task-name",
		ProjectID:		"project-id",
		ProjectExternalID:	"project-external-id",
		ProjectExternalName:	"project-external-name",
		TeamUserID:		 bson.NewObjectId().Hex(),
		TeamID:			"team",
	}
	s.repo.CreateTimer(timer)

	newTimerData := &models.Timer{
		TaskName:		"new-task-name",
		ProjectID:		"new-project-id",
		ProjectExternalID:	"new-project-external-id",
		ProjectExternalName:	"new-project-external-name",
	}

	err := s.service.UpdateUserTimer(user, timer, newTimerData)
	// Should return error, and didn't change timer
	s.NotNil(err)
	s.Equal(err.Error(), "update forbidden")
	s.NotEqual(timer.TaskName, newTimerData.TaskName)
	s.NotEqual(timer.ProjectID, newTimerData.ProjectID)
	s.NotEqual(timer.ProjectExternalID, newTimerData.ProjectExternalID)
	s.NotEqual(timer.ProjectExternalName, newTimerData.ProjectExternalName)
}

func (s *TimerServiceTestSuite) TestUpdateUserTimerForSlackOwner(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		TeamID:		"team",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
			IsOwner:  true,
		},
	}

	timer := &models.Timer{
		ID:			bson.NewObjectId(),
		TaskName:		"task-name",
		TeamID:			"team",
		ProjectID:		"project-id",
		ProjectExternalID:	"project-external-id",
		ProjectExternalName:	"project-external-name",
		TeamUserID:		bson.NewObjectId().Hex(),
		CreatedAt:		utils.PT("2016 Dec 20 10:35:00"),
		Minutes:		20,
		ActualMinutes:		20,
	}
	s.repo.CreateTimer(timer)

	newTimerData := &models.Timer{
		ID:			bson.NewObjectId(),
		TaskName:		"new-task-name",
		TeamID:			"new-team-id",
		ProjectID:		"new-project-id",
		ProjectExternalID:	"new-project-external-id",
		ProjectExternalName:	"new-project-external-name",
		TeamUserID:		bson.NewObjectId().Hex(),
		CreatedAt:		utils.PT("2016 Dec 20 12:50:00"),
		Minutes:		40,
		ActualMinutes:		40,
		Edits:	[]*models.TimeEdit{
			{TeamUserID: user.ID.Hex(), CreatedAt: time.Now(), Minutes: 10},
		},
	}

	err := s.service.UpdateUserTimer(user, timer, newTimerData)
	s.Nil(err)
	s.Equal(timer.TaskName, newTimerData.TaskName)
	s.Equal(timer.ProjectID, newTimerData.ProjectID)
	s.Equal(timer.ProjectExternalID, newTimerData.ProjectExternalID)
	s.Equal(timer.ProjectExternalName, newTimerData.ProjectExternalName)
}

func (s *TimerServiceTestSuite) TestDeleteUserTimer(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		TeamID:		"team",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	timerMinutes := 10
	timer := &models.Timer{
		ID:			bson.NewObjectId(),
		TaskName:		"task-name",
		TeamID:			"team",
		ProjectID:		"project-id",
		ProjectExternalID:	"project-external-id",
		ProjectExternalName:	"project-external-name",
		TeamUserID:		user.ID.Hex(),
		CreatedAt:		time.Now().Add(time.Duration(-timerMinutes) * time.Minute),
		Minutes:		timerMinutes,
	}
	s.repo.CreateTimer(timer)

	err := s.service.DeleteUserTimer(user, timer)
	s.Nil(err)

	s.NotNil(timer.DeletedAt)
}

func (s *TimerServiceTestSuite) TestDeleteUserActiveTimer(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		TeamID:		"team",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	timerMinutes := 10
	timer := &models.Timer{
		ID:			bson.NewObjectId(),
		TaskName:		"task-name",
		TeamID:			"team",
		ProjectID:		"project-id",
		ProjectExternalID:	"project-external-id",
		ProjectExternalName:	"project-external-name",
		TeamUserID:		user.ID.Hex(),
		CreatedAt:		time.Now().Add(time.Duration(-timerMinutes) * time.Minute),
	}
	s.repo.CreateTimer(timer)

	err := s.service.DeleteUserTimer(user, timer)
	s.Nil(err)

	s.NotNil(timer.DeletedAt)
	s.Equal(timer.Minutes, timerMinutes)
	s.Equal(timer.ActualMinutes, timerMinutes)
}

func (s *TimerServiceTestSuite) TestDeleteUserTimerWithWrongUserID(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		TeamID:		"team",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	timerMinutes := 10
	timer := &models.Timer{
		ID:			bson.NewObjectId(),
		TaskName:		"task-name",
		TeamID:			"team",
		ProjectID:		"project-id",
		ProjectExternalID:	"project-external-id",
		ProjectExternalName:	"project-external-name",
		TeamUserID:		bson.NewObjectId().Hex(),
		CreatedAt:		time.Now().Add(time.Duration(-timerMinutes) * time.Minute),
	}
	s.repo.CreateTimer(timer)

	err := s.service.DeleteUserTimer(user, timer)
	s.NotNil(err)
	s.Equal(err.Error(), "delete forbidden")
}

func (s *TimerServiceTestSuite) TestUserMonthStatistics(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		TeamID:		"team",
		SlackUserInfo: &slack.User{},
	}
	startDate := utils.PT("2016 Dec 01 00:00:00")
	minutes := 10
	days := 31
	date := "2016-12-1"

	for i := 0; i < days + 1; i++ {
		m := minutes + i
		finished := startDate.AddDate(0, 0, i).Add(time.Minute * time.Duration(m))
		s.repo.CreateTimer(&models.Timer{
			ID:			bson.NewObjectId(),
			TeamID:			"team",
			ProjectID:		"project",
			ProjectExternalName:	"project_name",
			TeamUserID:		user.ID.Hex(),
			CreatedAt:		startDate.AddDate(0, 0, i),
			FinishedAt:		&finished,
			TeamUserTZOffset:	user.SlackUserInfo.TZOffset,
			Minutes:		m,
			ActualMinutes:		m,
		})
	}

	result, err := s.service.UserMonthStatistics(user, date)
	s.Nil(err)
	s.Len(result, days)
	s.Equal(result[0].Minutes, 10)
	s.Equal(result[days - 1].Minutes, 40)
}

func (s *TimerServiceTestSuite) TestUserMonthStatisticsWithWrongDate(t *testing.T) {
	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		TeamID:		"team",
		SlackUserInfo: &slack.User{},
	}
	dates := [...]string{"", "2016-12-", "2016/12/1", "not a date"}

	for _, date := range dates {
		result, err := s.service.UserMonthStatistics(user, date)
		s.Nil(result)
		s.NotNil(err)
	}
}

func (s *TimerServiceTestSuite) SetUpSuite() {
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
