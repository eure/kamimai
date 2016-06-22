package main

var (
	upCmd = &Cmd{
		Name:  "up",
		Usage: "",
		Run:   doUpCmd,
	}
)

func doUpCmd(cmd *Cmd, args ...string) error {
	return nil
}
