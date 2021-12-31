package main

import (
	"fmt"
	"flag"
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
	// TODO append sudo to the beginning of the string
	fmt.Println("Running", c.name, "with:", c.cmd)

	var cmdStr string = c.cmd
	if c.sudo {
		cmdStr = "sudo " + cmdStr
	}

	cmd := exec.Command("bash", "-c", cmdStr)

	// redirecto command output to parent (this go process)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

type InstallOptions struct {
	node bool // indicates if NodeJS should be installed
	java bool // indicates if Java should be installed
	go_ bool // indicates if Go should be installed
	math bool // indicates if math software should be installed
	c bool // indicates if C Development Tools and Libraries should be installed
}

func getOptions() *InstallOptions {
	// define flags
	var opts InstallOptions
	var installAllLangs bool
	flag.BoolVar(&opts.node, "node", false, "Indicate if NodeJS should be installed. Default: false")
	flag.BoolVar(&opts.java, "java", false, "Indicate if Java 11 (devel) should be installed. Default: false")
	flag.BoolVar(&opts.go_, "go", false, "Indicate if Go should be installed. Default: false")
	flag.BoolVar(&opts.c, "c", false, "Indicate if dev tools and libraries for C should be installed. Default: false")
	flag.BoolVar(&opts.math, "math", false, "Indicate if math tools (octave, rstudio, texstudio) should be installed. Default: false")
	flag.BoolVar(
		&installAllLangs, 
		"langs", 
		false, 
		"Indicate if you want to install all languages (NodeJS, Java, Go, C). Default: false. Overrides -node, -java, -go, -c",
	)
	
	// parse CLI flags
	flag.Parse()

	if installAllLangs {
		opts.node, opts.java, opts.go_, opts.c = true, true, true, true
	}
	
	return &opts;
}

func main() {
	var opts *InstallOptions = getOptions()
	var pendingCmds []*Cmd = []*Cmd{
		NewCmd("Update system", "dnf update -y"),
		NewCmd(
			"Install useful programs", 
			"dnf install -y latte-dock terminator thunderbird vim firewalld firewall-config keepassxc chromium kgpg stacer dia git",
		),
		NewCmd("Disable services", "systemctl disable --now cups sshd geoclue"),
		NewCmd("Set firewall in block zone", "firewall-cmd --set-default-zone=block"),
		NewCmd("Remove Package managers", "dnf install -y flatpak snapd"),
		NewCmd("Install Sec & net tools", "dnf install -y nmap wireshark clamav clamav-update rkhunter"),
		NewCmd("Lock root account", "usermod -L root"),
	}

	if opts.node {
		pendingCmds = append(pendingCmds, NewCmd("Install Node JS", "dnf install -y nodejs"))
	}
	if opts.java {
		pendingCmds = append(pendingCmds, NewCmd(
			"Install Java",
			"dnf install -y java-11-openjdk-src java-11-openjdk-javadoc java-11-openjdk-devel maven",
		))
	}
	if opts.go_ {
		pendingCmds = append(pendingCmds, NewCmd("Install Go", "dnf install -y golang"))
	}
	if opts.math {
		pendingCmds = append(pendingCmds, NewCmd("Install Math utilites", "dnf install -y octave rstudio texstudio"))
	}
	if opts.c {
		pendingCmds = append(pendingCmds, NewCmd("Install C Dev tools", "dnf group install -y \"C Development Tools and Libraries\""))
	}

	for _, c := range pendingCmds {
		c.Run()
	}

	fmt.Println("-- Final recommendations --")
	fmt.Println("Install librewolf https://librewolf.net Don't use firefox")
	fmt.Println("Set static IP")
	fmt.Println("Configure the DNS to block adds")

	fmt.Println("\nHave a nice experience with your new Fedora installation 🙃")
}
