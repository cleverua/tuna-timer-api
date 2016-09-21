package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/pavlo/slack-time/utils"
	. "gopkg.in/check.v1"
	"time"
	"github.com/pavlo/slack-time/models"
	"gopkg.in/mgo.v2/bson"
)




// Team User
func (s *TimerServiceTestSuite) TestTotalMinutesForTodayAddsTimeForUnfinishedTask(c *C) {
	now := time.Now()

	offsetDuration1, _ := time.ParseDuration("20m")
	firstTimerStartedAt := now.Add(offsetDuration1 * -1) // 20 minutes ago

	offsetDuration2, _ := time.ParseDuration("5m")
	secondTimerStartedAt := now.Add(offsetDuration2 * -1) // 5 minutes ago

	s.repository.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "u",
		TaskHash:   	  "t",
		CreatedAt:        now.Add(offsetDuration1 * -1),
		FinishedAt:		  &firstTimerStartedAt,
		Minutes:  		  10,
	})

	timer, _ := s.repository.createTimer(&models.Timer{
		ID:               bson.NewObjectId(),
		TeamID:           "teamID",
		ProjectID:        "projectID",
		TeamUserID:       "u",
		TaskHash:   	  "t",
		CreatedAt:        secondTimerStartedAt,
		FinishedAt:		  nil,
		Minutes:  		  0,
	})

	c.Assert(s.service.TotalMinutesForTaskToday(timer), Equals, 15)
}


// Suite lifecycle and callbacks
func (s *TimerServiceTestSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.service = NewTimerService(s.session)
	s.repository = NewTimerRepository(s.session)
}

func (s *TimerServiceTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *TimerServiceTestSuite) SetUpTest(c *C) {
	time.Local = time.UTC
	utils.TruncateTables(s.session)
}

func TestTimerService(t *testing.T) { TestingT(t) }

type TimerServiceTestSuite struct {
	env        *utils.Environment
	session    *mgo.Session
	repository *TimerRepository
	service    *TimerService
}

var _ = Suite(&TimerServiceTestSuite {})
