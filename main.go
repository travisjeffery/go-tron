package main

import "os"
import . "github.com/visionmedia/go-debug"
import "github.com/codegangsta/cli"
import "github.com/travisjeffery/tron/cmd"
import . "github.com/travisjeffery/tron/checklist"
import "github.com/DHowett/go-plist"
import "fmt"
import "log"
import "os/user"
import "io/ioutil"
import "path/filepath"
import "encoding/json"
import "time"
import "bitbucket.org/kardianos/osext"

var debug = Debug("tron")

func pullReports() {
	os.Chdir(reportsDir())
	cmd.New("git pull -v").Exec()
}

func pushReports() {
	os.Chdir(reportsDir())
	cmd.New("git pull --rebase").Exec()
	cmd.New("git push").Exec()
}

func recordReport() {
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal(err)
	}

	userId := fmt.Sprintf("%s@%s", currentUser().Username, hostname)

	os.Chdir(reportsDir())
	cmd.New("git pull --rebase").Exec()

	successes, failures := run()
	failuresCount := len(failures)
	totalCount := len(successes) + failuresCount

	var status string

	if failuresCount == 0 {
		status = ":white_check_mark:"
	} else {
		status = ":warning:"
	}

	summary := fmt.Sprintf("%s %s - %d tests, %d failures", status, userId, totalCount, failuresCount)

	d := map[string]interface{}{
		"summary": summary,
		"date":    fmt.Sprintf("%s", time.Now()),
		"tests": map[string][]string{
			"successes": successes,
			"failures":  failures,
		},
	}

	p := filepath.Join(reportsDir(), userId)

	f, err := os.Create(p)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(f).Encode(d)

	cmd.New("git add").WithArg(userId).Exec()
	cmd.New("git commit -m").WithArg(summary).Exec()
}

func run() ([]string, []string) {
	return Checklist.Run()
}

func initReport() {
	p := filepath.Join(currentUser().HomeDir, ".tron", "reports_git_url")

	if _, err := os.Stat(p); err != nil {
		log.Fatal(fmt.Sprintf("tron needs a git url in ~/.tron/%s", "reports_git_url"))
	}

	url, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(reportsDir()); err == nil {
		debug("resetting reports")
		os.Chdir(reportsDir())
		cmd.New("git reset --hard origin/master").Exec()
	} else {
		debug("cloning reports")
		cmd.New("git clone -q").WithArgs(string(url), reportsDir()).Exec()
	}
}

func reportsDir() string {
	return filepath.Join(tronDir(), "reports")
}

func installLaunchAgent() {
	label := "com.travisjeffery.tron"

	p := filepath.Join(currentUser().HomeDir, "Library", "LaunchAgents", fmt.Sprintf("%s.plist", label))

	if _, err := os.Stat(p); err == nil {
		cmd.New("launchctl unload").WithArg(p).Exec()
	}

	f, err := os.Create(p)

	if err != nil {
		log.Fatal(err)
	}

	execPath, err := osext.Executable()

	if err != nil {
		log.Fatal(err)
	}

	encoder := plist.NewEncoder(f)
	encoder.Encode(map[string]interface{}{
		"Label":            label,
		"ProgramArguments": []string{execPath, "report"},
		"StartCalendarInterval": map[string]int{
			"Hour":   15,
			"Minute": 0,
		},
	})

	cmd.New("launchctl load").WithArg(p).Exec()
	cmd.New("launchctl start").WithArg(label).Exec()

	println("tron is installed and will report daily.")
}

func currentUser() *user.User {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func tronDir() string {
	return filepath.Join(currentUser().HomeDir, ".tron")
}

func main() {
	app := cli.NewApp()
	app.Name = "tron"
	app.Commands = []cli.Command{
		{
			Name:  "report",
			Usage: "Pulls, runs checklist, and pushes results to GitHub",
			Action: func(c *cli.Context) {
				debug("tron report")
				pullReports()
				initReport()
				recordReport()
				pushReports()
			},
		},
		{
			Name:  "run",
			Usage: "Runs checklists",
			Action: func(c *cli.Context) {
				debug("tron run")
				cmd.Stdout = nil
				cmd.Stderr = nil
				cmd.Stdin = nil
				run()
			},
		},
		{
			Name:  "pull",
			Usage: "Pulls the latest reported results from GitHub",
			Action: func(c *cli.Context) {
				debug("tron pull")
				pullReports()
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
