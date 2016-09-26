package data

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (s *TimerServiceTestSuite) TestGetActiveTimer(c *C) {

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
	c.Assert(err, IsNil)
	c.Assert(timer, NotNil)
	c.Assert(timer.Minutes, Equals, 20)
}

func (s *TimerServiceTestSuite) TestStopTimer(c *C) {
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

	c.Assert(err, IsNil)
	c.Assert(timer, NotNil)

	s.service.StopTimer(timer)

	loadedTimer, err := s.repo.findByID(id.Hex())
	c.Assert(err, IsNil)

	c.Assert(loadedTimer.Minutes, Equals, 20)
	c.Assert(loadedTimer.FinishedAt, NotNil)
}

func (s *TimerServiceTestSuite) TestStartTimer(c *C) {
	timer, err := s.service.StartTimer("team", "project", "user", "task")
	c.Assert(err, IsNil)
	c.Assert(timer, NotNil)

	loadedTimer, err := s.repo.findByID(timer.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(loadedTimer, NotNil)

	c.Assert(loadedTimer.TeamID, Equals, "team")
	c.Assert(loadedTimer.ProjectID, Equals, "project")
	c.Assert(loadedTimer.TeamUserID, Equals, "user")
	c.Assert(loadedTimer.TaskName, Equals, "task")
	c.Assert(loadedTimer.TaskHash, NotNil)
	c.Assert(loadedTimer.CreatedAt, NotNil)
	c.Assert(loadedTimer.FinishedAt, IsNil)
	c.Assert(loadedTimer.DeletedAt, IsNil)
	c.Assert(loadedTimer.Minutes, Equals, 0)
}

func (s *TimerServiceTestSuite) TestTotalMinutesForTodayAddsTimeForUnfinishedTask(c *C) {
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

	c.Assert(s.service.TotalMinutesForTaskToday(timer), Equals, 15)
}

func (s *TimerServiceTestSuite) TestTotalMinutesForUserToday(c *C) {
	now := time.Now()

	duration, _ := time.ParseDuration("5m")
	secondTimerStartedAt := now.Add(duration * -1) // 5 minutes ago

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

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  secondTimerStartedAt,
		FinishedAt: nil,
		Minutes:    0,
	})

	c.Assert(s.service.TotalUserMinutesForDay("user", time.Now()), Equals, 15)
	c.Assert(s.service.TotalUserMinutesForDay("this user has no tasks", time.Now()), Equals, 0)
}

func (s *TimerServiceTestSuite) TestTotalMinutesForUserTodayWhenOneTasksLastsSinceYesterday(c *C) {
	now := time.Now()

	duration, _ := time.ParseDuration("25h")
	startedAt := now.Add(duration * -1) // 25 hours ago ago

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  startedAt,
		FinishedAt: nil,
		Minutes:    0,
	})

	actual := now.Hour()*60 + now.Minute()
	c.Assert(s.service.TotalUserMinutesForDay("user", time.Now()), Equals, actual)
}

func (s *TimerServiceTestSuite) TestGetCompletedTasksForDay(c *C) {

	now := time.Now()

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 10:35:00"),
		FinishedAt: &now,
		Minutes:    2,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 12 23:59:59"),
		FinishedAt: &now,
		Minutes:    3,
	})

	s.repo.createTimer(&models.Timer{
		ID:         bson.NewObjectId(),
		TeamID:     "team",
		ProjectID:  "project",
		TeamUserID: "user",
		TaskHash:   "task",
		CreatedAt:  utils.PT("2016 Sep 13 00:00:00"),
		FinishedAt: &now,
		Minutes:    7,
	})

	v, err := s.service.GetCompletedTasksForDay("user", utils.PT("2016 Sep 12 23:59:59"))
	c.Assert(err, IsNil)
	c.Assert(len(v), Equals, 1)
	c.Assert(v[0].Minutes, Equals, 5)

	v, err = s.service.GetCompletedTasksForDay("user", utils.PT("2016 Sep 13 23:59:59"))
	c.Assert(err, IsNil)
	c.Assert(len(v), Equals, 1)
	c.Assert(v[0].Minutes, Equals, 7)

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
	s.repo = NewTimerRepository(s.session)
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
	env     *utils.Environment
	session *mgo.Session
	repo    *TimerRepository
	service *TimerService
}

var _ = Suite(&TimerServiceTestSuite{})
