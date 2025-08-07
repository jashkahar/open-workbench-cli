package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jashkahar/open-workbench-platform/internal/generator"
	"github.com/jashkahar/open-workbench-platform/internal/generator/docker"
	"github.com/jashkahar/open-workbench-platform/internal/generator/terraform"
	manifestPkg "github.com/jashkahar/open-workbench-platform/internal/manifest"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Generate deployment configuration from workbench.yaml",
	Long: `Generate deployment configuration from your workbench.yaml file.

This command supports multiple deployment targets:
- docker: Generate Docker Compose configuration for local development
- terraform: Generate Terraform configuration for cloud infrastructure

Interactive Mode (no target specified):
  om compose

Direct Mode (with target specified):
  om compose --target docker
  om compose --target terraform

Examples:
  # Interactive mode - prompts for target selection
  om compose

  # Direct mode with Docker target
  om compose --target docker

  # Direct mode with Terraform target
  om compose --target terraform

The generated configuration will be based on your workbench.yaml file and
the selected target. For Terraform generation, ensure you have configured
environments in your workbench.yaml file.`,
	RunE: runCompose,
}

// initComposeCommand registers the compose command with the root command
func initComposeCommand() {
	if rootCmd != nil {
		rootCmd.AddCommand(composeCmd)
	}

	// Add target flag
	composeCmd.Flags().String("target", "", "Deployment target (docker, terraform)")
	// Add environment flag for Terraform
	composeCmd.Flags().String("env", "", "Environment name (dev, staging, prod) - required for Terraform")
}

func runCompose(cmd *cobra.Command, args []string) error {
	// Find workbench.yaml
	workbenchPath := "workbench.yaml"
	if _, err := os.Stat(workbenchPath); os.IsNotExist(err) {
		return fmt.Errorf("workbench.yaml not found in current directory. Please run this command from your project root")
	}

	fmt.Println("📖 Loading workbench.yaml...")

	// Load and parse workbench.yaml
	manifest, err := loadWorkbenchManifest(workbenchPath)
	if err != nil {
		return fmt.Errorf("failed to load workbench.yaml: %w", err)
	}

	fmt.Printf("✅ Loaded project: %s\n", manifest.Metadata.Name)

	// Get target from flag or prompt user
	target, err := getTarget(cmd)
	if err != nil {
		return fmt.Errorf("failed to get target: %w", err)
	}

	// For Terraform, handle environment configuration
	if target == "terraform" {
		if err := handleTerraformEnvironment(cmd, manifest); err != nil {
			return fmt.Errorf("failed to configure environment: %w", err)
		}
	}

	// Create generator registry
	registry := generator.NewRegistry()

	// Register generators
	dockerGen := docker.NewGenerator()
	terraformGen := terraform.NewGenerator()

	if err := registry.Register(dockerGen); err != nil {
		return fmt.Errorf("failed to register Docker generator: %w", err)
	}

	if err := registry.Register(terraformGen); err != nil {
		return fmt.Errorf("failed to register Terraform generator: %w", err)
	}

	// Get the selected generator
	gen, err := registry.Get(target)
	if err != nil {
		return fmt.Errorf("failed to get generator '%s': %w", target, err)
	}

	fmt.Printf("🔧 Using %s generator: %s\n", target, gen.Description())

	// Generate configuration
	if err := gen.Generate(manifest); err != nil {
		return fmt.Errorf("failed to generate %s configuration: %w", target, err)
	}

	return nil
}

// handleTerraformEnvironment handles environment configuration for Terraform generation
func handleTerraformEnvironment(cmd *cobra.Command, manifest *manifestPkg.WorkbenchManifest) error {
	// Check if environments are already configured
	if len(manifest.Environments) > 0 {
		fmt.Println("✅ Environments already configured in workbench.yaml")
		return nil
	}

	// Get environment from flag or prompt user
	envName, err := getEnvironment(cmd)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	// Create default environment with all services
	manifest.Environments = map[string]manifestPkg.Environment{
		envName: {
			Provider: "aws",       // Default provider
			Region:   "us-east-1", // Default region
			Config: map[string]string{
				"services": strings.Join(getServiceNames(manifest.Services), ","),
			},
		},
	}

	// Save updated workbench.yaml
	if err := saveWorkbenchManifest(manifest, "."); err != nil {
		return fmt.Errorf("failed to save workbench.yaml: %w", err)
	}

	fmt.Printf("✅ Added '%s' environment to workbench.yaml\n", envName)
	return nil
}

// getEnvironment gets the environment name from flag or prompts user
func getEnvironment(cmd *cobra.Command) (string, error) {
	// Check if environment is provided via flag
	envName, err := cmd.Flags().GetString("env")
	if err != nil {
		return "", err
	}

	// If environment is provided, validate it
	if envName != "" {
		validEnvs := []string{"dev", "staging", "prod"}
		for _, valid := range validEnvs {
			if envName == valid {
				return envName, nil
			}
		}
		return "", fmt.Errorf("invalid environment '%s'. Valid environments are: %s", envName, strings.Join(validEnvs, ", "))
	}

	// Interactive mode - prompt user for environment
	var envChoice string
	prompt := &survey.Select{
		Message: "Which environment would you like to configure?",
		Options: []string{
			"dev - Development environment",
			"staging - Staging environment",
			"prod - Production environment",
		},
		Help: "Select the environment for Terraform configuration",
	}

	if err := survey.AskOne(prompt, &envChoice); err != nil {
		return "", fmt.Errorf("failed to get environment selection: %w", err)
	}

	// Extract environment from choice
	if strings.Contains(envChoice, "dev") {
		return "dev", nil
	} else if strings.Contains(envChoice, "staging") {
		return "staging", nil
	} else if strings.Contains(envChoice, "prod") {
		return "prod", nil
	}

	return "", fmt.Errorf("invalid environment selection")
}

// getServiceNames extracts service names from the services map
func getServiceNames(services map[string]manifestPkg.Service) []string {
	var names []string
	for name := range services {
		names = append(names, name)
	}
	return names
}

// getTarget gets the target from flag or prompts user
func getTarget(cmd *cobra.Command) (string, error) {
	// Check if target is provided via flag
	target, err := cmd.Flags().GetString("target")
	if err != nil {
		return "", err
	}

	// If target is provided, validate it
	if target != "" {
		validTargets := []string{"docker", "terraform"}
		for _, valid := range validTargets {
			if target == valid {
				return target, nil
			}
		}
		return "", fmt.Errorf("invalid target '%s'. Valid targets are: %s", target, strings.Join(validTargets, ", "))
	}

	// Interactive mode - prompt user for target
	var targetChoice string
	prompt := &survey.Select{
		Message: "Which target would you like to compose for?",
		Options: []string{
			"docker - Generate Docker Compose configuration for local development",
			"terraform - Generate Terraform configuration for cloud infrastructure",
		},
		Help: "Select the deployment target for your configuration",
	}

	if err := survey.AskOne(prompt, &targetChoice); err != nil {
		return "", fmt.Errorf("failed to get target selection: %w", err)
	}

	// Extract target from choice
	if strings.Contains(targetChoice, "docker") {
		return "docker", nil
	} else if strings.Contains(targetChoice, "terraform") {
		return "terraform", nil
	}

	return "", fmt.Errorf("invalid target selection")
}

// loadWorkbenchManifest loads and parses the workbench.yaml file
func loadWorkbenchManifest(path string) (*manifestPkg.WorkbenchManifest, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read workbench.yaml: %w", err)
	}

	var manifest manifestPkg.WorkbenchManifest
	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse workbench.yaml: %w", err)
	}

	return &manifest, nil
}
