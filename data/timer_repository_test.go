package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"fmt"
	"github.com/nlopes/slack"
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/pavlo/gosuite"
	"gopkg.in/tylerb/is.v1"
)

func TestTimerRepository(t *testing.T) {
	gosuite.Run(t, &TimerRepositoryTestSuite{Is: is.New(t)})
}

func (s *TimerRepositoryTestSuite) TestUpdate(t *testing.T) {

	project := &models.Project{
		ID:                  bson.NewObjectId(),
		ExternalProjectName: "project",
		ExternalProjectID:   "0987654321",
	}

	user := &models.TeamUser{
		ID:             bson.NewObjectId(),
		ExternalUserID: "user",
		SlackUserInfo: &slack.User{
			TZOffset: 10800,
		},
	}

	timer, err := s.repo.create("team", project, user, "task")
	s.Nil(err)
	s.NotNil(timer)
	s.Equal(timer.Minutes, 0)
	s.Equal(timer.ModelVersion, models.ModelVersionTimer)

	timer.Minutes = 50

	err = s.repo.update(timer)
	s.Nil(err)

	loadedTimer, err := s.repo.findByID(timer.ID.Hex())
	s.Nil(err)
	s.Equal(loadedTimer.Minutes, 50)
}

func (s *TimerRepositoryTestSuite) TestCreateTimer(t *testing.T) {
	project := &models.Project{
		ID:                  bson.NewObjectId(),
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

	timer, err := s.repo.create("team", project, user, "task")
	s.Nil(err)
	s.NotNil(timer)

	timerFromDB, err := s.repo.findByID(timer.ID.Hex())
	s.Nil(err)

	s.NotNil(timerFromDB.CreatedAt) //todo check this
	s.Nil(timerFromDB.DeletedAt) //todo check this
	s.Nil(timerFromDB.FinishedAt) //todo check this
	s.Equal(timerFromDB.Minutes, 0)
	s.Equal(timerFromDB.TeamID, "team")
	s.Equal(timerFromDB.ProjectID, project.ID.Hex())
	s.Equal(timerFromDB.ProjectExternalID, "0987654321")
	s.Equal(timerFromDB.ProjectExternalName, "project")
	s.Equal(timerFromDB.TeamUserID, userID.Hex())
	s.Equal(timerFromDB.TaskName, "task")
	s.Equal(timerFromDB.TeamUserTZOffset, 10800)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserNotExist(t *testing.T) {
	timer, err := s.repo.findActiveByTeamAndUser("does not", "matter")
	s.Nil(err)
	s.Nil(timer)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserExists(t *testing.T) {

	newID := bson.NewObjectId()
	timer := &models.Timer{
		ID:         newID,
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		CreatedAt:  time.Now(),
		TaskName:   "task",
		Minutes:    0,
	}
	s.repo.createTimer(timer)

	timerFromDB, err := s.repo.findActiveByTeamAndUser("team", "user")
	s.Nil(err)
	s.NotNil(timerFromDB)
	s.Equal(timerFromDB.ID.Hex(), newID.Hex())
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserButAlreadyFinished(t *testing.T) {

	newID := bson.NewObjectId()
	finishedAt := time.Now()
	timer := &models.Timer{
		ID:         newID,
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		CreatedAt:  finishedAt,
		FinishedAt: &finishedAt,
		TaskName:   "task",
		Minutes:    0,
	}
	s.repo.createTimer(timer)

	timerFromDB, err := s.repo.findActiveByTeamAndUser("team", "user")
	s.Nil(err)
	s.Nil(timerFromDB)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserButAlreadyDeleted(t *testing.T) {

	newID := bson.NewObjectId()
	deletedAt := time.Now()
	timer := &models.Timer{
		ID:         newID,
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		CreatedAt:  deletedAt,
		DeletedAt:  &deletedAt,
		TaskName:   "task",
		Minutes:    0,
	}
	s.repo.createTimer(timer)

	timerFromDB, err := s.repo.findActiveByTeamAndUser("team", "user")
	s.Nil(err)
	s.Nil(timerFromDB)
}

func (s *TimerRepositoryTestSuite) TestTotalMinutesMethods(t *testing.T) {

	now := time.Now()
	// creates 10 timers one minute each
	for i := 10; i < 20; i++ {
		createdAt := utils.PT(fmt.Sprintf("2016 Sep %d 12:35:00", i))
		s.repo.createTimer(&models.Timer{
			ID:         bson.NewObjectId(),
			TeamID:     "team",
			ProjectID:  "project",
			TeamUserID: "user",
			TaskHash:   "task",
			CreatedAt:  createdAt,
			FinishedAt: &now,
			Minutes:    1,
		})
	}

	// let's add a few more task for different users and tasks
	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "another task",
		CreatedAt:  utils.PT("2016 Sep 12 10:35:00"),
		FinishedAt: &now,
		Minutes:    1,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "another user",
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 13 19:35:00"),
		FinishedAt: &now,
		Minutes:    1,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "another user",
		TaskHash:   "another task",
		CreatedAt:  utils.PT("2016 Sep 14 19:35:00"),
		FinishedAt: &now,
		Minutes:    1,
	})

	// all tasks
	m := s.repo.totalMinutesForTaskAndUser("task", "user", utils.PT("2016 Sep 09 12:35:00"), utils.PT("2016 Sep 21 12:35:00"))
	s.Equal(m, 10)

	// one year later than any of the tasks
	m = s.repo.totalMinutesForTaskAndUser("task", "user", utils.PT("2017 Sep 09 12:35:00"), utils.PT("2017 Sep 21 12:35:00"))
	s.Equal(m, 0)

	// should get one for 10th, one for 11th and one for 12th because the endDate is one minute after the third time
	m = s.repo.totalMinutesForTaskAndUser("task", "user", utils.PT("2016 Sep 10 10:00:00"), utils.PT("2016 Sep 12 12:36:00"))
	s.Equal(m, 3)

	m = s.repo.totalMinutesForUser("user", utils.PT("2016 Sep 09 12:35:00"), utils.PT("2016 Sep 21 12:35:00"))
	s.Equal(m, 11) // 10 regular and one outstanding timer

	m = s.repo.totalMinutesForUser("user", utils.PT("2017 Sep 09 12:35:00"), utils.PT("2017 Sep 21 12:35:00"))
	s.Equal(m, 0)

	m = s.repo.totalMinutesForUser("user", utils.PT("2016 Sep 12 00:00:00"), utils.PT("2016 Sep 12 23:59:59"))
	s.Equal(m, 2)
}

func (s *TimerRepositoryTestSuite) TestCompletedTasksForUser(t *testing.T) {

	now := time.Now()

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task-hash1",
		TaskName:   "task-name1",
		CreatedAt:  utils.PT("2016 Sep 25 12:35:00"),
		FinishedAt: &now,
		Minutes:    5,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task-hash1",
		TaskName:   "task-name1",
		CreatedAt:  utils.PT("2016 Sep 25 12:40:00"),
		FinishedAt: &now,
		Minutes:    10,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task-hash2",
		TaskName:   "task-name2",
		CreatedAt:  utils.PT("2016 Sep 25 12:50:00"),
		FinishedAt: &now,
		Minutes:    20,
	})

	m, err := s.repo.completedTasksForUser("user", utils.PT("2016 Sep 25 12:35:00"), utils.PT("2016 Sep 25 12:45:00"))
	s.Nil(err)

	s.Equal(len(m), 1) // only the `task-hash1` one given the time frame
	s.Equal(m[0].Minutes, 15)
	s.Equal(m[0].Name, "task-name1")

	m, err = s.repo.completedTasksForUser("user", utils.PT("2016 Sep 25 12:35:00"), utils.PT("2016 Sep 25 15:00:00"))
	s.Nil(err)

	s.Equal(len(m), 2)
	s.Equal(m[0].Minutes, 15)
	s.Equal(m[0].Name, "task-name1")
	s.Equal(m[1].Minutes, 20)
	s.Equal(m[1].Name, "task-name2")
}

func (s *TimerRepositoryTestSuite) TestFindActiveByTimezoneOffset(t *testing.T) {
	s.repo.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		FinishedAt:       nil,
		DeletedAt:        nil,
		TeamUserTZOffset: 10,
		TaskHash:         "match",
	})
	s.repo.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		FinishedAt:       nil,
		DeletedAt:        nil,
		TeamUserTZOffset: 10,
		TaskHash:         "match",
	})
	s.repo.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		FinishedAt:       nil,
		DeletedAt:        nil,
		TeamUserTZOffset: 20,
		TaskHash:         "not match",
	})

	now := time.Now()
	s.repo.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		FinishedAt:       &now,
		DeletedAt:        nil,
		TeamUserTZOffset: 10,
		TaskHash:         "not match",
	})

	timers, err := s.repo.findActiveByTimezoneOffset(10)
	s.Nil(err)
	s.Equal(len(timers), 2)

	for _, timer := range timers {
		s.Equal(timer.TaskHash, "match")
	}
}

func (s *TimerRepositoryTestSuite) SetUpSuite() {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.repo = NewTimerRepository(s.session)
}

func (s *TimerRepositoryTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *TimerRepositoryTestSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *TimerRepositoryTestSuite) TearDown() {}

type TimerRepositoryTestSuite struct {
	*is.Is
	env     *utils.Environment
	session *mgo.Session
	repo    *TimerRepository
}
