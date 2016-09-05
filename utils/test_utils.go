package utils

import (
	"fmt"
)

// Count - returns the number of persistence entiries
func Count(env *Environment, aType interface{}) int {
	count := 0
	env.OrmDB.Model(aType).Count(&count)
	return count
}

// TruncateTables - clears database tables, supposed to be run in test's setup method
func TruncateTables(env *Environment) {
	tablesToTruncate := []string{"projects", "team_users", "teams", "tasks", "timers", "slack_commands"}
	for _, tableName := range tablesToTruncate {
		env.OrmDB.Exec(fmt.Sprintf("truncate table %s cascade", tableName))
	}
}
