package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/kaneshin/kamimai/core"
	_ "github.com/kaneshin/kamimai/driver"
)

var (
	cmds = []*Cmd{
		upCmd,
		downCmd,
	}

	help    = flag.String("help", "", "show help")
	version = flag.String("version", "", "print the version")

	config *core.Config
)

func main() {
	flag.Usage = usage
	flag.Parse()
	os.Exit(run(flag.Args()))
}

func run(args []string) int {

	if len(args) == 0 {
		flag.Usage()
		return 1
	}

	var cmd *Cmd
	name := args[0]
	for _, c := range cmds {
		if c.Name == name {
			cmd = c
			break
		}
	}

	if cmd == nil {
		fmt.Printf("error: unknown command %q\n", name)
		flag.Usage()
		return 1
	}

	if err := cmd.Exec(args[1:]); err != nil {
		panic(err)
		return 1
	}

	return 0
}

func usage() {
	params := map[string]interface{}{
		"name":        "kamimai",
		"description": "kamimai is a database migration management system.",
		"usage":       "kamimai [global options] command [command options] [arguments...]",
		"version":     "0.0.1",
		"author":      "kaneshin <kaneshin0120@gmail.com>",
		"cmds": []map[string]interface{}{
			map[string]interface{}{
				"name":    "up",
				"summary": "",
			},
			map[string]interface{}{
				"name":    "down",
				"summary": "",
			},
		},
	}
	opts := []map[string]interface{}{}
	flag.VisitAll(func(f *flag.Flag) {
		opt := map[string]interface{}{
			"name":    f.Name,
			"summary": f.Usage,
		}
		opts = append(opts, opt)
	})
	params["opts"] = opts
	helpTemplate.Execute(os.Stdout, params)
}

var helpTemplate = template.Must(template.New("usage").Parse(`
NAME:
  {{.name}} - {{.description}}

USAGE:
  {{.usage}}

VERSION:
  {{.version}}

AUTHOR(S):
  {{.author}}

COMMANDS:{{range .cmds}}
  {{.name | printf "%-18s"}} {{.summary}}{{end}}

GLOBAL OPTIONS:{{range .opts}}
  {{.name | printf "--%-16s"}} {{.summary}}{{end}}
`))
