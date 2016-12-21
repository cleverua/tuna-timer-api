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
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/nlopes/slack"
)

func TestJwtToken(t *testing.T) {
	gosuite.Run(t, &JwtTokenTestSuite{Is: is.New(t)})
}

func (s *JwtTokenTestSuite) TestNewUserToken(t *testing.T) {
	jwtToken, err := NewUserToken(s.pass.TeamUserID, s.session)
	s.Nil(err)

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":	 s.user.ID,
		"is_team_admin": s.user.SlackUserInfo.IsAdmin,
		"name":		 s.user.ExternalUserName,
		"image48":	 s.user.SlackUserInfo.Profile.Image48,
		"team_id":	 s.team.ID,
		"ext_team_id":	 s.team.ExternalTeamID,
		"ext_team_name": s.team.ExternalTeamName,

	})
	verificationToken, err := newToken.SignedString([]byte("TODO: Extract me in config/env"))

	s.Nil(err)
	s.Equal(jwtToken, verificationToken)
}

func (s *JwtTokenTestSuite) TestNewUserTokenFail(t *testing.T) {
	id := bson.NewObjectId().Hex()
	jwtToken, err := NewUserToken(id, s.session)

	s.Err(err)
	s.Equal(err.Error(), mgo.ErrNotFound.Error())
	s.Zero(jwtToken)
}

type JwtTokenTestSuite struct {
	*is.Is
	env     *utils.Environment
	session *mgo.Session
	user    *models.TeamUser
	pass    *models.Pass
	team	*models.Team

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
	passRepository := data.NewPassRepository(s.session)
	userRepository := data.NewUserRepository(s.session)
	teamRepository := data.NewTeamRepository(s.session)

	var err error

	//Create team
	s.team, err = teamRepository.CreateTeam("ExtTeamID", "ExtTeamName")
	s.Nil(err)

	//Create user
	s.user = &models.TeamUser{
		TeamID:           s.team.ID.Hex(),
		ExternalUserID:   "ext-user-id",
		ExternalUserName: "user-name",
		SlackUserInfo:    &slack.User{
			IsAdmin: true,
		},
	}
	_, err = userRepository.Save(s.user)
	s.Nil(err)

	//Create pass
	s.pass = &models.Pass{
		ID:           bson.NewObjectId(),
		Token:        "token",
		TeamUserID:   s.user.ID.Hex(),
	}
	err = passRepository.Insert(s.pass)
	s.Nil(err)
}

func (s *JwtTokenTestSuite) TearDown() {}
