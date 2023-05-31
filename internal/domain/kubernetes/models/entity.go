package models

type ContainerMetric struct {
	ContainerName string `json:"container_name"`
	MemoryUsage   string `json:"memory_usage"`
	CPUUsage      string `json:"cpu_usage"`
}

type PodMetric struct {
	PodName    string            `json:"pod_name"`
	PodSatus   string            `json:"pod_status"`
	Containers []ContainerMetric `json:"containers"`
}
