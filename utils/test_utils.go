package utils

import (
	"gopkg.in/mgo.v2"
	"time"
)

const testTimeParseLayout = "2006 Jan 02 15:04:05"

// TruncateTables - clears database tables, supposed to be run in test's setup method
func TruncateTables(session *mgo.Session) {
	tablesToTruncate := []string{
		MongoCollectionTeams,
		MongoCollectionTimers,
		MongoCollectionTeamUsers,
		MongoCollectionPasses,
	}

	for _, tableName := range tablesToTruncate {
		//log.Printf("Truncating table: %s", tableName)
		session.DB("").C(tableName).RemoveAll(nil)
	}
}

// stands for parse time
func PT(value string) time.Time {
	result, _ := time.Parse(testTimeParseLayout, value)
	return result
}
