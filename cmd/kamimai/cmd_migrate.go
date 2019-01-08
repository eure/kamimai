package main

import (
	"strconv"

	"github.com/mr04vv/kamimai/internal/cast"
)

var (
	migrateCmd = &Cmd{
		Name:  "migrate",
		Usage: "apply or rollback the number of n migrations",
		Run:   doMigrateCmd,
	}
)

func doMigrateCmd(cmd *Cmd, args ...string) error {

	if len(args) == 0 {
		return nil
	}

	// FIXME:
	// -1 couldn't be passed from args.
	val := cast.Int(args[0])

	switch {
	case val > 0:
		// kamimai up n
		return doUpCmd(upCmd, strconv.Itoa(val))

	case val < 0:
		// kamimai down n
		return doDownCmd(downCmd, strconv.Itoa(-val))

	}

	return nil
}
