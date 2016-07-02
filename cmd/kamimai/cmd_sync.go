package main

import (
	"database/sql"

	"github.com/kaneshin/kamimai/core"
)

var (
	syncCmd = &Cmd{
		Name:  "sync",
		Usage: "apply all migrations",
		Run:   doSyncCmd,
	}
)

func doSyncCmd(cmd *Cmd, args ...string) error {

	// driver
	driver := core.GetDriver(config.Driver())
	if err := driver.Open(config.Dsn()); err != nil {
		return err
	}

	current, err := driver.Version().Current()
	if err != nil {
		return err
	}

	// generate a service
	svc := core.NewService(config).
		WithVersion(current).
		WithDriver(driver)

	return driver.Transaction(func(tx *sql.Tx) error {
		// Sync all migrations
		return svc.Sync()
	})
}
