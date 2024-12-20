package booter

import "fmt"

// Registra is a service registration function.
type Registra func() (interface{}, error)

// Booter is a simple service container that supports booting services with dependencies.
type Booter struct {
	registry map[string]Registra
	booted   map[string]interface{}
	aliases  map[string]string
	booting  []string
}

// NewBooterWithCached creates a new Booter with a cached map of booted services.
func NewBooterWithCached(booted map[string]interface{}) *Booter {
	return &Booter{
		registry: make(map[string]Registra),
		aliases:  make(map[string]string),
		booted:   booted,
		booting:  []string{},
	}
}

// NewBooter creates a new Booter.
func NewBooter() *Booter {
	return NewBooterWithCached(map[string]interface{}{})
}

// Cache caches instead of calling the service registration function.
func (b *Booter) Cache(svc string, instance interface{}) {
	b.booted[svc] = instance
}

// Register registers a service with a registration function.
func (b *Booter) Register(svc string, r Registra) {
	if _, ok := b.registry[svc]; ok {
		panic(fmt.Sprintf("service %s already registered", svc))
	}
	b.registry[svc] = r
}

// Get gets a service instance by name.
func (b *Booter) Get(svc string) (svcInstance interface{}) {
	// Check if the service is an alias
	if sourceSvc, ok := b.aliases[svc]; ok {
		// Get the source service
		return b.Get(sourceSvc)
	}

	if sInst, ok := b.booted[svc]; ok {
		return sInst
	}

	r, ok := b.registry[svc]

	if !ok {
		panic("service [" + svc + "] not registered")
	}

	// Check for circular dependencies
	for _, s := range b.booting {
		if s == svc {
			panic(fmt.Sprintf("circular dependency detected: %s", svc))
		}
	}

	b.booting = append(b.booting, svc)
	sInst, err := r()
	b.booting = b.booting[:len(b.booting)-1]

	if err != nil {
		panic(err)
	}

	b.booted[svc] = sInst

	return sInst
}

// Alias creates an alias for a service name.
func (b *Booter) Alias(source string, alias string) {

	// Check if source and alias are the same
	if source == alias {
		panic("source and alias cannot be the same")
	}

	// Check if source is already an alias
	// If it is, set the alias to the source of the alias
	if origin, sourceIsAlias := b.aliases[source]; sourceIsAlias {
		b.Alias(origin, alias)
		return
	}

	// Otherwise, set the alias to the source
	b.aliases[alias] = source
}
