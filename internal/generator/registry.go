package generator

import (
	"fmt"
	"sync"
)

// registry implements the Registry interface
type registry struct {
	generators map[string]Generator
	mutex      sync.RWMutex
}

// NewRegistry creates a new generator registry
func NewRegistry() Registry {
	return &registry{
		generators: make(map[string]Generator),
	}
}

// Register adds a new generator to the registry
func (r *registry) Register(generator Generator) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := generator.Name()
	if name == "" {
		return fmt.Errorf("generator name cannot be empty")
	}

	if _, exists := r.generators[name]; exists {
		return fmt.Errorf("generator '%s' is already registered", name)
	}

	r.generators[name] = generator
	return nil
}

// Get retrieves a generator by name
func (r *registry) Get(name string) (Generator, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	generator, exists := r.generators[name]
	if !exists {
		return nil, fmt.Errorf("generator '%s' not found", name)
	}

	return generator, nil
}

// List returns all available generators
func (r *registry) List() []Generator {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	generators := make([]Generator, 0, len(r.generators))
	for _, generator := range r.generators {
		generators = append(generators, generator)
	}

	return generators
}

// Names returns all available generator names
func (r *registry) Names() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.generators))
	for name := range r.generators {
		names = append(names, name)
	}

	return names
}
