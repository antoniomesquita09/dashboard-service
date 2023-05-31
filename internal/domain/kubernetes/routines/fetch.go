package routines

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"dashboard-service/internal/domain/kubernetes/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func MakeKubernetesRoutine(seconds int64) {
	for {
		fetchKubernetesMetrics()
		// Wait for 5 seconds before making the next API call
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func fetchKubernetesMetrics() {
	// Get the kubeconfig file path
	homeDir := homedir.HomeDir()
	kubeconfig := filepath.Join(homeDir, ".kube", "config")

	// Build the Kubernetes client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Error building kubeconfig: %v", err)
	}

	// Create the Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating Kubernetes client: %v", err)
	}

	// Set the namespace you want to monitor
	namespace := "default"

	// Get the list of pods in the specified namespace
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error getting pod list: %v", err)
	}

	// Create the Kubernetes Metrics client
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes Metrics client: %v", err)
	}

	clusterMetrics := []models.PodMetric{}

	// Retrieve and print memory usage and CPU information for each pod
	for _, pod := range podList.Items {
		containerMetrics := []models.ContainerMetric{}

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

			containerMetric := models.ContainerMetric{
				ContainerName: containerName,
				MemoryUsage:   memoryUsage.String(),
				CPUUsage:      cpuUsage.String(),
			}

			containerMetrics = append(containerMetrics, containerMetric)

		}

		podMetric := models.PodMetric{
			PodName:    pod.Name,
			PodSatus:   string(pod.Status.Phase),
			Containers: containerMetrics,
		}

		clusterMetrics = append(clusterMetrics, podMetric)
	}
}
