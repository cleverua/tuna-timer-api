package web

import (
	"gopkg.in/mgo.v2"
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/dgrijalva/jwt-go"
)

type JwtToken struct{
	Token string `json:"jwt"`
}

func NewToken(userId string, session *mgo.Session) (string, error) {
	user_service := data.NewUserService(session)
	user, user_err := user_service.FindByID(userId)

	if user_err != nil {
		return "", user_err
	}

	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"team_id": user.TeamID,
		"user_id": user.ID,
		"is_team_admin": user.SlackUserInfo.IsAdmin,
		"image48": user.SlackUserInfo.Profile.Image48,
	})

	signed_token, err := jwt_token.SignedString([]byte("TODO: Extract me in config/env"))
	return signed_token, err
}
