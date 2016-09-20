package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/pavlo/slack-time/utils"
	. "gopkg.in/check.v1"
	"time"
)




// Team User
func (s *TimerServiceTestSuite) TestTotalMinutesForToday(c *C) {

	s.service.TotalMinutesForToday(nil)
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
	service    *TimerService
}

var _ = Suite(&TimerServiceTestSuite {})
