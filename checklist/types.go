package checklist

import "github.com/fatih/color"

type Test struct {
	description string
	hint        string
	run         func() bool
}

type Suite struct {
	description string
	skip        func() bool
	Tests       []Test
	failures    []Test
	successes   []Test
}

type Runner struct {
	successes []string
	failures  []string
	Suites    []Suite
}

func (r *Runner) Run() (successes, failures []string) {
	successes = []string{}
	failures = []string{}
	successColor := color.New(color.FgGreen)
	failureColor := color.New(color.FgRed).Add(color.Bold)

	for _, s := range r.Suites {
		for _, c := range s.Tests {
			if c.run() {
				successes = append(successes, c.description)
				successColor.Println(" ✓", c.description)
			} else {
				failures = append(failures, c.description)
				failureColor.Println(" ✗", c.description)
			}
		}
	}

	return
}
