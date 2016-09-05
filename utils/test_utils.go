package utils

import (
	"fmt"
)

func Count(env *Environment, aType interface{}) int {
	count := 0
	env.OrmDB.Model(aType).Count(&count)
	return count
}

func TruncateTables(env *Environment) {
	tablesToTruncate := []string{"teams", "team_users", "projects"}
	for _, tableName := range tablesToTruncate {
		env.OrmDB.Exec(fmt.Sprintf("truncate table %s cascade", tableName))
	}
}
