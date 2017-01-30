package models

type TaskAggregation struct {
	TaskHash            string `bson:"task_hash"`
	ProjectExternalName string `bson:"project_ext_name"`
	ProjectExternalID   string `bson:"project_ext_id"`
	Name                string `bson:"task_name"`
	Minutes             int    `bson:"minutes"`
}

type UserStatisticsAggregation struct {
	Day		int	 `json:"day" bson:"day"`
	Minutes		int	 `json:"minutes" bson:"minutes"`
	ProjectsNames	[]string `json:"projects_names" bson:"projects_names"`
}
