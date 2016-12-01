package web

import (
	"testing"
	"github.com/cleverua/tuna-timer-api/utils"
	"github.com/cleverua/tuna-timer-api/models"
	"gopkg.in/tylerb/is.v1"
	"log"
	"gopkg.in/mgo.v2"
	"github.com/pavlo/gosuite"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

func TestJwtToken(t *testing.T) {
	gosuite.Run(t, &JwtTokenTestSuite{Is: is.New(t)})
}

func (s *JwtTokenTestSuite) TestNewUserToken(t *testing.T) {
	jwt_token, err := NewUserToken(s.pass.TeamUserID, s.session)
	s.Nil(err)

	new_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"image48": s.user.SlackUserInfo.Profile.Image48,
		"team_id": s.user.TeamID,
		"user_id": s.user.ID,
		"is_team_admin": s.user.SlackUserInfo.IsAdmin,
	})
	verification_token, err := new_token.SignedString([]byte("TODO: Extract me in config/env"))

	s.Nil(err)
	s.Equal(jwt_token, verification_token)
}

func (s *JwtTokenTestSuite) TestNewUserTokenFail(t *testing.T) {
	id := bson.NewObjectId().Hex()
	jwt_token, err := NewUserToken(id, s.session)

	s.Err(err)
	s.Equal(err.Error(), "user doesn't exist")
	s.Zero(jwt_token)
}

type JwtTokenTestSuite struct {
	*is.Is
	env        *utils.Environment
	session    *mgo.Session
	user       models.TeamUser
	pass       models.Pass
}
func (s *JwtTokenTestSuite) SetUpSuite() {
	e := utils.NewEnvironment(utils.TestEnv, "1.0.0")

	session, err := utils.ConnectToDatabase(e.Config)
	if err != nil {
		log.Fatal("Failed to connect to DB!")
	}

	s.session = session.Clone()
	e.MigrateDatabase(session)
	s.env = e
}

func (s *JwtTokenTestSuite) TearDownSuite() {
	s.session.Close()
}

func (s *JwtTokenTestSuite) SetUp() {
	//Clear Database
	utils.TruncateTables(s.session)

	//Seed Database
	teamUser, err := utils.Create(&models.TeamUser{}, s.session)
	s.Nil(err)
	s.user = teamUser.(models.TeamUser)

	pass, err := utils.Create(&models.Pass{}, s.session)
	s.Nil(err)
	s.pass = pass.(models.Pass)
}

func (s *JwtTokenTestSuite) TearDown() {}
