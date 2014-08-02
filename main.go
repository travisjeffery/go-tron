package main

import "os"
import "github.com/codegangsta/cli"

func main() {
	app := cli.NewApp()
	app.Name = "tron"
	app.Run(os.Args)
}
