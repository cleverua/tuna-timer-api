package models

type TaskAggregation struct {
	TaskHash            string `bson:"task_hash"`
	ProjectExternalName string `bson:"project_ext_name"`
	ProjectExternalID   string `bson:"project_ext_id"`
	Name                string `bson:"task_name"`
	Minutes             int    `bson:"minutes"`
}

type UserReportAggregation struct {
	Day		int8	 `bson:"day"`
	Projects	[]string `bson:"projects_names"`
	total		int	 `bson:"total"`
}
