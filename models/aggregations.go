package models

type TaskAggregation struct {
	ProjectID string `bson:"project_id"`
	Name      string `bson:"task_name"`
	Minutes   int    `bson:"minutes"`
}
