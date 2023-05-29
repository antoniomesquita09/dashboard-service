package handlers

import (
	"context"
	"log"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Error struct {
	Message string `json:"message"`
}

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

func FetchKubernetesMetrics(c echo.Context) error {
	// Get the kubeconfig file path
	homeDir := homedir.HomeDir()
	kubeconfig := filepath.Join(homeDir, ".kube", "config")

	// Build the Kubernetes client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Error building kubeconfig: %v", err)
		error := Error{
			Message: "Error building kubeconfig",
		}
		return c.JSON(http.StatusInternalServerError, error)
	}

	// Create the Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating Kubernetes client: %v", err)
		error := Error{
			Message: "Error creating Kubernetes client",
		}
		return c.JSON(http.StatusInternalServerError, error)
	}

	// Set the namespace you want to monitor
	namespace := "default"

	// Get the list of pods in the specified namespace
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error getting pod list: %v", err)
		error := Error{
			Message: "Error getting pod list",
		}
		return c.JSON(http.StatusNotFound, error)
	}

	// Create the Kubernetes Metrics client
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes Metrics client: %v", err)
		error := Error{
			Message: "Error creating Kubernetes Metrics client",
		}
		return c.JSON(http.StatusInternalServerError, error)
	}

	clusterMetrics := []PodMetric{}

	// Retrieve and print memory usage and CPU information for each pod
	for _, pod := range podList.Items {
		containerMetrics := []ContainerMetric{}

		// Retrieve pod metrics
		podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), pod.Name, v1.GetOptions{})
		if err != nil {
			log.Printf("Error getting metrics for pod %s: %v\n", pod.Name, err)
			continue
		}

		for _, container := range podMetrics.Containers {
			containerName := container.Name
			memoryUsage := container.Usage["memory"]
			cpuUsage := container.Usage["cpu"]

			containerMetric := ContainerMetric{
				ContainerName: containerName,
				MemoryUsage:   memoryUsage.String(),
				CPUUsage:      cpuUsage.String(),
			}

			containerMetrics = append(containerMetrics, containerMetric)

		}

		podMetric := PodMetric{
			PodName:    pod.Name,
			PodSatus:   string(pod.Status.Phase),
			Containers: containerMetrics,
		}

		clusterMetrics = append(clusterMetrics, podMetric)
	}

	return c.JSON(http.StatusOK, clusterMetrics)
}
