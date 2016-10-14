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

	p1Test, err := s.repository.FindActivePassByToken("token")
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

	p1Test, err := s.repository.FindActivePassByToken("token")
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

	p1Test, err := s.repository.FindActivePassByToken("token")
	c.Assert(err, IsNil)
	c.Assert(p1Test, IsNil)
}

func (s *PassRepositoryTestSuite) TestFindActiveByUserID(c *C) {

	now := time.Now()

	userID := bson.NewObjectId()
	s.userRepository.save(&models.TeamUser{
		ID: userID,
	})

	p1 := &models.Pass{ // a good one
		ID:         bson.NewObjectId(),
		Token:      "p1token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(5 * time.Minute),
		ClaimedAt:  nil,
		TeamUserID: userID.Hex(),
	}

	p2 := &models.Pass{ // already claimed
		ID:        bson.NewObjectId(),
		Token:     "p2token",
		CreatedAt: now,
		ExpiresAt: now.Add(5 * time.Minute),
		ClaimedAt: &now,
	}

	p3 := &models.Pass{ // belongs to another user
		ID:         bson.NewObjectId(),
		Token:      "p3token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(5 * time.Minute),
		ClaimedAt:  &now,
		TeamUserID: "another-user",
	}

	s.repository.insert(p1)
	s.repository.insert(p2)
	s.repository.insert(p3)

	pass, err := s.repository.FindActiveByUserID(userID.Hex())
	c.Assert(err, IsNil)
	c.Assert(pass, NotNil)

	c.Assert(pass.Token, Equals, "p1token")
	c.Assert(pass.Token, Equals, "p1token")
}

func (s *PassRepositoryTestSuite) TestRemoveExpiredPasses(c *C) {

	now := time.Now()

	p1 := &models.Pass{ //should be removed as its expiresAt is in the past
		ID:         bson.NewObjectId(),
		Token:      "p1token",
		CreatedAt:  now.Add(-5 * time.Minute),
		ExpiresAt:  now.Add(-3 * time.Minute),
		ClaimedAt:  nil,
		TeamUserID: "user-id",
	}

	p2 := &models.Pass{ //should NOT be removed as its expiresAt is in the future
		ID:         bson.NewObjectId(),
		Token:      "p2token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(5 * time.Minute),
		ClaimedAt:  nil,
		TeamUserID: "user-id",
	}

	claimedAt := now.Add(2 * time.Minute)
	p3 := &models.Pass{ //should NOT be removed as it is claimed
		ID:         bson.NewObjectId(),
		Token:      "p3token",
		CreatedAt:  now,
		ExpiresAt:  now.Add(5 * time.Minute),
		ClaimedAt:  &claimedAt,
		TeamUserID: "user-id",
	}

	err := s.repository.insert(p1)
	c.Assert(err, IsNil)

	err = s.repository.insert(p2)
	c.Assert(err, IsNil)

	err = s.repository.insert(p3)
	c.Assert(err, IsNil)

	err = s.repository.removeExpiredPasses()
	c.Assert(err, IsNil)

	p1, err = s.repository.findByID(p1.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(p1, IsNil) // not found

	p2, err = s.repository.findByID(p2.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(p2, NotNil) // found

	p3, err = s.repository.findByID(p3.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(p3, NotNil) // found
}

func (s *PassRepositoryTestSuite) TestFindByID(c *C) {
	p1 := &models.Pass{
		ID:         bson.NewObjectId(),
		Token:      "p1token",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now(),
		ClaimedAt:  nil,
		TeamUserID: "user-id",
	}

	p2 := &models.Pass{
		ID:         bson.NewObjectId(),
		Token:      "p2token",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now(),
		ClaimedAt:  nil,
		TeamUserID: "user-id",
	}

	err := s.repository.insert(p1)
	c.Assert(err, IsNil)
	err = s.repository.insert(p2)
	c.Assert(err, IsNil)

	p, err := s.repository.findByID(p1.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
	c.Assert(p.Token, Equals, "p1token")

	p, err = s.repository.findByID(bson.NewObjectId().Hex())
	c.Assert(err, IsNil)
	c.Assert(p, IsNil)
}

func (s *PassRepositoryTestSuite) TestRemovePassesClaimedBefore(c *C) {

	now := time.Now()

	fiveMinutesInPast := now.Add(-5 * time.Minute)
	p1 := &models.Pass{
		ID:         bson.NewObjectId(),
		Token:      "p1token",
		CreatedAt:  now,
		ExpiresAt:  now,
		ClaimedAt:  &fiveMinutesInPast,
		TeamUserID: "user-id",
	}

	fiveMinutesInFuture := now.Add(5 * time.Minute)
	p2 := &models.Pass{
		ID:         bson.NewObjectId(),
		Token:      "p2token",
		CreatedAt:  now,
		ExpiresAt:  now,
		ClaimedAt:  &fiveMinutesInFuture,
		TeamUserID: "user-id",
	}

	s.repository.insert(p1)
	s.repository.insert(p2)

	err := s.repository.removePassesClaimedBefore(time.Now())
	c.Assert(err, IsNil)

	p1, err = s.repository.findByID(p1.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(p1, IsNil)

	p2, err = s.repository.findByID(p2.ID.Hex())
	c.Assert(err, IsNil)
	c.Assert(p2, NotNil)
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
	s.userRepository = NewUserRepository(session)
}

func (s *PassRepositoryTestSuite) TearDownSuite(c *C) {
	s.session.Close()
}

func (s *PassRepositoryTestSuite) SetUpTest(c *C) {
	utils.TruncateTables(s.session)
}

func TestPassRepository(t *testing.T) { TestingT(t) }

type PassRepositoryTestSuite struct {
	env            *utils.Environment
	session        *mgo.Session
	repository     *PassRepository
	userRepository *UserRepository
}

var _ = Suite(&PassRepositoryTestSuite{})
