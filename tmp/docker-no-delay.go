package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var dockerClient, _ = client.NewClientWithOpts(client.FromEnv)

const (
    minPort = 20000 // 🔹 Define min port range
    maxPort = 30000 // 🔹 Define max port range
)

// sanitizeUsername ensures the container name is valid for Docker
func sanitizeUsername(username string) string {
    sanitized := strings.ToLower(username)
    sanitized = strings.ReplaceAll(sanitized, "@", "-")
    sanitized = strings.ReplaceAll(sanitized, ".", "-")
    re := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
    sanitized = re.ReplaceAllString(sanitized, "")
    return sanitized
}

// findAvailablePort checks ports within a defined range
func findAvailablePort() int {
    rand.Seed(time.Now().UnixNano()) // Randomize seed for different users

    for i := 0; i < 50; i++ { // Try 50 times to find an open port
        port := minPort + rand.Intn(maxPort-minPort)
        addr := ":" + strconv.Itoa(port)
        listener, err := net.Listen("tcp", addr)
        if err == nil {
            listener.Close()
            return port // Return an available port
        }
    }

    log.Fatalf("No available ports found in range %d-%d", minPort, maxPort)
    return 0 // Should never reach here
}

// getOrSpawnContainer checks if a JupyterLab container for the user exists or creates a new one
func getOrSpawnContainer(userID string) *ContainerInfo {
    ctx := context.Background()
    containerName := "jupyter-" + sanitizeUsername(userID)
    containerPort := "8888/tcp"

    // 🔹 Step 1: Check if the container already exists
    filterArgs := filters.NewArgs()
    filterArgs.Add("name", containerName)

    containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
        All:     true,
        Filters: filterArgs,
    })
    if err != nil {
        log.Fatalf("Failed to list containers: %v", err)
    }

    if len(containers) > 0 {
        existingPort := int(containers[0].Ports[0].PublicPort) // ✅ Convert uint16 to int
        log.Printf("Found existing container for user %s: %s on port %d", userID, containerName, existingPort)
        return &ContainerInfo{
            ID:       containers[0].ID,
            HostPort: existingPort,
            URL:      "http://localhost:" + strconv.Itoa(existingPort),
        }
    }

    // 🔹 Step 2: Assign an available port within the defined range
    hostPort := findAvailablePort()

    portBindings := nat.PortMap{
        nat.Port(containerPort): []nat.PortBinding{
            {HostIP: "0.0.0.0", HostPort: strconv.Itoa(hostPort)}, // ✅ Convert int to string
        },
    }

    resp, err := dockerClient.ContainerCreate(ctx,
        &container.Config{
            Image: "jupyter/base-notebook",
            Env:   []string{"JUPYTER_ENABLE_LAB=yes"},
            Cmd:   []string{"start-notebook.sh", "--NotebookApp.token=''", "--NotebookApp.allow_origin='*'", "--NotebookApp.ip='0.0.0.0'", "--NotebookApp.port=8888"},
        },
        &container.HostConfig{
            PortBindings: portBindings,
        },
        nil, nil, containerName)

    if err != nil {
        log.Fatalf("Failed to create container: %v", err)
    }

    dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

    log.Printf("Started new container for user %s: %s on port %d", userID, containerName, hostPort)

    return &ContainerInfo{
        ID:       resp.ID,
        HostPort: hostPort, // ✅ Use int directly
        URL:      "http://localhost:" + strconv.Itoa(hostPort),
    }
}
