package main

import "os"
import . "github.com/visionmedia/go-debug"
import "github.com/codegangsta/cli"
import "github.com/travisjeffery/tron/cmd"
import "fmt"

var debug = Debug("tron")

func main() {

	app := cli.NewApp()
	app.Name = "tron"
	app.Commands = []cli.Command{
		{
			Name:  "report",
			Usage: "Pulls, runs checks, and pushes results to GitHub",
			Action: func(c *cli.Context) {

			},
		},
		{
			Name:  "pull",
			Usage: "Pulls the latest reported results from GitHub",
			Action: func(c *cli.Context) {

			},
		},
		{
			Name:  "install",
			Usage: "Schedules tron to report daily",
			Action: func(c *cli.Context) {
				debug("tron install")

				var plist = "$HOME/Library/LaunchAgents/com.travisjeffery.tron"

				if _, err := os.Stat(plist); err != nil {
					cmd.New(fmt.Sprintf("launchtl unload \"%s\"", plist))
				}
			},
		},
	}
	app.Run(os.Args)
}
