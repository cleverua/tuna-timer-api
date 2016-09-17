package utils

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Count - returns the number of persistence entiries
func Count(db *gorm.DB, aType interface{}) int {
	count := 0
	db.Model(aType).Count(&count)
	return count
}

// TruncateTables - clears database tables, supposed to be run in test's setup method
func TruncateTables(db *gorm.DB) {
	tablesToTruncate := []string{"projects", "team_users", "teams", "tasks", "timers", "slack_commands"}
	for _, tableName := range tablesToTruncate {
		db.Exec(fmt.Sprintf("truncate table %s cascade", tableName))
	}
}
