package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jashkahar/open-workbench-platform/internal/manifest"
)

// Generator implements the Generator interface for Terraform
type Generator struct{}

// NewGenerator creates a new Terraform generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Name returns the unique identifier for this generator
func (g *Generator) Name() string {
	return "terraform"
}

// Description returns a human-readable description of this generator
func (g *Generator) Description() string {
	return "Generate Terraform configuration for cloud infrastructure"
}

// Validate checks if the manifest is compatible with this generator
func (g *Generator) Validate(manifest *manifest.WorkbenchManifest) error {
	if manifest == nil {
		return fmt.Errorf("manifest cannot be nil")
	}

	if manifest.Metadata.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if len(manifest.Services) == 0 {
		return fmt.Errorf("at least one service is required")
	}

	// Check if environments are configured
	if len(manifest.Environments) == 0 {
		return fmt.Errorf("at least one environment must be configured for Terraform generation")
	}

	return nil
}

// Generate creates the Terraform configuration for the given manifest
func (g *Generator) Generate(manifest *manifest.WorkbenchManifest) error {
	// Validate the manifest
	if err := g.Validate(manifest); err != nil {
		return fmt.Errorf("manifest validation failed: %w", err)
	}

	fmt.Println("🔧 Generating Terraform configuration...")

	// Create terraform directory
	terraformDir := "terraform"
	if err := os.MkdirAll(terraformDir, 0755); err != nil {
		return fmt.Errorf("failed to create terraform directory: %w", err)
	}

	// Generate main.tf
	if err := g.generateMainTf(manifest, terraformDir); err != nil {
		return fmt.Errorf("failed to generate main.tf: %w", err)
	}

	// Generate variables.tf
	if err := g.generateVariablesTf(manifest, terraformDir); err != nil {
		return fmt.Errorf("failed to generate variables.tf: %w", err)
	}

	// Generate outputs.tf
	if err := g.generateOutputsTf(manifest, terraformDir); err != nil {
		return fmt.Errorf("failed to generate outputs.tf: %w", err)
	}

	// Generate terraform.tfvars.example
	if err := g.generateTfvarsExample(manifest, terraformDir); err != nil {
		return fmt.Errorf("failed to generate terraform.tfvars.example: %w", err)
	}

	fmt.Println("✅ Generated Terraform configuration")

	// Print success message with instructions
	printTerraformSuccessMessage()

	return nil
}

func (g *Generator) generateMainTf(manifest *manifest.WorkbenchManifest, terraformDir string) error {
	content := `# Terraform configuration for ` + manifest.Metadata.Name + `

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC and networking
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.project_name}-vpc"
  }
}

resource "aws_subnet" "public" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = var.public_subnet_cidr
  availability_zone = var.availability_zone

  tags = {
    Name = "${var.project_name}-public-subnet"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project_name}-igw"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${var.project_name}-public-rt"
  }
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

# Security groups
resource "aws_security_group" "app" {
  name_prefix = "${var.project_name}-app-"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-app-sg"
  }
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.project_name}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "${var.project_name}-cluster"
  }
}

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "${var.project_name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.app.id]
  subnets            = [aws_subnet.public.id]

  tags = {
    Name = "${var.project_name}-alb"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# Services
`

	// Add service-specific resources
	for serviceName, service := range manifest.Services {
		content += g.generateServiceResources(serviceName, service)
	}

	// Add component-specific resources
	for componentName, component := range manifest.Components {
		content += g.generateComponentResources(componentName, component)
	}

	return os.WriteFile(filepath.Join(terraformDir, "main.tf"), []byte(content), 0644)
}

func (g *Generator) generateServiceResources(serviceName string, service manifest.Service) string {
	// Generate ECS service
	ecsService := fmt.Sprintf(`
# Service: %s
resource "aws_ecs_service" "%s" {
  name            = "%s"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.%s.arn
  desired_count   = var.%s_desired_count

  network_configuration {
    subnets         = [aws_subnet.public.id]
    security_groups = [aws_security_group.app.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.%s.arn
    container_name   = "%s"
    container_port   = %d
  }

  depends_on = [aws_lb_listener.http]

  tags = {
    Name = "%s"
  }
}
`, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, service.Port, serviceName)

	// Generate ECS task definition
	taskDefinition := fmt.Sprintf(`
resource "aws_ecs_task_definition" "%s" {
  family                   = "%s"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                     = var.%s_cpu
  memory                  = var.%s_memory

  container_definitions = jsonencode([
    {
      name  = "%s"
      image = var.%s_image
      portMappings = [
        {
          containerPort = %d
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "NODE_ENV"
          value = "production"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/%s"
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = {
    Name = "%s"
  }
}
`, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, service.Port, serviceName, serviceName)

	// Generate load balancer target group
	targetGroup := fmt.Sprintf(`
resource "aws_lb_target_group" "%s" {
  name     = "%s-tg"
  port     = %d
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }

  tags = {
    Name = "%s-tg"
  }
}
`, serviceName, serviceName, service.Port, serviceName)

	return ecsService + taskDefinition + targetGroup
}

func (g *Generator) generateComponentResources(componentName string, component manifest.Component) string {
	content := fmt.Sprintf(`
# Component: %s
resource "aws_ecs_service" "%s" {
  name            = "%s"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.%s.arn
  desired_count   = var.%s_desired_count

  network_configuration {
    subnets         = [aws_subnet.public.id]
    security_groups = [aws_security_group.app.id]
  }

  tags = {
    Name = "%s"
  }
}

resource "aws_ecs_task_definition" "%s" {
  family                   = "%s"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                     = var.%s_cpu
  memory                  = var.%s_memory

  container_definitions = jsonencode([
    {
      name  = "%s"
      image = var.%s_image
      portMappings = [
        {
          containerPort = 80
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "NODE_ENV"
          value = "production"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/%s"
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = {
    Name = "%s"
  }
}

`, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName)

	return content
}

func (g *Generator) generateVariablesTf(manifest *manifest.WorkbenchManifest, terraformDir string) error {
	content := `# Variables for ` + manifest.Metadata.Name + `

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "` + manifest.Metadata.Name + `"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidr" {
  description = "CIDR block for public subnet"
  type        = string
  default     = "10.0.1.0/24"
}

variable "availability_zone" {
  description = "Availability zone"
  type        = string
  default     = "us-east-1a"
}

`

	// Add variables for each service
	for serviceName := range manifest.Services {
		content += fmt.Sprintf(`
variable "%s_desired_count" {
  description = "Desired count for %s service"
  type        = number
  default     = 1
}

variable "%s_cpu" {
  description = "CPU units for %s service"
  type        = number
  default     = 256
}

variable "%s_memory" {
  description = "Memory for %s service"
  type        = number
  default     = 512
}

variable "%s_image" {
  description = "Docker image for %s service"
  type        = string
  default     = "nginx:alpine"
}

`, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName)
	}

	// Add variables for each component
	for componentName := range manifest.Components {
		content += fmt.Sprintf(`
variable "%s_desired_count" {
  description = "Desired count for %s component"
  type        = number
  default     = 1
}

variable "%s_cpu" {
  description = "CPU units for %s component"
  type        = number
  default     = 256
}

variable "%s_memory" {
  description = "Memory for %s component"
  type        = number
  default     = 512
}

variable "%s_image" {
  description = "Docker image for %s component"
  type        = string
  default     = "nginx:alpine"
}

`, componentName, componentName, componentName, componentName, componentName, componentName, componentName, componentName)
	}

	return os.WriteFile(filepath.Join(terraformDir, "variables.tf"), []byte(content), 0644)
}

func (g *Generator) generateOutputsTf(manifest *manifest.WorkbenchManifest, terraformDir string) error {
	content := `# Outputs for ` + manifest.Metadata.Name + `

output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS name"
  value       = aws_lb.main.dns_name
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

`

	// Add outputs for each service
	for serviceName := range manifest.Services {
		content += fmt.Sprintf(`
output "%s_service_name" {
  description = "%s service name"
  value       = aws_ecs_service.%s.name
}

output "%s_task_definition_arn" {
  description = "%s task definition ARN"
  value       = aws_ecs_task_definition.%s.arn
}

`, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName)
	}

	return os.WriteFile(filepath.Join(terraformDir, "outputs.tf"), []byte(content), 0644)
}

func (g *Generator) generateTfvarsExample(manifest *manifest.WorkbenchManifest, terraformDir string) error {
	content := `# Example terraform.tfvars for ` + manifest.Metadata.Name + `

aws_region = "us-east-1"
project_name = "` + manifest.Metadata.Name + `"
vpc_cidr = "10.0.0.0/16"
public_subnet_cidr = "10.0.1.0/24"
availability_zone = "us-east-1a"

`

	// Add example values for each service
	for serviceName := range manifest.Services {
		content += fmt.Sprintf(`
# %s service configuration
%s_desired_count = 1
%s_cpu = 256
%s_memory = 512
%s_image = "nginx:alpine"

`, serviceName, serviceName, serviceName, serviceName, serviceName)
	}

	// Add example values for each component
	for componentName := range manifest.Components {
		content += fmt.Sprintf(`
# %s component configuration
%s_desired_count = 1
%s_cpu = 256
%s_memory = 512
%s_image = "nginx:alpine"

`, componentName, componentName, componentName, componentName, componentName)
	}

	return os.WriteFile(filepath.Join(terraformDir, "terraform.tfvars.example"), []byte(content), 0644)
}

func printTerraformSuccessMessage() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Successfully generated Terraform configuration!")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n📁 Generated files:")
	fmt.Println("  • terraform/main.tf - Main Terraform configuration")
	fmt.Println("  • terraform/variables.tf - Variable definitions")
	fmt.Println("  • terraform/outputs.tf - Output definitions")
	fmt.Println("  • terraform/terraform.tfvars.example - Example variable values")

	fmt.Println("\n🚀 To deploy your infrastructure:")
	fmt.Println("  1. cd terraform")
	fmt.Println("  2. terraform init")
	fmt.Println("  3. cp terraform.tfvars.example terraform.tfvars")
	fmt.Println("  4. Edit terraform.tfvars with your values")
	fmt.Println("  5. terraform plan")
	fmt.Println("  6. terraform apply")

	fmt.Println("\n📋 Additional commands:")
	fmt.Println("  terraform destroy     # Destroy all resources")
	fmt.Println("  terraform output      # View outputs")
	fmt.Println("  terraform state list  # List all resources")

	fmt.Println("\n💡 Tips:")
	fmt.Println("  • Review and customize the generated Terraform files")
	fmt.Println("  • Set up AWS credentials before running terraform")
	fmt.Println("  • Use terraform.tfvars for environment-specific values")
	fmt.Println("  • Consider using remote state storage for team collaboration")

	fmt.Println("\n🎉 Your cloud infrastructure configuration is ready!")
}
