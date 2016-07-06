package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/eure/kamimai/core"
	_ "github.com/eure/kamimai/driver"
)

var (
	cmds = []*Cmd{
		createCmd,
		upCmd,
		downCmd,
		syncCmd,
		// migrateCmd,
	}

	help    = flag.String("help", "", "show help")
	dirPath = flag.String("path", "", "migration dir containing config")
	env     = flag.String("env", "", "config environment to use")
	dryRun  = flag.Bool("dry-run", false, "")

	config *core.Config
)

func init() {
	log.SetPrefix("kamimai: ")
	log.SetFlags(0)
}

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
		log.Fatal(err)
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
	}

	params["opts"] = func() (list []map[string]interface{}) {
		flag.VisitAll(func(f *flag.Flag) {
			opt := map[string]interface{}{
				"name":    f.Name,
				"summary": f.Usage,
			}
			list = append(list, opt)
		})
		return
	}()

	params["cmds"] = func() (list []map[string]interface{}) {
		for _, c := range cmds {
			cmd := map[string]interface{}{
				"name":    c.Name,
				"summary": c.Usage,
			}
			list = append(list, cmd)
		}
		return
	}()

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
