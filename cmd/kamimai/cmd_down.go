package main

import (
	"database/sql"

	"github.com/kaneshin/kamimai/core"
	"github.com/kaneshin/kamimai/internal/cast"
	"github.com/kaneshin/kamimai/internal/direction"
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
			if err := svc.Prev(1); err != nil {
				return err
			}
			return nil
		}

		val := cast.Int(args[0])
		if val == 0 {
			return nil
		}
		return svc.Prev(val)
	})
}
