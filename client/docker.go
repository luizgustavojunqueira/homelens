package client

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"

	"homelens/shared"
)

const socketPath = "/var/run/docker.sock"

type DockerContainerRead struct {
	Names  []string
	State  string
	Image  string
	Status string
	Ports  []struct {
		PrivatePort int
		PublicPort  int
		Type        string
	}
}

func readDockerContainers() []shared.DockerContainer {
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return nil
	}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	resp, err := httpc.Get("http://localhost/containers/json")
	if err != nil {
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var rawContainers []DockerContainerRead
	if err := json.NewDecoder(resp.Body).Decode(&rawContainers); err != nil {
		return nil
	}

	var cleanContainers []shared.DockerContainer
	for _, raw := range rawContainers {
		cleanName := ""

		if len(raw.Names) > 0 {
			cleanName = strings.TrimPrefix(raw.Names[0], "/")
		}

		var cleanPorts []shared.DockerPort
		for _, p := range raw.Ports {
			cleanPorts = append(cleanPorts, shared.DockerPort{
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			})
		}

		cleanContainers = append(cleanContainers, shared.DockerContainer{
			Name:   cleanName,
			State:  raw.State,
			Image:  raw.Image,
			Status: raw.Status,
			Ports:  cleanPorts,
		})
	}

	return cleanContainers
}
