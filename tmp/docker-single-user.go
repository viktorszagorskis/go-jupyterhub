package main

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)



var dockerClient, _ = client.NewClientWithOpts(client.FromEnv)

// sanitizeUsername ensures the container name is valid for Docker
func sanitizeUsername(username string) string {
    sanitized := strings.ToLower(username)
    sanitized = strings.ReplaceAll(sanitized, "@", "-")
    sanitized = strings.ReplaceAll(sanitized, ".", "-")
    re := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
    sanitized = re.ReplaceAllString(sanitized, "")
    return sanitized
}

// getOrSpawnContainer checks if a JupyterLab container for the user exists or creates a new one
func getOrSpawnContainer(userID string) *ContainerInfo {
    ctx := context.Background()
    containerName := "jupyter-" + sanitizeUsername(userID)
    containerPort := "8888/tcp"
    hostPort := "18888"

    // 🔹 Step 1: Check if the container already exists
    containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{All: true})
    if err != nil {
        log.Fatalf("Failed to list containers: %v", err)
    }

    for _, c := range containers {
        for _, name := range c.Names {
            if name == "/"+containerName {
                log.Printf("Found existing container for user %s: %s", userID, containerName)
                return &ContainerInfo{
                    ID:       c.ID,
                    HostPort: 18888,
                    URL:      "http://localhost:18888",
                }
            }
        }
    }

    // 🔹 Step 3: If no existing container, create a new one
    portBindings := nat.PortMap{
        nat.Port(containerPort): []nat.PortBinding{
            {HostIP: "0.0.0.0", HostPort: hostPort},
        },
    }

    resp, err := dockerClient.ContainerCreate(ctx,
        &container.Config{
            Image: "jupyter/base-notebook",
            Env:   []string{"JUPYTER_ENABLE_LAB=yes"},
            Cmd:   []string{"start-notebook.sh", "--NotebookApp.token=''", "--NotebookApp.allow_origin='*'", "--NotebookApp.ip='0.0.0.0'", "--NotebookApp.port=8888"},
        },
        &container.HostConfig{
            PortBindings: portBindings, // ✅ Correctly map ports
        },
        nil, nil, containerName)

    if err != nil {
        log.Fatalf("Failed to create container: %v", err)
    }

    dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})


    log.Printf("Started new container for user %s: %s", userID, containerName)

    return &ContainerInfo{
        ID:       resp.ID,
        HostPort: 18888,
        URL:      "http://localhost:18888",
    }
}
