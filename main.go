package main

import "os"
import "github.com/codegangsta/cli"
import "github.com/travisjeffery/tron/pkg/reports"
import "github.com/travisjeffery/tron/pkg/cmd"

func main() {
	app := cli.NewApp()
	app.Name = "tron"
	addCommands(app)
	app.Run(os.Args)
}

func addCommands(app *cli.App) {
	r := reports.New()

	app.Commands = []cli.Command{
		{
			Name:  "report",
			Usage: "Pulls, runs checklist, and pushes results to GitHub",
			Action: func(c *cli.Context) {
				r.Report()
			},
		},
		{
			Name:  "run",
			Usage: "Runs checklists",
			Action: func(c *cli.Context) {
				cmd.Stdout = nil
				cmd.Stderr = nil
				cmd.Stdin = nil
				r.Run()
			},
		},
		{
			Name:  "pull",
			Usage: "Pulls the latest reported results from GitHub",
			Action: func(c *cli.Context) {
				r.Pull()
			},
		},
		{
			Name:  "install",
			Usage: "Schedules tron to report daily",
			Action: func(c *cli.Context) {
				r.Install()
			},
		},
	}

}
