package jobs

import (
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"log"
)

type ClearPasses struct {
	env     *utils.Environment
	session *mgo.Session
}

func NewClearPasses(env *utils.Environment, session *mgo.Session) *ClearPasses {
	return &ClearPasses{
		env:     env,
		session: session,
	}
}

func (j *ClearPasses) Run() {
	log.Println("ClearPasses launched!")

	service := data.NewPassService(j.session)
	service.RemoveStalePasses()

	log.Println("ClearPasses finished!")
}
