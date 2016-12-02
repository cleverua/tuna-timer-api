package web

import (
	"gopkg.in/mgo.v2"
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/dgrijalva/jwt-go"
	"errors"
)

type JwtToken struct{
	Token string `json:"jwt"`
}

func NewUserToken(userId string, session *mgo.Session) (string, error) {
	userService := data.NewUserService(session)
	user, userErr := userService.FindByID(userId)

	if userErr == nil && user == nil {
		return "", errors.New("user doesn't exist")
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"team_id": user.TeamID,
		"user_id": user.ID,
		"is_team_admin": user.SlackUserInfo.IsAdmin,
		"image48": user.SlackUserInfo.Profile.Image48,
	})

	signedToken, err := jwtToken.SignedString([]byte("TODO: Extract me in config/env"))
	return signedToken, err
}
