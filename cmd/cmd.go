package cmd

import (
    "fmt"
    "os/exec"
    "log"
    "os"
)

type Cmd struct {
	name string
	cmd string
	sudo bool
}

// Creates a new Cmd with sudo: true and the given name and cmd
// cmd: the actual command to be executed
func NewCmd(name string, cmd string) *Cmd {
	return &Cmd{name, cmd, true}
}

func (c *Cmd) String() string{
	return fmt.Sprintf("%s. Command: %s\n", c.name, c.cmd)
}

func (c *Cmd) Run() {
	fmt.Println("\n\033[97;1m", "--- Running", c.name, "with:", c.cmd, "--- \033[0m")

	var cmdStr string = c.cmd
	if c.sudo {
		cmdStr = "sudo " + cmdStr
	}
	cmd := exec.Command("bash", "-c", cmdStr)

	// redirec to command output to parent (this go process)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

