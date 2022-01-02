package main

import (
	"fmt"
	"flag"
    "github.com/BenjaminGuzman/fedora-new-install/cmd"
)

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
	var pendingCmds []*cmd.Cmd = []*cmd.Cmd{
		cmd.NewCmd("Update system", "dnf update -y"),
		cmd.NewCmd(
			"Install useful programs", 
			"dnf install -y latte-dock terminator thunderbird vim firewalld firewall-config keepassxc chromium kgpg stacer dia git",
		),
		cmd.NewCmd("Disable services", "systemctl disable --now cups sshd geoclue"),
        cmd.NewCmd("Set firewall default zone to block", "firewall-cmd --set-default-zone=block"),
		cmd.NewCmd("Remove package managers", "dnf install -y flatpak snapd"),
		cmd.NewCmd("Install sec & net tools", "dnf install -y nmap wireshark clamav clamav-update rkhunter"),
		cmd.NewCmd("Lock root account", "usermod -L root"),
	}

	if opts.node {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install Node JS", "dnf install -y nodejs"))
	}
	if opts.java {
		pendingCmds = append(pendingCmds, cmd.NewCmd(
			"Install Java",
			"dnf install -y java-11-openjdk-src java-11-openjdk-javadoc java-11-openjdk-devel maven",
		))
	}
	if opts.go_ {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install Go", "dnf install -y golang"))
	}
	if opts.math {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install Math utilites", "dnf install -y octave rstudio texstudio"))
	}
	if opts.c {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install C Dev tools", "dnf group install -y \"C Development Tools and Libraries\""))
	}

	for _, c := range pendingCmds {
		c.Run()
	}

    fmt.Println("\n------------------------------")
    fmt.Println("All commands were run")
	fmt.Println("-- Final recommendations --")
	fmt.Println("1. Install librewolf https://librewolf.net Don't use firefox")
	fmt.Println("2. Set static IP")
	fmt.Println("3. Configure the DNS to block adds. See https://github.com/Ultimate-Hosts-Blacklist/Ultimate.Hosts.Blacklist")
    fmt.Println("4. Disable file search")

	fmt.Println("\nHave a nice experience with your new Fedora installation 🙃")
}
