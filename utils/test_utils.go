package utils

import (
	"log"

	"gopkg.in/mgo.v2"

	"github.com/jinzhu/gorm"
)

// Count - returns the number of persistence entiries
func Count(db *gorm.DB, aType interface{}) int {
	count := 0
	db.Model(aType).Count(&count)
	return count
}

// TruncateTables - clears database tables, supposed to be run in test's setup method
func TruncateTables(session *mgo.Session) {
	tablesToTruncate := []string{"teams", "timers"}
	for _, tableName := range tablesToTruncate {
		log.Printf("Truncating table: %s", tableName)
		session.DB("").C(tableName).RemoveAll(nil)
	}
}
