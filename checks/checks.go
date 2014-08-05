package checks

import "github.com/travisjeffery/tron/cmd"
import "os"
import "fmt"
import "strings"

type check struct {
	description string
	hint        string
	run         func() bool
}

type suite struct {
	description string
	skip        func() bool
	checks      []check
	failures    []check
	successes   []check
}

type runner struct {
	successes []string
	failures  []string
	suites    []suite
}

func (r *runner) Run() (successes, failures []string) {
	successes = []string{}
	failures = []string{}

	for _, s := range r.suites {
		for _, c := range s.checks {
			if c.run() {
				successes = append(successes, c.description)
				fmt.Println(":white_check_mark:", c.description)
			} else {
				failures = append(failures, c.description)
				fmt.Println(":warning:", c.description)
			}
		}
	}

	return
}

var Runner = runner{
	suites: []suite{
		suite{
			description: "1Password",
			skip: func() bool {
				return ipassword() == ""
			},
			successes: []check{},
			failures:  []check{},
			checks: []check{
				check{
					description: "1Password locks when the app closes",
					hint:        "",
					run: func() bool {
						output, status := cmd.New("defaults").WithArgs("read", "-app", ipassword(), "LockOnMainAppExit").ExecOutput()

						return status == nil && strings.TrimSpace(output) == "1"
					},
				},
			},
		},
	},
}

func ipassword() string {
	paths := []string{
		"/Applications/1Password 4.app",
		"/Applications/1Password.app",
		"/Applications/1Password.localized/1Password.app",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
