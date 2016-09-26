package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
)

// TeamService todo
type TeamService struct {
	session    *mgo.Session
	repository TeamRepositoryInterface
}

// NewTeamService todo
func NewTeamService(session *mgo.Session) *TeamService {
	return &TeamService{
		session:    session,
		repository: NewTeamRepository(session),
	}
}

// EnsureTeamSetUp creates Team, User and Project if either is not in database yet
func (s *TeamService) EnsureTeamSetUp(slackCommand *models.SlackCustomCommand) (*models.Team, *models.Project, *models.TeamUser, error) {

	team, err := s.repository.findByExternalID(slackCommand.TeamID)
	if err != nil {
		return nil, nil, nil, err
	}

	if team == nil {
		team, err = s.repository.createTeam(slackCommand.TeamID, slackCommand.TeamDomain)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	var reloadTeam = false

	existingProject := s.findProject(team, slackCommand.ChannelID)
	if existingProject == nil {
		err = s.repository.addProject(team, slackCommand.ChannelID, slackCommand.ChannelName)
		if err != nil {
			return nil, nil, nil, err
		}
		reloadTeam = true
	}

	existingUser := s.findUser(team, slackCommand.UserID)
	if existingUser == nil {
		err = s.repository.addUser(team, slackCommand.UserID, slackCommand.UserName)
		if err != nil {
			return nil, nil, nil, err
		}
		reloadTeam = true
	}

	if reloadTeam {
		// not catching the error here since we've once already created or loaded the Team successfully
		team, _ = s.repository.findByExternalID(slackCommand.TeamID)
		existingProject = s.findProject(team, slackCommand.ChannelID)
		existingUser = s.findUser(team, slackCommand.UserID)
	}
	return team, existingProject, existingUser, nil
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

func (s *TeamService) findUser(team *models.Team, externalUserID string) *models.TeamUser {
	var result *models.TeamUser
	for _, user := range team.Users {
		if user.ExternalUserID == externalUserID {
			result = user
			break
		}
	}
	return result
}
