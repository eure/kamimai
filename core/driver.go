package core

import (
	"sync"
)

type (
	// Driver interface.
	Driver interface {
		// Open is the first function to be called.
		// Check the dsn string and open and verify any connection
		// that has to be made.
		Open(string) error

		// Close is the last function to be called.
		// Close any open connection here.
		Close() error

		// Version returns a version interface.
		Version() Version
	}
)

var (
	registry = struct {
		mu      sync.RWMutex
		drivers map[string]Driver
	}{drivers: make(map[string]Driver)}
)

// RegisterDriver a driver so it can be created from its name. Drivers should
// call this from an init() function so that they registers themselvse on
// import
func RegisterDriver(name string, d Driver) {
	if d == nil {
		panic("driver: register driver is nil")
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, dup := registry.drivers[name]; dup {
		panic("sql: register called twice for driver " + name)
	}
	registry.drivers[name] = d
}

// GetDriver retrieves a registered driver by name.
func GetDriver(name string) Driver {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return registry.drivers[name]
}
