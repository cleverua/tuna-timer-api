package data

import (
	"github.com/tuna-timer/tuna-timer-api/models"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PassRepository struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func NewPassRepository(session *mgo.Session) *PassRepository {
	return &PassRepository{
		session:    session,
		collection: session.DB("").C(utils.MongoCollectionPasses),
	}
}

func (r *PassRepository) FindByToken(token string) (*models.Pass, error) {
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

func (r *PassRepository) insert(pass *models.Pass) error {
	return r.collection.Insert(pass)
}

func (r *PassRepository) update(pass *models.Pass) error {
	return r.collection.Update(bson.M{"_id": pass.ID}, pass)
}
