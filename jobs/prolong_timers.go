package jobs

import (
	"fmt"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
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
}

