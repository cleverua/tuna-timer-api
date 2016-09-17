package commands

import (
	"context"
	"log"

	"github.com/pavlo/slack-time/data"
	"github.com/pavlo/slack-time/models"
	"github.com/pavlo/slack-time/utils"
)

//Stop - handles the '/timer stop` command received from Slack
type Stop struct {
}

// cases:
// 1. Successfully stopped a timer
// 2. No currently ticking timer existed
// 3. Any other errors

// Handle - SlackCustomCommandHandler interface
func (c *Stop) Handle(ctx context.Context, slackCommand models.SlackCustomCommand) *SlackCustomCommandHandlerResult {

	dataService := data.CreateDataService()
	db := utils.GetDBTransactionFromContext(ctx)
	team, user, project := dataService.CreateTeamAndUserAndProject(db, slackCommand)

	log.Printf("%v %v %v", team, user, project)
	return nil
}
