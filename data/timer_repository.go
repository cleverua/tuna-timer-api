package data

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/cleverua/tuna-timer-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const timersCollectionName = "timers"

// TimerRepository todo
type TimerRepository struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// NewTimerRepository todo
func NewTimerRepository(session *mgo.Session) *TimerRepository {
	return &TimerRepository{
		session:    session,
		collection: session.DB("").C(timersCollectionName),
	}
}

func (r *TimerRepository) findByID(timerID string) (*models.Timer, error) {
	result := &models.Timer{}
	err := r.collection.FindId(bson.ObjectIdHex(timerID)).One(&result)
	return result, err
}

func (r *TimerRepository) findActiveByTimezoneOffset(timezoneOffset int) ([]*models.Timer, error) {
	result := []*models.Timer{}

	err := r.collection.Find(bson.M{
		"tz_offset":   timezoneOffset,
		"finished_at": nil,
		"deleted_at":  nil}).All(&result)

	return result, err
}

func (r *TimerRepository) findActiveByTeamAndUser(teamID, userID string) (*models.Timer, error) {

	result := &models.Timer{}

	err := r.collection.Find(bson.M{
		"team_id":      teamID,
		"team_user_id": userID,
		"finished_at":  nil,
		"deleted_at":   nil}).One(result)

	if err != nil && err == mgo.ErrNotFound {
		result = nil
		err = nil
	}
	return result, err
}

func (r *TimerRepository) findActiveByUser(userID string) (*models.Timer, error) {

	result := &models.Timer{}

	err := r.collection.Find(bson.M{
		"team_user_id": userID,
		"finished_at":  nil,
		"deleted_at":   nil}).One(result)

	if err != nil && err == mgo.ErrNotFound {
		result = nil
		err = nil
	}
	return result, err
}

func (r *TimerRepository) create(teamID string, project *models.Project, user *models.TeamUser, taskName string) (*models.Timer, error) {

	timer := &models.Timer{
		ID:                  bson.NewObjectId(),
		TeamID:              teamID,
		ProjectID:           project.ID.Hex(),
		ProjectExternalName: project.ExternalProjectName,
		ProjectExternalID:   project.ExternalProjectID,
		TeamUserID:          user.ID.Hex(),
		TeamUserTZOffset:    user.SlackUserInfo.TZOffset,
		CreatedAt:           time.Now(),
		TaskName:            taskName,
		TaskHash:            taskSHA256(teamID, project.ID.Hex(), taskName),
		Minutes:             0,
		ModelVersion:        models.ModelVersionTimer,
	}

	return r.CreateTimer(timer)
}

/*
db.getCollection('data').aggregate([
{
	$match: {
		'task_hash': 'nisl',
		'created_at' : {
			$gte: ISODate("2014-09-30T03:44:54.000Z"),
			$lt: ISODate("2017-09-30T03:44:54.000Z")
		},
		'deleted_at': null,
		'team_user_id': '12341234'
	}
},
{
	$group: {
		_id: { task_hash: '$task_hash' },
		minutes: { $sum: '$minutes' },
		total_timers: { $sum: 1 },
	}
}])

*/
func (r *TimerRepository) totalMinutesForTaskAndUser(taskHash, userID string, startDate, endDate time.Time) int {
	pipeConfig := []map[string]interface{}{
		{
			"$match": bson.M{
				"task_hash":    taskHash,
				"team_user_id": userID,
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":          bson.M{"task_hash": "$task_hash"},
				"minutes":      bson.M{"$sum": "$minutes"},
				"total_timers": bson.M{"$sum": 1},
			},
		},
	}

	var result map[string]interface{}
	err := r.collection.Pipe(pipeConfig).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Printf("Error: %s", err)
	}

	if result == nil {
		return 0
	}

	return result["minutes"].(int)
}

func (r *TimerRepository) totalMinutesForUser(userID string, startDate, endDate time.Time) int {
	pipeConfig := []map[string]interface{}{
		{
			"$match": bson.M{
				"team_user_id": userID,
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":          bson.M{"user_id": "$team_user_id"},
				"minutes":      bson.M{"$sum": "$minutes"},
				"total_timers": bson.M{"$sum": 1},
			},
		},
	}

	var result map[string]interface{}
	err := r.collection.Pipe(pipeConfig).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Printf("Error: %s", err)
	}

	if result == nil {
		return 0
	}

	return result["minutes"].(int)
}

func (r *TimerRepository) completedTasksForUser(userID string, startDate, endDate time.Time) ([]*models.TaskAggregation, error) {

	pipeConfig := []map[string]interface{}{
		{
			"$match": bson.M{
				"team_user_id": userID,
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
				"finished_at": bson.M{"$ne": nil},
				"deleted_at":  nil,
			},
		},
		{
			"$sort": bson.M{"created_at": -1},
		},
		{
			"$group": bson.M{
				"_id":     bson.M{"task_name": "$task_name", "task_hash": "$task_hash", "project_ext_name": "$project_ext_name", "project_ext_id": "$project_ext_id"},
				"minutes": bson.M{"$sum": "$minutes"},
			},
		},
		{
			"$project": bson.M{
				"_id":              0,
				"task_name":        "$_id.task_name",
				"minutes":          "$minutes",
				"task_hash":        "$_id.task_hash",
				"project_ext_name": "$_id.project_ext_name",
				"project_ext_id":   "$_id.project_ext_id",
			},
		},
	}

	var results []*models.TaskAggregation
	err := r.collection.Pipe(pipeConfig).All(&results)
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	return results, nil
}

func (r *TimerRepository) CreateTimer(timer *models.Timer) (*models.Timer, error) {
	err := r.collection.Insert(timer)
	return timer, err
}

func (r *TimerRepository) update(timer *models.Timer) error {
	return r.collection.UpdateId(timer.ID, timer)
}

// split into two - hash and trim?
func taskSHA256(teamID, projectID, taskName string) string {
	hashSeed := fmt.Sprintf("%s%s%s", teamID, projectID, taskName)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hashSeed)))[0:6]
}

func (r *TimerRepository) findUserTasksByRange(userID string, startDate, endDate time.Time) ([]*models.Timer, error) {
	var results []*models.Timer

	err := r.collection.Find(bson.M{
		"team_user_id": userID,
		"created_at":  	bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}).Sort("created_at").All(&results)

	return results, err
}
