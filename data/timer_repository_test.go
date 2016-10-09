package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"fmt"
	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (s *TimerRepositoryTestSuite) TestUpdate(c *C) {

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
	c.Assert(err, IsNil)
	c.Assert(timer, NotNil)
	c.Assert(timer.Minutes, Equals, 0)
	c.Assert(timer.ModelVersion, Equals, models.ModelVersionTimer)

	timer.Minutes = 50

	err = s.repo.update(timer)
	c.Assert(err, IsNil)

	loadedTimer, err := s.repo.findByID(timer.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(loadedTimer.Minutes, Equals, 50)
}

func (s *TimerRepositoryTestSuite) TestCreateTimer(c *C) {
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
	c.Assert(err, IsNil)
	c.Assert(timer, NotNil)

	timerFromDB, err := s.repo.findByID(timer.ID.Hex())

	c.Assert(err, IsNil)
	c.Assert(timerFromDB.CreatedAt, NotNil)
	c.Assert(timerFromDB.DeletedAt, IsNil)
	c.Assert(timerFromDB.FinishedAt, IsNil)
	c.Assert(timerFromDB.Minutes, Equals, 0)
	c.Assert(timerFromDB.TeamID, Equals, "team")
	c.Assert(timerFromDB.ProjectID, Equals, project.ID.Hex())
	c.Assert(timerFromDB.ProjectExternalID, Equals, "0987654321")
	c.Assert(timerFromDB.ProjectExternalName, Equals, "project")
	c.Assert(timerFromDB.TeamUserID, Equals, userID.Hex())
	c.Assert(timerFromDB.TaskName, Equals, "task")
	c.Assert(timerFromDB.TeamUserTZOffset, Equals, 10800)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserNotExist(c *C) {
	timer, err := s.repo.findActiveByTeamAndUser("does not", "matter")
	c.Assert(err, IsNil)
	c.Assert(timer, IsNil)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserExists(c *C) {

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
	c.Assert(err, IsNil)
	c.Assert(timerFromDB, NotNil)
	c.Assert(timerFromDB.ID.Hex(), Equals, newID.Hex())
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserButAlreadyFinished(c *C) {

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
	c.Assert(err, IsNil)
	c.Assert(timerFromDB, IsNil)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserButAlreadyDeleted(c *C) {

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
	c.Assert(err, IsNil)
	c.Assert(timerFromDB, IsNil)
}

func (s *TimerRepositoryTestSuite) TestTotalMinutesMethods(c *C) {

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
	c.Assert(m, Equals, 10)

	// one year later than any of the tasks
	m = s.repo.totalMinutesForTaskAndUser("task", "user", utils.PT("2017 Sep 09 12:35:00"), utils.PT("2017 Sep 21 12:35:00"))
	c.Assert(m, Equals, 0)

	// should get one for 10th, one for 11th and one for 12th because the endDate is one minute after the third time
	m = s.repo.totalMinutesForTaskAndUser("task", "user", utils.PT("2016 Sep 10 10:00:00"), utils.PT("2016 Sep 12 12:36:00"))
	c.Assert(m, Equals, 3)

	m = s.repo.totalMinutesForUser("user", utils.PT("2016 Sep 09 12:35:00"), utils.PT("2016 Sep 21 12:35:00"))
	c.Assert(m, Equals, 11) // 10 regular and one outstanding timer

	m = s.repo.totalMinutesForUser("user", utils.PT("2017 Sep 09 12:35:00"), utils.PT("2017 Sep 21 12:35:00"))
	c.Assert(m, Equals, 0)

	m = s.repo.totalMinutesForUser("user", utils.PT("2016 Sep 12 00:00:00"), utils.PT("2016 Sep 12 23:59:59"))
	c.Assert(m, Equals, 2)
}

func (s *TimerRepositoryTestSuite) TestCompletedTasksForUser(c *C) {

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
	c.Assert(err, IsNil)
	c.Assert(len(m), Equals, 1) // only the `task-hash1` one given the time frame
	c.Assert(m[0].Minutes, Equals, 15)
	c.Assert(m[0].Name, Equals, "task-name1")

	m, err = s.repo.completedTasksForUser("user", utils.PT("2016 Sep 25 12:35:00"), utils.PT("2016 Sep 25 15:00:00"))
	c.Assert(err, IsNil)
	c.Assert(len(m), Equals, 2)
	c.Assert(m[0].Minutes, Equals, 15)
	c.Assert(m[0].Name, Equals, "task-name1")

	c.Assert(m[1].Minutes, Equals, 20)
	c.Assert(m[1].Name, Equals, "task-name2")
}

func (s *TimerRepositoryTestSuite) TestFindActiveByTimezoneOffset(c *C) {
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
	c.Assert(err, IsNil)
	c.Assert(len(timers), Equals, 2)

	for _, timer := range timers {
		c.Assert(timer.TaskHash, Equals, "match")
	}
}

// Suite lifecycle and callbacks
func (s *TimerRepositoryTestSuite) SetUpSuite(c *C) {
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

func (s *TimerRepositoryTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *TimerRepositoryTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestTimerRepository(t *testing.T) { TestingT(t) }

type TimerRepositoryTestSuite struct {
	env     *utils.Environment
	session *mgo.Session
	repo    *TimerRepository
}

var _ = Suite(&TimerRepositoryTestSuite{})
