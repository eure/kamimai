package main

import (
	"flag"

	"github.com/kaneshin/kamimai/core"
)

type (
	// A Cmd executes a command
	Cmd struct {
		Name  string
		Usage string
		Run   func(*Cmd, ...string) error

		flag flag.FlagSet
	}
)

// Exec executes a command with arguments.
func (c *Cmd) Exec(args []string) error {
	c.flag.Parse(args)

	// Load config
	config = core.MustNewConfig(*dirPath).WithEnv(*env)
	return c.Run(c, c.flag.Args()...)
}
