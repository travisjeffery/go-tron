package checks

type Check struct {
	description string
	hint        string
	run         func() bool
}

type Suite struct {
	description string
	checks      []Check
}
