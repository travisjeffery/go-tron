package checklist

import "fmt"

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

	for _, s := range r.Suites {
		for _, c := range s.Tests {
			if c.run() {
				successes = append(successes, c.description)
				fmt.Println(":white_test_mark:", c.description)
			} else {
				failures = append(failures, c.description)
				fmt.Println(":warning:", c.description)
			}
		}
	}

	return
}
