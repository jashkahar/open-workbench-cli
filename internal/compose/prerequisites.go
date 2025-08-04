package compose

import (
	"fmt"
	"os/exec"
	"runtime"
)

// PrerequisiteChecker handles validation of required external tools
type PrerequisiteChecker struct{}

// NewPrerequisiteChecker creates a new prerequisite checker
func NewPrerequisiteChecker() *PrerequisiteChecker {
	return &PrerequisiteChecker{}
}

// CheckDocker checks if Docker is installed and available
func (pc *PrerequisiteChecker) CheckDocker() error {
	_, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("Docker is not installed or not available in PATH. Please install Docker from https://docs.docker.com/get-docker/")
	}
	return nil
}

// CheckDockerCompose checks if Docker Compose is installed and available
func (pc *PrerequisiteChecker) CheckDockerCompose() error {
	// Try docker-compose first (legacy)
	_, err := exec.LookPath("docker-compose")
	if err == nil {
		return nil
	}

	// Try docker compose (newer version)
	_, err = exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("Docker is not installed or not available in PATH")
	}

	// Check if docker compose plugin is available
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Docker Compose is not available. Please install Docker Compose or ensure you have Docker Desktop with Compose support")
	}

	return nil
}

// CheckAllPrerequisites checks all required prerequisites
func (pc *PrerequisiteChecker) CheckAllPrerequisites() error {
	if err := pc.CheckDocker(); err != nil {
		return err
	}

	if err := pc.CheckDockerCompose(); err != nil {
		return err
	}

	return nil
}

// GetDockerComposeCommand returns the appropriate docker-compose command for the system
func (pc *PrerequisiteChecker) GetDockerComposeCommand() string {
	// Check if docker-compose is available
	_, err := exec.LookPath("docker-compose")
	if err == nil {
		return "docker-compose"
	}

	// Use docker compose (newer version)
	return "docker compose"
}

// GetPlatformSpecificInstructions returns platform-specific installation instructions
func (pc *PrerequisiteChecker) GetPlatformSpecificInstructions() string {
	switch runtime.GOOS {
	case "windows":
		return `To install Docker on Windows:
1. Download Docker Desktop from https://docs.docker.com/desktop/install/windows/
2. Install and start Docker Desktop
3. Ensure Docker Desktop is running before using 'om compose'`
	case "darwin":
		return `To install Docker on macOS:
1. Download Docker Desktop from https://docs.docker.com/desktop/install/mac/
2. Install and start Docker Desktop
3. Ensure Docker Desktop is running before using 'om compose'`
	default:
		return `To install Docker on Linux:
1. Follow the instructions at https://docs.docker.com/engine/install/
2. Install Docker Compose: https://docs.docker.com/compose/install/
3. Ensure Docker service is running before using 'om compose'`
	}
}
