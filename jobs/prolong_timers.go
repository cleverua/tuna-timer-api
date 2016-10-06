package jobs

import (
	"fmt"
	"github.com/tuna-timer/tuna-timer-api/data"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"time"
)

type ProlongTimersJob struct {
	env     *utils.Environment
	session *mgo.Session
}

func NewProlongTimersJob(env *utils.Environment, session *mgo.Session) *ProlongTimersJob {
	return &ProlongTimersJob{
		env:     env,
		session: session,
	}
}

func (j *ProlongTimersJob) Run() {
	fmt.Println("ProlongTimersJob launched!")

	now := time.Now()

	service := data.NewTimerService(j.session)
	service.CompleteActiveTimersAtMidnight(&now)

	fmt.Println("ProlongTimersJob finished!")
}
