package main

import (
	"github.com/kaneshin/kamimai/core"
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
	current, err := driver.Version().Current()
	if err != nil {
		return err
	}

	// migration files
	mig := core.NewMigration(config).WithVersion(current)
	_ = mig

	// All
	// if err := mig.Down(); err != nil {
	// 	return err
	// }

	// Just one
	// if err := mig.Prev(); err != nil {
	// 	return err
	// }

	return nil
}
