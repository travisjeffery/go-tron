package main

import "os"
import . "github.com/visionmedia/go-debug"
import "github.com/codegangsta/cli"
import "github.com/travisjeffery/tron/cmd"
import "github.com/DHowett/go-plist"
import "fmt"
import "log"
import "os/user"
import "io/ioutil"
import "path/filepath"

const reportsGitURL = "reports_git_url"

var debug = Debug("tron")

func main() {

	app := cli.NewApp()
	app.Name = "tron"
	app.Commands = []cli.Command{
		{
			Name:  "report",
			Usage: "Pulls, runs checks, and pushes results to GitHub",
			Action: func(c *cli.Context) {
				debug("tron report")
				// TODO: check for updates to tron and get new checks
				// pull reports
				pullReports()
				initReport()
				recordReport()
			},
		},
		{
			Name:  "pull",
			Usage: "Pulls the latest reported results from GitHub",
			Action: func(c *cli.Context) {
				debug("tron pull")
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

func pullReports() {
	os.Chdir(reportsDir())
	cmd.New("git").WithArgs("pull", "-v").Exec()
}


func recordReport() {
	total, failures := runChecks()

	var status string

	if failures == 0 {
		status = ":white_check_mark:"
	} else {
		status = ":warning:"
	}

	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal(err)
	}

	id := fmt.Sprintf("%s@%s", currentUser().Username, hostname)

	os.Chdir(reportsDir())
	cmd.New("git").WithArgs("pull", "--rebase")

	d := map[string]interface{}{
		"summary": fmt.Sprintf("%s %s: %d checks, %d failures", status, id, total, failures),
		"date":    fmt.Sprintf("%s", time.Now()),
		"checks": map[string][]string{
			"successes": []string{"check that travis is cool successful obv"},
			"failures":  []string{},
		},
	}

	p := filepath.Join(reportsDir(), id)

	f, err := os.Create(p)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(f).Encode(d)
}

func initReport() {
	p := filepath.Join(currentUser().HomeDir, ".tron", reportsGitURL)

	if _, err := os.Stat(p); err != nil {
		log.Fatal(fmt.Sprintf("tron needs a git url in ~/.tron/%s", reportsGitURL))
	}

	url, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(reportsDir()); err == nil {
		debug("resetting reports")
		os.Chdir(reportsDir())
		cmd.New("git").WithArgs("reset", "--hard", "origin/master").Exec()
	} else {
		debug("cloning reports")
		cmd.New("git").WithArgs("clone", "-q", string(url), reportsDir()).Exec()
	}
}

func reportsDir() string {
	return filepath.Join(tronDir(), "reports")
}

func installLaunchAgent() {
	label := "com.travisjeffery.tron"

	p := filepath.Join(currentUser().HomeDir, "Library", "LaunchAgents", label)

	if _, err := os.Stat(p); err == nil {
		cmd.New(fmt.Sprintf("launchtl unload \"%s\"", p)).Exec()
	}

	f, err := os.Create(p)

	if err != nil {
		log.Fatal(err)
	}

	encoder := plist.NewEncoder(f)
	encoder.Encode(map[string]interface{}{
		"Label":             label,
		"StandardOutPath":   "",
		"StandardErrorPath": "",
		"ProgramArguments":  []string{"tron", "report"},
		"StartCalendarInterval": map[string]int{
			"Hour":   15,
			"Minute": 0,
		},
	})

	cmd.New(fmt.Sprintf("launchtl load \"%s\"", p)).Exec()
	cmd.New(fmt.Sprintf("launchtl start \"%s\"", p)).Exec()

	println("tron is installed and will report daily.")
}

func runChecks() (total int, failures int) {
	return 10, 0
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
