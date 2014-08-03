package main

import "os"
import . "github.com/visionmedia/go-debug"
import "github.com/codegangsta/cli"
import "github.com/travisjeffery/tron/cmd"
import "github.com/DHowett/go-plist"
import "fmt"
import "log"

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
				installLaunchAgent()
			},
		},
	}
	app.Run(os.Args)
}

func installLaunchAgent() {
	label := "com.travisjeffery.tron"
	p := fmt.Sprintf("$HOME/Library/LaunchAgents/%s", label)

	if _, err := os.Stat(p); err != nil {
		cmd.New(fmt.Sprintf("launchtl unload \"%s\"", p)).Exec()
	}

	f, err := os.Create(p)

	if err != nil {
		log.Fatal(err)
	}

	encoder := plist.NewEncoder(f)
	encoder.Encode(map[string]string{
		"Label":             label,
		"StandardOutPath":   "",
		"StandardErrorPath": "",
	})
	encoder.Encode(map[string][]string{
		// tron needs to be the path to the tron executable here
		"ProgramArguments": []string{"tron", "report"},
	})
	encoder.Encode(map[string]map[string]int{
		"StartCalendarInterval": map[string]int{
			"Hour":   15,
			"Minute": 0,
		},
	})

	cmd.New(fmt.Sprintf("launchtl load \"%s\"", p)).Exec()
	cmd.New(fmt.Sprintf("launchtl start \"%s\"", p)).Exec()

	println("tron is installed and will report daily.")
}
