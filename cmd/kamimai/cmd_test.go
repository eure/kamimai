package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCmd     *Cmd
	initTestCmd = func() {
		testCmd = &Cmd{
			Name:  "test",
			Usage: "",
			Run: func(cmd *Cmd, args ...string) error {
				for _, v := range args {
					if v == "error" {
						return fmt.Errorf("just error")
					}
				}
				return nil
			},
		}
	}
)

func TestCmd(t *testing.T) {
	assert := assert.New(t)

	var err error

	*dirPath = "../../examples/testdata"
	*env = "development"
	args := []string{}

	initTestCmd()
	err = testCmd.Exec(args)
	assert.NoError(err)
	assert.Equal(0, len(testCmd.flag.Args()))

	// should be return an error.
	initTestCmd()
	err = testCmd.Exec(append(args, "error"))
	assert.Error(err)
	assert.Equal(1, len(testCmd.flag.Args()))
}
