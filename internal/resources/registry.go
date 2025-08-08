package resources

import (
	"fmt"
	"sync"
)

// Registry manages all available resource blueprints
type Registry struct {
	blueprints map[string]ResourceBlueprint
	mutex      sync.RWMutex
}

// NewRegistry creates a new resource registry
func NewRegistry() *Registry {
	registry := &Registry{
		blueprints: make(map[string]ResourceBlueprint),
	}

	// Register default resource blueprints
	registry.registerDefaultBlueprints()

	return registry
}

// Get retrieves a resource blueprint by name
func (r *Registry) Get(name string) (ResourceBlueprint, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	blueprint, exists := r.blueprints[name]
	if !exists {
		return ResourceBlueprint{}, fmt.Errorf("resource blueprint '%s' not found", name)
	}

	return blueprint, nil
}

// List returns all available resource blueprints
func (r *Registry) List() []ResourceBlueprint {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	blueprints := make([]ResourceBlueprint, 0, len(r.blueprints))
	for _, blueprint := range r.blueprints {
		blueprints = append(blueprints, blueprint)
	}

	return blueprints
}

// ListByCategory returns resource blueprints filtered by category
func (r *Registry) ListByCategory(category string) []ResourceBlueprint {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var blueprints []ResourceBlueprint
	for _, blueprint := range r.blueprints {
		if blueprint.Category == category {
			blueprints = append(blueprints, blueprint)
		}
	}

	return blueprints
}

// Names returns all available resource blueprint names
func (r *Registry) Names() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.blueprints))
	for name := range r.blueprints {
		names = append(names, name)
	}

	return names
}

// Categories returns all available resource categories
func (r *Registry) Categories() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	categories := make(map[string]bool)
	for _, blueprint := range r.blueprints {
		categories[blueprint.Category] = true
	}

	result := make([]string, 0, len(categories))
	for category := range categories {
		result = append(result, category)
	}

	return result
}

// registerDefaultBlueprints registers the default resource blueprints
func (r *Registry) registerDefaultBlueprints() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Database resources
	r.blueprints["postgres-db"] = ResourceBlueprint{
		Name:        "postgres-db",
		Description: "A PostgreSQL Database",
		Category:    "database",
		DockerComposeSnippet: `
    image: postgres:{{.Version}}
    environment:
      - POSTGRES_DB={{.DatabaseName}}
      - POSTGRES_USER={{.Username}}
      - POSTGRES_PASSWORD={{.Password}}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "{{.Port}}:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U {{.Username}} -d {{.DatabaseName}}"]
      interval: 10s
      timeout: 5s
      retries: 5`,
		TerraformModule: "modules/aws/rds-postgres",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "PostgreSQL version", Type: "select", Required: true, Default: "15", Options: []string{"13", "14", "15", "16"}},
			{Name: "databaseName", Description: "Database name", Type: "string", Required: true, Default: "app"},
			{Name: "username", Description: "Database username", Type: "string", Required: true, Default: "postgres"},
			{Name: "password", Description: "Database password", Type: "string", Required: true},
			{Name: "port", Description: "Database port", Type: "number", Required: false, Default: 5432},
		},
	}

	r.blueprints["mysql-db"] = ResourceBlueprint{
		Name:        "mysql-db",
		Description: "A MySQL Database",
		Category:    "database",
		DockerComposeSnippet: `
    image: mysql:{{.Version}}
    environment:
      - MYSQL_DATABASE={{.DatabaseName}}
      - MYSQL_USER={{.Username}}
      - MYSQL_PASSWORD={{.Password}}
      - MYSQL_ROOT_PASSWORD={{.RootPassword}}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "{{.Port}}:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5`,
		TerraformModule: "modules/aws/rds-mysql",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "MySQL version", Type: "select", Required: true, Default: "8.0", Options: []string{"5.7", "8.0"}},
			{Name: "databaseName", Description: "Database name", Type: "string", Required: true, Default: "app"},
			{Name: "username", Description: "Database username", Type: "string", Required: true, Default: "app"},
			{Name: "password", Description: "Database password", Type: "string", Required: true},
			{Name: "rootPassword", Description: "Root password", Type: "string", Required: true},
			{Name: "port", Description: "Database port", Type: "number", Required: false, Default: 3306},
		},
	}

	r.blueprints["mongodb"] = ResourceBlueprint{
		Name:        "mongodb",
		Description: "A MongoDB Database",
		Category:    "database",
		DockerComposeSnippet: `
    image: mongo:{{.Version}}
    environment:
      - MONGO_INITDB_DATABASE={{.DatabaseName}}
      - MONGO_INITDB_ROOT_USERNAME={{.Username}}
      - MONGO_INITDB_ROOT_PASSWORD={{.Password}}
    volumes:
      - mongodb_data:/data/db
    ports:
      - "{{.Port}}:27017"
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5`,
		TerraformModule: "modules/aws/documentdb",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "MongoDB version", Type: "select", Required: true, Default: "7.0", Options: []string{"6.0", "7.0"}},
			{Name: "databaseName", Description: "Database name", Type: "string", Required: true, Default: "app"},
			{Name: "username", Description: "Database username", Type: "string", Required: true, Default: "admin"},
			{Name: "password", Description: "Database password", Type: "string", Required: true},
			{Name: "port", Description: "Database port", Type: "number", Required: false, Default: 27017},
		},
	}

	// Cache resources
	r.blueprints["redis-cache"] = ResourceBlueprint{
		Name:        "redis-cache",
		Description: "A Redis Cache",
		Category:    "cache",
		DockerComposeSnippet: `
    image: redis:{{.Version}}
    command: redis-server --requirepass {{.Password}}
    volumes:
      - redis_data:/data
    ports:
      - "{{.Port}}:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5`,
		TerraformModule: "modules/aws/elasticache-redis",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "Redis version", Type: "select", Required: true, Default: "7.2", Options: []string{"6.2", "7.0", "7.2"}},
			{Name: "password", Description: "Redis password", Type: "string", Required: true},
			{Name: "port", Description: "Redis port", Type: "number", Required: false, Default: 6379},
		},
	}

	r.blueprints["memcached"] = ResourceBlueprint{
		Name:        "memcached",
		Description: "A Memcached Cache",
		Category:    "cache",
		DockerComposeSnippet: `
    image: memcached:{{.Version}}
    ports:
      - "{{.Port}}:11211"
    healthcheck:
      test: ["CMD", "memcached-tool", "localhost:11211", "stats"]
      interval: 10s
      timeout: 5s
      retries: 5`,
		TerraformModule: "modules/aws/elasticache-memcached",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "Memcached version", Type: "select", Required: true, Default: "1.6", Options: []string{"1.6"}},
			{Name: "port", Description: "Memcached port", Type: "number", Required: false, Default: 11211},
		},
	}

	// Storage resources
	r.blueprints["s3-bucket"] = ResourceBlueprint{
		Name:        "s3-bucket",
		Description: "An AWS S3 Bucket",
		Category:    "storage",
		DockerComposeSnippet: `
    image: minio/minio:{{.Version}}
    command: server /data --console-address ":9001"
    environment:
      - MINIO_ROOT_USER={{.AccessKey}}
      - MINIO_ROOT_PASSWORD={{.SecretKey}}
    volumes:
      - minio_data:/data
    ports:
      - "{{.Port}}:9000"
      - "9001:9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 10s
      timeout: 5s
      retries: 5`,
		TerraformModule: "modules/aws/s3",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "MinIO version", Type: "select", Required: true, Default: "latest", Options: []string{"latest", "RELEASE.2023-10-07T15-07-38Z"}},
			{Name: "accessKey", Description: "Access key", Type: "string", Required: true, Default: "minioadmin"},
			{Name: "secretKey", Description: "Secret key", Type: "string", Required: true},
			{Name: "port", Description: "MinIO port", Type: "number", Required: false, Default: 9000},
		},
	}

	// Message queue resources
	r.blueprints["rabbitmq"] = ResourceBlueprint{
		Name:        "rabbitmq",
		Description: "A RabbitMQ Message Queue",
		Category:    "message-queue",
		DockerComposeSnippet: `
    image: rabbitmq:{{.Version}}-management
    environment:
      - RABBITMQ_DEFAULT_USER={{.Username}}
      - RABBITMQ_DEFAULT_PASS={{.Password}}
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    ports:
      - "{{.Port}}:5672"
      - "15672:15672"
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5`,
		TerraformModule: "modules/aws/mq-rabbitmq",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "RabbitMQ version", Type: "select", Required: true, Default: "3.12", Options: []string{"3.11", "3.12"}},
			{Name: "username", Description: "RabbitMQ username", Type: "string", Required: true, Default: "admin"},
			{Name: "password", Description: "RabbitMQ password", Type: "string", Required: true},
			{Name: "port", Description: "RabbitMQ port", Type: "number", Required: false, Default: 5672},
		},
	}

	r.blueprints["kafka"] = ResourceBlueprint{
		Name:        "kafka",
		Description: "An Apache Kafka Message Queue",
		Category:    "message-queue",
		DockerComposeSnippet: `
    image: confluentinc/cp-kafka:{{.Version}}
    environment:
      KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,PLAINTEXT_HOST://:29092
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:29092
      KAFKA_CFG_PROCESS_ROLES: broker
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: '1@kafka:29093'
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:29093,PLAINTEXT_HOST://:29092
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:29092
      KAFKA_CFG_PROCESS_ROLES: broker
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: '1@kafka:29093'
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:29093,PLAINTEXT_HOST://:29092
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:29092
      KAFKA_CFG_PROCESS_ROLES: broker
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: '1@kafka:29093'
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
    ports:
      - "{{.Port}}:9092"
      - "29092:29092"
    depends_on:
      - zookeeper`,
		TerraformModule: "modules/aws/msk",
		Parameters: []ResourceParameter{
			{Name: "version", Description: "Kafka version", Type: "select", Required: true, Default: "7.4.0", Options: []string{"7.4.0", "7.3.0"}},
			{Name: "port", Description: "Kafka port", Type: "number", Required: false, Default: 9092},
		},
	}
}
