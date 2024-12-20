package booter

import "fmt"

// Registra is a service registration function.
type Registra func() (interface{}, error)

// Container is a simple service container that supports booting services with dependencies.
type Container struct {
	registry map[string]Registra
	booted   map[string]interface{}
	aliases  map[string]string
	booting  []string
}

// NewBooterWithCached creates a new Booter with a cached map of booted services.
func NewBooterWithCached(booted map[string]interface{}) *Container {
	return &Container{
		registry: make(map[string]Registra),
		aliases:  make(map[string]string),
		booted:   booted,
		booting:  []string{},
	}
}

// NewContainer creates a new Booter.
func NewContainer() *Container {
	return NewBooterWithCached(map[string]interface{}{})
}

// Cache caches instead of calling the service registration function.
func (c *Container) Cache(svc string, instance interface{}) {
	c.booted[svc] = instance
}

// Register registers a service with a registration function.
func (c *Container) Register(svc string, r Registra) {
	if _, ok := c.registry[svc]; ok {
		panic(fmt.Sprintf("service %s already registered", svc))
	}
	c.registry[svc] = r
}

// Get gets a service instance by name.
func (c *Container) Get(svc string) (svcInstance interface{}) {
	// Check if the service is an alias
	if sourceSvc, ok := c.aliases[svc]; ok {
		// Get the source service
		return c.Get(sourceSvc)
	}

	if sInst, ok := c.booted[svc]; ok {
		return sInst
	}

	r, ok := c.registry[svc]

	if !ok {
		panic("service [" + svc + "] not registered")
	}

	// Check for circular dependencies
	for _, s := range c.booting {
		if s == svc {
			panic(fmt.Sprintf("circular dependency detected: %s", svc))
		}
	}

	c.booting = append(c.booting, svc)
	sInst, err := r()
	c.booting = c.booting[:len(c.booting)-1]

	if err != nil {
		panic(err)
	}

	c.booted[svc] = sInst

	return sInst
}

// Alias creates an alias for a service name.
func (c *Container) Alias(source string, alias string) {

	// Check if source and alias are the same
	if source == alias {
		panic("source and alias cannot be the same")
	}

	// Check if source is already an alias
	// If it is, set the alias to the source of the alias
	if origin, sourceIsAlias := c.aliases[source]; sourceIsAlias {
		c.Alias(origin, alias)
		return
	}

	// Otherwise, set the alias to the source
	c.aliases[alias] = source
}
