package jobs

import (
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type StopTimersAtMidnight struct {
	env     *utils.Environment
	session *mgo.Session
}

func NewStopTimersAtMidnight(env *utils.Environment, session *mgo.Session) *StopTimersAtMidnight {
	return &StopTimersAtMidnight{
		env:     env,
		session: session,
	}
}

func (j *StopTimersAtMidnight) Run() {
	log.Println("ProlongTimersJob launched!")

	now := time.Now()

	service := data.NewTimerService(j.session)
	service.CompleteActiveTimersAtMidnight(&now)

	log.Println("ProlongTimersJob finished!")
}
