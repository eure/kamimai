package main

import (
	"github.com/kaneshin/kamimai/core"
)

var (
	upCmd = &Cmd{
		Name:  "up",
		Usage: "",
		Run:   doUpCmd,
	}
)

func doUpCmd(cmd *Cmd, args ...string) error {

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
	// if err := mig.Up(); err != nil {
	// 	return err
	// }

	// Just one
	// if err := mig.Next(); err != nil {
	// 	return err
	// }

	return nil
}
