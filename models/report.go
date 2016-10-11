package models

// StartCommandInventory collect everything that StartCommand creates, modifies or touches
// This report instance will be sent to a UITheme to format a slack reply to the Start Command
type StartCommandReport struct {
	Team                             *Team
	Project                          *Project
	TeamUser                         *TeamUser
	Pass                             *Pass
	StoppedTimer                     *Timer
	StartedTimer                     *Timer
	AlreadyStartedTimer              *Timer
	StoppedTaskTotalForToday         int
	StartedTaskTotalForToday         int
	AlreadyStartedTimerTotalForToday int
	Resumed                          bool
	UserTotalForToday                int
}

type StopCommandReport struct {
	Team                     *Team
	Project                  *Project
	TeamUser                 *TeamUser
	Pass                     *Pass
	StoppedTimer             *Timer
	StoppedTaskTotalForToday int
	UserTotalForToday        int
}

type StatusCommandReport struct {
	Team                             *Team
	Project                          *Project
	TeamUser                         *TeamUser
	Pass                             *Pass
	Tasks                            []*TaskAggregation
	AlreadyStartedTimer              *Timer
	AlreadyStartedTimerTotalForToday int
	PeriodName                       string // `today`, `yesterday`, `MM-DD-YYYY`
	UserTotalForPeriod               int
}
