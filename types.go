package main

// ContainerInfo stores information about a running JupyterLab container
type ContainerInfo struct {
    ID       string
    HostPort int
    URL      string
}
