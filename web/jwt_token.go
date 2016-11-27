package web

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"github.com/tuna-timer/tuna-timer-api/data"
	"github.com/dgrijalva/jwt-go"
)

func NewToken(pass *models.Pass, session *mgo.Session) (string, error) {
	user_service := data.NewUserService(session)
	user, user_err := user_service.FindByID(pass.TeamUserID)

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
