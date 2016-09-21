package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/pavlo/slack-time/utils"
	. "gopkg.in/check.v1"
	"github.com/pavlo/slack-time/models"
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
)

const timerRepositoryTestSuiteTimeParseLayout = "2006 Jan 02 15:04:05"

func (s *TimerRepositoryTestSuite) TestUpdate(c *C) {

	timer, err := s.repository.create("teamID", "projectID", "userID", "taskName")
	c.Assert(err, IsNil)
	c.Assert(timer, NotNil)
	c.Assert(timer.Minutes, Equals, 0)

	timer.Minutes = 50

	err = s.repository.update(timer)
	c.Assert(err, IsNil)

	loadedTimer, err := s.repository.findByID(timer.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(loadedTimer.Minutes, Equals, 50)
}

func (s *TimerRepositoryTestSuite) TestCreateTimer(c *C) {
	timer, err := s.repository.create("teamID", "projectID", "userID", "taskName")
	c.Assert(err, IsNil)
	c.Assert(timer, NotNil)

	timerFromDB, err := s.repository.findByID(timer.ID.Hex())

	c.Assert(err, IsNil)
	c.Assert(timerFromDB.CreatedAt, NotNil)
	c.Assert(timerFromDB.DeletedAt, IsNil)
	c.Assert(timerFromDB.FinishedAt, IsNil)
	c.Assert(timerFromDB.Minutes, Equals, 0)
	c.Assert(timerFromDB.TeamID, Equals, "teamID")
	c.Assert(timerFromDB.ProjectID, Equals, "projectID")
	c.Assert(timerFromDB.TeamUserID, Equals, "userID")
	c.Assert(timerFromDB.TaskName, Equals, "taskName")
	// it is sha256 of "teamIDprojectIDtaskName" is bac825e73328d01c534aacf90e16db89d102bb7bda818ee0be1f42f59f043025
	c.Assert(timerFromDB.TaskHash, Equals, "bac825")
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserNotExist(c *C) {
	timer, err := s.repository.findActiveByTeamAndUser("does not", "matter")
	c.Assert(err, IsNil)
	c.Assert(timer, IsNil)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserExists(c *C) {

	newID := bson.NewObjectId()
	timer := &models.Timer{
		ID:               newID,
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "userID",
		CreatedAt:        time.Now(),
		TaskName:   	  "taskName",
		Minutes:  		  0,
	}
	s.repository.createTimer(timer)


	timerFromDB, err := s.repository.findActiveByTeamAndUser("teamID", "userID")
	c.Assert(err, IsNil)
	c.Assert(timerFromDB, NotNil)
	c.Assert(timerFromDB.ID.Hex(), Equals, newID.Hex())
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserButAlreadyFinished(c *C) {

	newID := bson.NewObjectId()
	finishedAt := time.Now()
	timer := &models.Timer{
		ID:               newID,
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "userID",
		CreatedAt:        finishedAt,
		FinishedAt:		  &finishedAt,
		TaskName:   	  "taskName",
		Minutes:  		  0,
	}
	s.repository.createTimer(timer)

	timerFromDB, err := s.repository.findActiveByTeamAndUser("teamID", "userID")
	c.Assert(err, IsNil)
	c.Assert(timerFromDB, IsNil)
}

func (s *TimerRepositoryTestSuite) TestFindActiveTimerByTeamAndUserButAlreadyDeleted(c *C) {

	newID := bson.NewObjectId()
	deletedAt := time.Now()
	timer := &models.Timer{
		ID:               newID,
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "userID",
		CreatedAt:        deletedAt,
		DeletedAt:		  &deletedAt,
		TaskName:   	  "taskName",
		Minutes:  		  0,
	}
	s.repository.createTimer(timer)

	timerFromDB, err := s.repository.findActiveByTeamAndUser("teamID", "userID")
	c.Assert(err, IsNil)
	c.Assert(timerFromDB, IsNil)
}


func (s *TimerRepositoryTestSuite) TestTotalMinutesForTaskAndUser(c *C) {

	now := time.Now()
	// creates 10 timers one minute each
	for i := 10; i < 20; i++ {
		createdAt := s.pt(fmt.Sprintf("2016 Sep %d 12:35:00", i))
		s.repository.createTimer(&models.Timer{
			ID:               bson.NewObjectId(),
			TeamID:           "teamID",
			ProjectID:        "projectID",
			TeamUserID:       "u",
			TaskHash:   	  "t",
			CreatedAt:        createdAt,
			FinishedAt:		  &now,
			Minutes:  		  1,
		})
	}

	// let's add a few more task for different users and tasks
	s.repository.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "u",
		TaskHash:   	  "another task",
		CreatedAt:        s.pt("2016 Sep 12 10:35:00"),
		FinishedAt:		  &now,
		Minutes:  		  1,
	})

	s.repository.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "another user",
		TaskHash:   	  "t",
		CreatedAt:        s.pt("2016 Sep 13 19:35:00"),
		FinishedAt:		  &now,
		Minutes:  		  1,
	})

	s.repository.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "another user",
		TaskHash:   	  "another task",
		CreatedAt:        s.pt("2016 Sep 14 19:35:00"),
		FinishedAt:		  &now,
		Minutes:  		  1,
	})

	// all tasks
	m := s.repository.totalMinutesForTaskAndUser("t", "u", s.pt("2016 Sep 09 12:35:00"), s.pt("2016 Sep 21 12:35:00"))
	c.Assert(m, Equals, 10)

	// one year later than any of the tasks
	m = s.repository.totalMinutesForTaskAndUser("t", "u", s.pt("2017 Sep 09 12:35:00"), s.pt("2017 Sep 21 12:35:00"))
	c.Assert(m, Equals, 0)

	// should get one for 10th, one for 11th and one for 12th because the endDate is one minute after the third time
	m = s.repository.totalMinutesForTaskAndUser("t", "u", s.pt("2016 Sep 10 10:00:00"), s.pt("2016 Sep 12 12:36:00"))
	c.Assert(m, Equals, 3)
}

// stands for Parse Time
func (s *TimerRepositoryTestSuite) pt(value string) time.Time {
	result, _ := time.Parse(timerRepositoryTestSuiteTimeParseLayout, value)
	return result
}

// Suite lifecycle and callbacks
func (s *TimerRepositoryTestSuite ) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.repository = NewTimerRepository(s.session)
}

func (s *TimerRepositoryTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *TimerRepositoryTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestTimerRepository(t *testing.T) { TestingT(t) }

type TimerRepositoryTestSuite struct {
	env        *utils.Environment
	session    *mgo.Session
	repository *TimerRepository
}

var _ = Suite(&TimerRepositoryTestSuite {})
