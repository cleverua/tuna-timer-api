package web

import (
	"gopkg.in/mgo.v2"
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/dgrijalva/jwt-go"
)

type JwtToken struct{
	Token string `json:"jwt"`
}

func NewUserToken(userId string, session *mgo.Session) (string, error) {
	userService := data.NewUserService(session)
	teamService := data.NewTeamService(session)

	user, err := userService.FindByID(userId)
	if err != nil {
		return "", err
	}

	userTeam, err := teamService.FindByID(user.TeamID)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":	 user.ID,
		"name":		 user.ExternalUserName,
		"is_team_admin": user.SlackUserInfo.IsAdmin,
		"image48":	 user.SlackUserInfo.Profile.Image48,
		"team_id":	 userTeam.ID,
		"ext_team_id":	 userTeam.ExternalTeamID,
		"ext_team_name": userTeam.ExternalTeamName,
	})

	signedToken, err := jwtToken.SignedString([]byte("TODO: Extract me in config/env"))
	return signedToken, err
}
