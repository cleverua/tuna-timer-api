package data

import (
	"errors"
	"github.com/nlopes/slack"
	"github.com/tuna-timer/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
)

// TeamService todo
type TeamService struct {
	repository     TeamRepositoryInterface
	userRepository *UserRepository
}

// NewTeamService todo
func NewTeamService(session *mgo.Session) *TeamService {
	return &TeamService{
		repository:     NewTeamRepository(session),
		userRepository: NewUserRepository(session),
	}
}

func (s *TeamService) CreateOrUpdateWithSlackOAuthResponse(slackOAuthResponse *slack.OAuthResponse) error {
	team, err := s.repository.FindByExternalID(slackOAuthResponse.TeamID)
	if err != nil {
		return err
	}

	if team == nil {
		team = &models.Team{
			ExternalTeamID: slackOAuthResponse.TeamID,
		}
	}

	team.ExternalTeamName = slackOAuthResponse.TeamName
	team.SlackOAuthResponse = slackOAuthResponse

	err = s.repository.save(team)
	if err != nil {
		return err
	}

	return nil
}

// EnsureTeamSetUp creates Team, User and Project if either is not in database yet
func (s *TeamService) EnsureTeamSetUp(slackCommand *models.SlackCustomCommand) (*models.Team, *models.Project, error) {

	team, err := s.repository.FindByExternalID(slackCommand.TeamID)
	if err != nil {
		return nil, nil, err
	}

	if team == nil {
		return nil, nil, errors.New("Team not found!")
	}

	var reloadTeam = false

	existingProject := s.findProject(team, slackCommand.ChannelID)
	if existingProject == nil {
		err = s.repository.addProject(team, slackCommand.ChannelID, slackCommand.ChannelName)
		if err != nil {
			return nil, nil, err
		}
		reloadTeam = true
	}

	if reloadTeam {
		// not catching the error here since we've once already created or loaded the Team successfully
		team, _ = s.repository.FindByExternalID(slackCommand.TeamID)
		existingProject = s.findProject(team, slackCommand.ChannelID)
	}
	return team, existingProject, nil
}

func (s *TeamService) findProject(team *models.Team, externalProjectID string) *models.Project {
	var result *models.Project
	for _, project := range team.Projects {
		if project.ExternalProjectID == externalProjectID {
			result = project
			break
		}
	}
	return result
}
