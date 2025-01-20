package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

type Container struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Health      string    `json:"health"`
	CPU         float64   `json:"cpu_usage"`
	Memory      int64     `json:"memory_usage"`
	RestartCount int      `json:"restart_count"`
}

type ContainerRequest struct {
	Name  string            `json:"name"`
	Image string            `json:"image"`
	Env   map[string]string `json:"env"`
	Ports map[string]string `json:"ports"`
}

type Orchestrator struct {
	client     *client.Client
	containers map[string]*Container
	mutex      sync.RWMutex
}

func NewOrchestrator() (*Orchestrator, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Orchestrator{
		client:     cli,
		containers: make(map[string]*Container),
	}, nil
}

func (o *Orchestrator) CreateContainer(w http.ResponseWriter, r *http.Request) {
	var req ContainerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare environment variables
	var env []string
	for k, v := range req.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Prepare port bindings
	portBindings := make(map[string][]types.PortBinding)
	exposedPorts := make(map[string]struct{})
	for containerPort, hostPort := range req.Ports {
		portBindings[containerPort] = []types.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		}
		exposedPorts[containerPort] = struct{}{}
	}

	// Create container
	resp, err := o.client.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        req.Image,
			Env:          env,
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			PortBindings: portBindings,
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		nil,
		nil,
		req.Name,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Start container
	if err := o.client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	container := &Container{
		ID:        resp.ID,
		Name:      req.Name,
		Image:     req.Image,
		Status:    "running",
		CreatedAt: time.Now(),
	}

	o.mutex.Lock()
	o.containers[resp.ID] = container
	o.mutex.Unlock()

	json.NewEncoder(w).Encode(container)
}

func (o *Orchestrator) ListContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := o.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []Container
	for _, c := range containers {
		stats, err := o.getContainerStats(c.ID)
		if err != nil {
			log.Printf("Error getting stats for container %s: %v", c.ID, err)
			continue
		}

		container := Container{
			ID:        c.ID,
			Name:      c.Names[0],
			Image:     c.Image,
			Status:    c.Status,
			CreatedAt: time.Unix(c.Created, 0),
			CPU:       stats.CPU,
			Memory:    stats.Memory,
		}
		result = append(result, container)
	}

	json.NewEncoder(w).Encode(result)
}

func (o *Orchestrator) StopContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	timeout := 10 * time.Second
	if err := o.client.ContainerStop(context.Background(), id, &timeout); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (o *Orchestrator) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := o.client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	o.mutex.Lock()
	delete(o.containers, id)
	o.mutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

type ContainerStats struct {
	CPU    float64
	Memory int64
}

func (o *Orchestrator) getContainerStats(id string) (*ContainerStats, error) {
	stats, err := o.client.ContainerStats(context.Background(), id, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var statsJSON types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
		return nil, err
	}

	cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - statsJSON.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)
	cpuPercent := 0.0
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	return &ContainerStats{
		CPU:    cpuPercent,
		Memory: int64(statsJSON.MemoryStats.Usage),
	}, nil
}

func (o *Orchestrator) monitorContainers() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		containers, err := o.client.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			log.Printf("Error listing containers: %v", err)
			continue
		}

		for _, c := range containers {
			stats, err := o.getContainerStats(c.ID)
			if err != nil {
				log.Printf("Error getting stats for container %s: %v", c.ID, err)
				continue
			}

			o.mutex.Lock()
			if container, exists := o.containers[c.ID]; exists {
				container.Status = c.Status
				container.CPU = stats.CPU
				container.Memory = stats.Memory
			}
			o.mutex.Unlock()
		}
	}
}

func main() {
	orchestrator, err := NewOrchestrator()
	if err != nil {
		log.Fatal(err)
	}

	// Start container monitoring
	go orchestrator.monitorContainers()

	r := mux.NewRouter()
	r.HandleFunc("/containers", orchestrator.CreateContainer).Methods("POST")
	r.HandleFunc("/containers", orchestrator.ListContainers).Methods("GET")
	r.HandleFunc("/containers/{id}", orchestrator.StopContainer).Methods("POST")
	r.HandleFunc("/containers/{id}", orchestrator.RemoveContainer).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Container orchestrator starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
