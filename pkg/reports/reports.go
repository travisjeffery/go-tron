package reports

import "os"
import . "github.com/visionmedia/go-debug"
import "github.com/travisjeffery/tron/pkg/cmd"
import . "github.com/travisjeffery/tron/pkg/utils"
import "github.com/DHowett/go-plist"
import "fmt"
import "os/user"
import "io/ioutil"
import "path/filepath"
import "encoding/json"
import "time"
import "bitbucket.org/kardianos/osext"
import "github.com/Merovius/go-tap"
import "bytes"

// import "bufio"

var debug = Debug("tron")

type Reports struct {
	Dir string
}

type Report struct {
	Summary string              `json:"summary"`
	Date    string              `json:"date"`
	Tests   map[string][]string `json:"tests"`
}

func New() *Reports {
	return &Reports{
		Dir: filepath.Join(filepath.Join(currentUser().HomeDir, ".tron"), "reports"),
	}
}

func (r *Reports) Pull() {
	debug("pull")

	os.Chdir(r.Dir)
	cmd.New("git pull -v").Exec()
}

func (r *Reports) Push() {
	os.Chdir(r.Dir)
	cmd.New("git pull --rebase").Exec()
	cmd.New("git push").Exec()
}

func (r *Reports) record() {
	hostname, err := os.Hostname()

	CheckErr(err)

	userId := fmt.Sprintf("%s@%s", currentUser().Username, hostname)
	os.Chdir(r.Dir)
	cmd.New("git pull --rebase").Exec()

	successes, failures := r.Run()
	failuresCount := len(failures)
	totalCount := len(successes) + failuresCount

	var status string

	if failuresCount == 0 {
		status = ":white_check_mark:"
	} else {
		status = ":warning:"
	}

	summary := fmt.Sprintf("%s %s - %d tests, %d failures", status, userId, totalCount, failuresCount)

	d := Report{
		Summary: summary,
		Date:    fmt.Sprintf("%s", time.Now()),
		Tests: map[string][]string{
			"successes": successes,
			"failures":  failures,
		},
	}

	p := filepath.Join(r.Dir, userId)

	f, err := os.Create(p)
	CheckErr(err)

	json.NewEncoder(f).Encode(d)

	cmd.New("git add").WithArg(userId).Exec()
	cmd.New("git commit -m").WithArg(summary).Exec()
}

func (r *Reports) Run() (successes []string, failures []string) {
	debug("run")

	// execPath, err := osext.Executable()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// arr := []string{"find", "-L", "-s", filepath.Join(execPath, "etc/checklists"), "-type", "f", "-name", "*.bats", "-print0"}

	arr := []string{"find", "-L", "-s", "/Users/tj/dev/go/src/github.com/travisjeffery/tron/etc/checklists", "-type", "f", "-name", "*.bats", "-print0"}

	find := cmd.NewWithArray(arr)
	xargs := cmd.NewWithArray([]string{"xargs", "-0", "/Users/tj/dev/go/src/github.com/travisjeffery/tron/vendor/bats/bin/bats", "--pretty"})
	findOut, err := find.Cmd.StdoutPipe()
	find.Cmd.Start()
	xargs.Cmd.Stdin = findOut

	xargsOutput, err := xargs.Cmd.Output()
	fmt.Println(string(xargsOutput))
	fmt.Println(err)
	// b := cmd.New(filepath.Join(execPath, "vendor/bats/bin/bats"))

	var batsout bytes.Buffer

	debug("batsout: %s", batsout.String())

	bs := bytes.NewReader(batsout.Bytes())

	parser, err := tap.NewParser(bs)
	CheckErr(err)

	s, err := parser.Suite()
	CheckErr(err)

	successes = []string{}
	failures = []string{}

	for _, t := range s.Tests {
		if t.Ok {
			successes = append(successes, t.Description)
		} else {
			failures = append(failures, t.Description)
		}
	}

	return
}

func (r *Reports) Install() {
	debug("install")

	label := "com.travisjeffery.tron"

	p := filepath.Join(currentUser().HomeDir, "Library", "LaunchAgents", fmt.Sprintf("%s.plist", label))

	if _, err := os.Stat(p); err == nil {
		cmd.New("launchctl unload").WithArg(p).Exec()
	}

	f, err := os.Create(p)
	CheckErr(err)

	execPath, err := osext.Executable()
	CheckErr(err)

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

func (r *Reports) init() {
	p := filepath.Join(currentUser().HomeDir, ".tron", "reports_git_url")

	url, err := ioutil.ReadFile(p)
	CheckErr(err)

	if _, err := os.Stat(r.Dir); err == nil {
		debug("resetting reports")

		os.Chdir(r.Dir)
		cmd.New("git reset --hard origin/master").Exec()
	} else {
		debug("cloning reports")

		cmd.New("git clone -q").WithArgs(string(url), r.Dir).Exec()
	}
}

func (r *Reports) Report() {
	debug("report")

	r.Pull()
	r.init()
	r.record()
	r.Push()
}

func currentUser() *user.User {
	user, err := user.Current()
	CheckErr(err)
	return user
}
