package main

var (
	downCmd = &Cmd{
		Name:  "down",
		Usage: "",
		Run:   doDownCmd,
	}
)

func doDownCmd(cmd *Cmd, args ...string) error {
	return nil
}
