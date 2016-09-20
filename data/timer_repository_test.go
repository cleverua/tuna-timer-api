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

	// creates 10 timers one minute each
	for i := 10; i < 20; i++ {
		createdAt := s.pt(fmt.Sprintf("2016 Sep %d 12:35:00", i))
		s.repository.createTimer(&models.Timer{
			ID:               bson.NewObjectId(),
			TeamID:           "teamID",
			ProjectID:        "projectID",
			TeamUserID:       "u",
			CreatedAt:        createdAt,
			FinishedAt:		  &createdAt,
			TaskHash:   	  "h",
			Minutes:  		  1,
		})
	}

	m := s.repository.totalMinutesForTaskAndUser("h", "u", s.pt("2016 Sep 09 12:35:00"), s.pt("2016 Sep 21 12:35:00"))
	c.Assert(m, Equals, 0)

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
