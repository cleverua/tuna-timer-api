package utils

import (
	"log"

	"gopkg.in/mgo.v2"
)

// TruncateTables - clears database tables, supposed to be run in test's setup method
func TruncateTables(session *mgo.Session) {
	tablesToTruncate := []string{"teams", "timers"}
	for _, tableName := range tablesToTruncate {
		log.Printf("Truncating table: %s", tableName)
		session.DB("").C(tableName).RemoveAll(nil)
	}
}
