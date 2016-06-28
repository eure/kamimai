package main

import (
	"database/sql"

	"github.com/kaneshin/kamimai/core"
	"github.com/kaneshin/kamimai/internal/cast"
)

var (
	downCmd = &Cmd{
		Name:  "down",
		Usage: "",
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

	if err := driver.Transaction(func(tx *sql.Tx) error {
		if len(args) == 0 {
			// Just one
			if err := svc.Prev(); err != nil {
				return err
			}
			return nil
		}

		for i := cast.Int(args[0]); i < 0; i++ {
			if err := svc.Prev(); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
