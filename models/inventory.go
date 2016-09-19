package models

// StartCommandInventory todo
type StartCommandInventory struct {
	Team                             *Team
	Project                          *Project
	TeamUser                         *TeamUser
	StoppedTimer                     *Timer
	StartedTimer                     *Timer
	AlreadyStartedTimer              *Timer
	StoppedTaskTotalForToday         int
	StartedTaskTotalForToday         int
	AlreadyStartedTimerTotalForToday int
	Resumed                          bool
	UserTotalForToday                int
}
