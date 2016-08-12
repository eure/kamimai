package main

import (
	"errors"
	"log"
	"os"

	"github.com/eure/kamimai/core"
)

var (
	createCmd = &Cmd{
		Name:  "create",
		Usage: "create a new migration files",
		Run:   doCreateCmd,
	}
)

func doCreateCmd(cmd *Cmd, args ...string) error {
	// arguments validation
	if len(args) < 1 {
		return errors.New("no file name specified.")
	}

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

	up, down, err := svc.NextMigration(args[len(args)-1])
	if err != nil {
		return err
	}

	// create migration files
	for _, v := range []*core.Migration{up, down} {
		if err := svc.MakeMigrationsDir(); err != nil {
			log.Fatal(err)
		}

		name := v.Name()
		if !*dryRun {
			if _, err := os.Create(name); err != nil {
				log.Fatal(err)
			}
		}

		// print filename on stdout
		log.Printf("created %s", name)
	}

	return nil
}
