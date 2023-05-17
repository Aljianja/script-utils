package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Image    string            `yaml:"image"`
	Labels   map[string]string `yaml:"labels"`
	Networks []string          `yaml:"networks"`
}

type DockerCompose struct {
	Version  string              `yaml:"version"`
	Services map[string]*Service `yaml:"services"`
	Networks map[string]struct{} `yaml:"networks"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to a docker-compose.yml file as an argument.")
	}

	sourceFilePath := os.Args[1]

	dc, err := parseDockerComposeFile(sourceFilePath)
	if err != nil {
		log.Fatalf("Error reading docker-compose file: %v", err)
	}

	err = updateServicesWithTraefikLabels(dc, reader)
	if err != nil {
		log.Fatalf("Error updating services: %v", err)
	}

	err = writeDockerComposeFile("docker-compose-updated.yml", dc)
	if err != nil {
		log.Fatalf("Error writing docker-compose file: %v", err)
	}
}

func parseDockerComposeFile(filePath string) (*DockerCompose, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	dc := &DockerCompose{}
	err = yaml.Unmarshal(file, dc)
	if err != nil {
		return nil, err
	}

	return dc, nil
}

func updateServicesWithTraefikLabels(dc *DockerCompose, reader *bufio.Reader) error {
	// Add Traefik network if not exists
	if dc.Networks == nil {
		dc.Networks = make(map[string]struct{})
	}

	if _, exists := dc.Networks["traefik"]; !exists {
		dc.Networks["traefik"] = struct{}{}
	}

	for name, service := range dc.Services {
		shouldUpdate, err := promptForConfirmation(fmt.Sprintf("Add Traefik labels to service '%s'? (Y/n): ", name), reader)
		if err != nil {
			return err
		}

		if shouldUpdate {
			rule, err := promptForInput("Enter the Traefik rule: ", reader)
			if err != nil {
				return err
			}

			entrypoint, err := promptForInput("Enter the Traefik entrypoint: ", reader)
			if err != nil {
				return err
			}

			if service.Labels == nil {
				service.Labels = make(map[string]string)
			}

			service.Labels["traefik.http.routers."+name+".rule"] = rule
			service.Labels["traefik.http.routers."+name+".entrypoints"] = entrypoint

			// Add Traefik network to the service
			service.Networks = append(service.Networks, "traefik")
		}
	}

	return nil
}

func writeDockerComposeFile(filePath string, dc *DockerCompose) error {
	out, err := yaml.Marshal(dc)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, out, 0644)
	if err != nil {
		return err
	}

	return nil
}

func promptForInput(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(text, "\n"), nil
}

func promptForConfirmation(prompt string, reader *bufio.Reader) (bool, error) {
	text, err := promptForInput(prompt, reader)
	if err != nil {
		return false, err
	}

	text = strings.ToLower(text)
	if text == "y" || text == "yes" {
		return true, nil
	} else if text == "n" || text == "no" {
		return false, nil
	} else {
		return false, errors.New("Invalid input")
	}
}
