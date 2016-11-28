package data

import (
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
	"gopkg.in/tylerb/is.v1"
	"github.com/pavlo/gosuite"
)

func TestPassRepository(t *testing.T) {
	gosuite.Run(t, &PassRepositoryTestSuite{Is: is.New(t)})
}

func (s *PassRepositoryTestSuite) TestFindByToken(t *testing.T) {
	p1 := &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "token",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(5 * time.Minute),
		ClaimedAt:    nil,
		ModelVersion: models.ModelVersionPass,
	}
	err := s.repository.insert(p1)
	s.Nil(err)

	p1Test, err := s.repository.FindActivePassByToken("token")
	s.Nil(err)
	s.NotNil(p1Test)

	s.Equal(p1.ID, p1Test.ID)
}

func (s *PassRepositoryTestSuite) TestFindByTokenDoesNotGetExpired(t *testing.T) {
	p1 := &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "token",
		CreatedAt:    time.Now().Add(-10 * time.Minute),
		ExpiresAt:    time.Now().Add(-5 * time.Minute),
		ClaimedAt:    nil,
		ModelVersion: models.ModelVersionPass,
	}
	err := s.repository.insert(p1)
	s.Nil(err)

	p1Test, err := s.repository.FindActivePassByToken("token")
	s.Nil(err)
	s.Nil(p1Test)
}

func (s *PassRepositoryTestSuite) TestFindByTokenDoesNotGetClaimed(t *testing.T) {
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
	s.Nil(err)

	p1Test, err := s.repository.FindActivePassByToken("token")
	s.Nil(err)
	s.Nil(p1Test)
}

func (s *PassRepositoryTestSuite) TestFindActiveByUserID(t *testing.T) {

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
	s.Nil(err)
	s.NotNil(pass)
	s.Equal("p1token", pass.Token)
}

func (s *PassRepositoryTestSuite) TestRemoveExpiredPasses(t *testing.T) {

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
	s.Nil(err)

	err = s.repository.insert(p2)
	s.Nil(err)

	err = s.repository.insert(p3)
	s.Nil(err)

	err = s.repository.removeExpiredPasses()
	s.Nil(err)

	p1, err = s.repository.findByID(p1.ID.Hex())
	s.Nil(err)
	s.Nil(p1)

	p2, err = s.repository.findByID(p2.ID.Hex())
	s.Nil(err)
	s.NotNil(p2)

	p3, err = s.repository.findByID(p3.ID.Hex())
	s.Nil(err)
	s.NotNil(p3)
}

func (s *PassRepositoryTestSuite) TestFindByID(t *testing.T) {
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
	s.Nil(err)

	err = s.repository.insert(p2)
	s.Nil(err)

	p, err := s.repository.findByID(p1.ID.Hex())
	s.Nil(err)
	s.NotNil(p)
	s.Equal("p1token", p.Token)

	p, err = s.repository.findByID(bson.NewObjectId().Hex())
	s.Nil(err)
	s.Nil(p)
}

func (s *PassRepositoryTestSuite) TestRemovePassesClaimedBefore(t *testing.T) {

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
	s.Nil(err)

	p1, err = s.repository.findByID(p1.ID.Hex())
	s.Nil(err)
	s.Nil(p1)

	p2, err = s.repository.findByID(p2.ID.Hex())
	s.Nil(err)
	s.NotNil(p2)
}

type PassRepositoryTestSuite struct {
	*is.Is
	env            *utils.Environment
	session        *mgo.Session
	repository     *PassRepository
	userRepository *UserRepository
}

func (s *PassRepositoryTestSuite) SetUpSuite() {
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

func (s *PassRepositoryTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *PassRepositoryTestSuite) SetUp() {
	utils.TruncateTables(s.session)
}

func (s *PassRepositoryTestSuite) TearDown() {}
