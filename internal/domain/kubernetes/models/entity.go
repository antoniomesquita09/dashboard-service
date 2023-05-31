package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ContainerModel struct {
	ContainerName string `json:"container_name"`
	MemoryUsage   string `json:"memory_usage"`
	CPUUsage      string `json:"cpu_usage"`
}

type PodModel struct {
	ID         primitive.ObjectID `json:"id,omitempty"`
	PodName    string             `json:"pod_name"`
	PodStatus  string             `json:"pod_status"`
	Containers []ContainerModel   `json:"containers"`
}
