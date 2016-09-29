package models

type TaskAggregation struct {
	ProjectID           string `bson:"project_id"`
	ProjectExternalName string `bson:"project_ext_name"`
	ProjectExternalID   string `bson:"project_ext_id"`
	Name                string `bson:"task_name"`
	Minutes             int    `bson:"minutes"`
}
