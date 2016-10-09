package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
)

func (s *PassRepositoryTestSuite) TestFindByToken(c *C) {
	p1 := &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "token",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(5 * time.Minute),
		ClaimedAt:    nil,
		ModelVersion: models.ModelVersionPass,
	}
	err := s.repository.insert(p1)
	c.Assert(err, IsNil)

	p1Test, err := s.repository.FindByToken("token")
	c.Assert(err, IsNil)
	c.Assert(p1Test, NotNil)

	c.Assert(p1.ID, Equals, p1Test.ID)
}

func (s *PassRepositoryTestSuite) TestFindByTokenDoesNotGetExpired(c *C) {
	p1 := &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "token",
		CreatedAt:    time.Now().Add(-10 * time.Minute),
		ExpiresAt:    time.Now().Add(-5 * time.Minute),
		ClaimedAt:    nil,
		ModelVersion: models.ModelVersionPass,
	}
	err := s.repository.insert(p1)
	c.Assert(err, IsNil)

	p1Test, err := s.repository.FindByToken("token")
	c.Assert(err, IsNil)
	c.Assert(p1Test, IsNil)
}

func (s *PassRepositoryTestSuite) TestFindByTokenDoesNotGetClaimed(c *C) {
	now := time.Now()

	p1 := &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "token",
		CreatedAt:    now,
		ExpiresAt:    now.Add(5 * time.Minute),
		ClaimedAt:    &now,
		ModelVersion: models.ModelVersionPass,
	}
	err := s.repository.insert(p1)
	c.Assert(err, IsNil)

	p1Test, err := s.repository.FindByToken("token")
	c.Assert(err, IsNil)
	c.Assert(p1Test, IsNil)
}

// Suite lifecycle and callbacks
func (s *PassRepositoryTestSuite) SetUpSuite(c *C) {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	e.MigrateDatabase(session)

	s.env = e
	s.session = session.Clone()
	s.repository = NewPassRepository(s.session)
}

func (s *PassRepositoryTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *PassRepositoryTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestPassRepository(t *testing.T) { TestingT(t) }

type PassRepositoryTestSuite struct {
	env        *utils.Environment
	session    *mgo.Session
	repository *PassRepository
}

var _ = Suite(&PassRepositoryTestSuite{})
