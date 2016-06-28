package core

import (
	"database/sql"
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

		// Ext returns the sql file extension used by path. The extension is the
		// suffix beginning at the final dot in the final element of path; it is
		// empty if there is no dot.
		Ext() string

		// Transaction starts a db transaction. The isolation level is dependent on the
		// driver.
		Transaction(func(*sql.Tx) error) error

		// Migrate is the heart of the driver.
		// It will receive a file which the driver should apply
		// to its backend or whatever. The migration function should use
		// the pipe channel to return any errors or other useful information.
		Migrate(*Migration) error

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
