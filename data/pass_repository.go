package data

import (
	"github.com/cleverua/tuna-timer-api/models"
	"github.com/cleverua/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PassRepository struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type StalePassesRemover interface {
	removeExpiredPasses() error
	removePassesClaimedBefore(date time.Time) error
}

func NewPassRepository(session *mgo.Session) *PassRepository {
	return &PassRepository{
		session:    session,
		collection: session.DB("").C(utils.MongoCollectionPasses),
	}
}

func (r *PassRepository) FindActivePassByToken(token string) (*models.Pass, error) {
	pass := &models.Pass{}
	err := r.collection.Find(bson.M{
		"token":      token,
		"expires_at": bson.M{"$gt": time.Now()},
		"claimed_at": nil,
	}).One(pass)

	if err != nil && err == mgo.ErrNotFound {
		pass = nil
		err = nil
	}
	return pass, err
}

func (r *PassRepository) FindActiveByUserID(userID string) (*models.Pass, error) {
	pass := &models.Pass{}

	err := r.collection.Find(bson.M{
		"team_user_id": userID,
		"expires_at":   bson.M{"$gt": time.Now()},
		"claimed_at":   nil,
	}).One(pass)

	if err != nil && err == mgo.ErrNotFound {
		pass = nil
		err = nil
	}
	return pass, err
}

func (r *PassRepository) removeExpiredPasses() error {
	return r.collection.Remove(bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
		"claimed_at": nil,
	})
}

func (r *PassRepository) removePassesClaimedBefore(date time.Time) error {
	return r.collection.Remove(bson.M{
		"claimed_at": bson.M{"$lt": date},
	})
}

func (r *PassRepository) insert(pass *models.Pass) error {
	return r.collection.Insert(pass)
}

func (r *PassRepository) update(pass *models.Pass) error {
	return r.collection.Update(bson.M{"_id": pass.ID}, pass)
}

func (r *PassRepository) findByID(passID string) (*models.Pass, error) {
	result := &models.Pass{}
	err := r.collection.FindId(bson.ObjectIdHex(passID)).One(&result)

	if err != nil && err == mgo.ErrNotFound {
		result = nil
		err = nil
	}

	return result, err
}
