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

	initTestCmd()
	err = testCmd.Exec(nil)
	assert.NoError(err)
	assert.Equal(0, len(testCmd.flag.Args()))

	initTestCmd()
	err = testCmd.Exec([]string{"hello"})
	assert.NoError(err)
	assert.Equal(1, len(testCmd.flag.Args()))

	// should be return an error.
	initTestCmd()
	err = testCmd.Exec([]string{"error"})
	assert.Error(err)
	assert.Equal(1, len(testCmd.flag.Args()))
}
