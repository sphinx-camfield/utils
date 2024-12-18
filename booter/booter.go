package booter

import "fmt"

type Registra func() (interface{}, error)

type Booter struct {
	registry map[string]Registra
	booted   map[string]interface{}

	// The booting sequence of services.
	// This is useful for checking circular dependencies.
	booting []string
}

func NewBooterWithCached(booted map[string]interface{}) *Booter {
	return &Booter{
		registry: make(map[string]Registra),
		booted:   booted,
		booting:  []string{},
	}
}

func NewBooter() *Booter {
	return NewBooterWithCached(map[string]interface{}{})
}

func (b *Booter) cache(svc string, instance interface{}) {
	b.booted[svc] = instance
}

func (b *Booter) Register(svc string, r Registra) {
	if _, ok := b.registry[svc]; ok {
		panic(fmt.Sprintf("service %s already registered", svc))
	}
	b.registry[svc] = r
}

func (b *Booter) Get(svc string) (svcInstance interface{}, err error) {
	if sInst, ok := b.booted[svc]; ok {
		return sInst, nil
	}

	r, ok := b.registry[svc]

	if !ok {
		return nil, fmt.Errorf("service %s not registered", svc)
	}

	// Check for circular dependencies
	for _, s := range b.booting {
		if s == svc {
			return nil, fmt.Errorf("circular dependency detected: %v", append(b.booting, svc))
		}
	}

	b.booting = append(b.booting, svc)
	sInst, err := r()
	b.booting = b.booting[:len(b.booting)-1]

	if err != nil {
		return nil, err
	}

	b.booted[svc] = sInst

	return sInst, nil
}

func (b *Booter) MustGet(svc string) (svcInstance interface{}) {
	sInst, err := b.Get(svc)
	if err != nil {
		panic(err)
	}
	return sInst
}
