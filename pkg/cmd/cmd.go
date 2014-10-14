package cmd

import "fmt"
import "os/exec"
import "strings"
import "github.com/kballard/go-shellquote"
import "log"
import "os"
import "io"
import . "github.com/visionmedia/go-debug"

var debug = Debug("cmd")

type Cmd struct {
	Name string
	Args []string
	Cmd  *exec.Cmd
}

func (cmd Cmd) String() string {
	return fmt.Sprintf("%s %s", cmd.Name, strings.Join(cmd.Args, " "))
}

func (cmd *Cmd) WithArg(arg string) *Cmd {
	if arg != "" {
		cmd.Args = append(cmd.Args, arg)
	}
	return cmd
}

func (cmd *Cmd) WithArgs(args ...string) *Cmd {
	for _, arg := range args {
		cmd.WithArg(arg)
	}
	return cmd
}

func (cmd *Cmd) ExecOutput() (string, error) {
	output, err := cmd.Cmd.CombinedOutput()
	return string(output), err
}

func (cmd *Cmd) Exec() error {
	return cmd.Cmd.Run()
}

func New(cmd string) *Cmd {
	cmds, err := shellquote.Split(cmd)

	if err != nil {
		log.Fatal(err)
	}

	return NewWithArray(cmds)
}

func (cmd *Cmd) StdinPipe() (io.WriteCloser, error) {
	return cmd.Cmd.StdinPipe()
}

func NewWithArray(cmd []string) *Cmd {
	name := cmd[0]

	binary, _ := exec.LookPath(name)

	args := make([]string, 0)
	for _, arg := range cmd[1:] {
		args = append(args, arg)
	}

	debug("exec.Command(%s, %s)", binary, args)

	c := exec.Command(binary, args...)

	return &Cmd{Name: name, Args: args, Cmd: c}
}

var Stdout = os.Stdout
var Stderr = os.Stderr
var Stdin = os.Stdin
