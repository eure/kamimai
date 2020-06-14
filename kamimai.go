package kamimai

import (
	"database/sql"

	"github.com/eure/kamimai/core"
)

// Version represents kamimai's semantic version.
const Version = "v0.4.2"

// Current returns the current migration version.
func Current(c *core.Config) (uint64, error) {

	// driver
	driver := core.GetDriver(c.Driver())
	if err := driver.Open(c.Dsn()); err != nil {
		return 0, err
	}

	// current version
	return driver.Version().Current()
}

// Sync applies non-applied migration files.
func Sync(c *core.Config) error {

	// driver
	driver := core.GetDriver(c.Driver())
	if err := driver.Open(c.Dsn()); err != nil {
		return err
	}

	// current version
	current, err := driver.Version().Current()
	if err != nil {
		return err
	}

	// create a service
	svc := core.NewService(c).
		WithVersion(current).
		WithDriver(driver)

	return driver.Transaction(func(tx *sql.Tx) error {
		// Sync all migrations
		return svc.Sync()
	})
}

// Up applies up migration files.
func Up(c *core.Config) error {

	// driver
	driver := core.GetDriver(c.Driver())
	if err := driver.Open(c.Dsn()); err != nil {
		return err
	}

	// current version
	current, err := driver.Version().Current()
	if err != nil {
		return err
	}

	// create a service
	svc := core.NewService(c).
		WithVersion(current).
		WithDriver(driver)

	return driver.Transaction(func(tx *sql.Tx) error {
		// Up migrations
		return svc.Up()
	})
}

// Down applies down migration files.
func Down(c *core.Config) error {

	// driver
	driver := core.GetDriver(c.Driver())
	if err := driver.Open(c.Dsn()); err != nil {
		return err
	}

	// current version
	current, err := driver.Version().Current()
	if err != nil {
		return err
	}

	// create a service
	svc := core.NewService(c).
		WithVersion(current).
		WithDriver(driver)

	return driver.Transaction(func(tx *sql.Tx) error {
		// Down migrations
		return svc.Down()
	})
}
