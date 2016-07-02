package main

import (
	"database/sql"

	"github.com/kaneshin/kamimai/core"
	"github.com/kaneshin/kamimai/internal/cast"
)

var (
	upCmd = &Cmd{
		Name:  "up",
		Usage: "apply all available migrations",
		Run:   doUpCmd,
	}
)

func doUpCmd(cmd *Cmd, args ...string) error {

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

		if len(args) == 0 {
			// All
			if err := svc.Up(); err != nil {
				return err
			}
			return nil
		}

		val := cast.Int(args[0])
		if val == 0 {
			return nil
		}
		return svc.Next(val)
	})
}
