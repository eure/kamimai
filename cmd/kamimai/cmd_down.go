package main

import (
	"database/sql"

	"github.com/eure/kamimai/core"
	"github.com/eure/kamimai/internal/cast"
	"github.com/eure/kamimai/internal/direction"
)

var (
	downCmd = &Cmd{
		Name:  "down",
		Usage: "rollback the latest applied migration",
		Run:   doDownCmd,
	}
)

func doDownCmd(cmd *Cmd, args ...string) error {

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

		if version != 0 {
			return svc.Apply(direction.Down, version)
		}

		if len(args) == 0 {
			// Just one
			return svc.Prev(1)
		}

		val := cast.Int(args[0])
		if val == 0 {
			return nil
		}
		return svc.Prev(val)
	})
}
