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

	// All
	// if err := svc.Up(); err != nil {
	// 	return err
	// }

	// Just one
	if err := svc.Next(); err != nil {
		return err
	}

	return nil
}
