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

	svc := core.NewService(config).WithVersion(current)
	_ = svc

	// All
	// if err := svc.Down(); err != nil {
	// 	return err
	// }

	// Just one
	// if err := svc.Prev(); err != nil {
	// 	return err
	// }

	return nil
}
