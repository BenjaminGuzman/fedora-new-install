package main

import (
	"flag"
	"fmt"
	"github.com/BenjaminGuzman/fedora-new-install/cmd"
)

type InstallOptions struct {
	node      bool // indicates if NodeJS should be installed
	java      bool // indicates if Java should be installed
	go_       bool // indicates if Go should be installed
	math      bool // indicates if math software should be installed
	c         bool // indicates if C Development Tools and Libraries should be installed
	librewolf bool
	ohMyZsh   bool
}

func getOptions() *InstallOptions {
	// define flags
	var opts InstallOptions
	var installAllLangs bool
	flag.BoolVar(&opts.node, "node", false, "Indicate if NodeJS should be installed. Default: false")
	flag.BoolVar(&opts.java, "java", false, "Indicate if Java 11 & 17 (devel) should be installed. Default: false")
	flag.BoolVar(&opts.go_, "go", false, "Indicate if Go should be installed. Default: false")
	flag.BoolVar(&opts.c, "c", false, "Indicate if dev tools and libraries for C should be installed. Default: false")
	flag.BoolVar(&opts.math, "math", false, "Indicate if math tools (octave, rstudio, texstudio) should be installed. Default: false")
	flag.BoolVar(
		&installAllLangs,
		"langs",
		false,
		"Indicate if you want to install all languages (NodeJS, Java, Go, C). Overrides --node, --java, --go, --c. Default: false",
	)
	flag.BoolVar(&opts.librewolf, "librewolf", false, "Indicate if LibreWolf should be installed. Default: false")
	flag.BoolVar(&opts.ohMyZsh, "ohmyzsh", false, "Indicate if Oh My Zsh should be installed. Default: false")

	// parse CLI flags
	flag.Parse()

	if installAllLangs {
		opts.node, opts.java, opts.go_, opts.c = true, true, true, true
	}

	return &opts
}

func main() {
	var opts *InstallOptions = getOptions()
	var pendingCmds []*cmd.Cmd = []*cmd.Cmd{
		cmd.NewCmd("Update system", "dnf update -y --refresh"),
		cmd.NewCmd(
			"Install useful programs",
			"dnf install -y terminator thunderbird vim firewalld firewall-config chromium kgpg stacer dia git htop peek",
		),
		cmd.NewCmd("Disable services", "systemctl disable --now cups sshd geoclue"),
		cmd.NewCmd("Set firewall default zone to block", "firewall-cmd --set-default-zone=block"),
		cmd.NewCmd("Remove package managers", "dnf remove -y flatpak snapd"),
		cmd.NewCmd("Install security & network tools", "dnf install -y nmap wireshark clamav clamav-update rkhunter"),
		cmd.NewCmd("Lock root account", "usermod --expiredate 1 -L root"),
	}

	if opts.node {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install Node JS", "dnf install -y nodejs"))
	}
	if opts.java {
		pendingCmds = append(
			pendingCmds,
			cmd.NewCmd("Install maven", "dnf install -y maven"),
			cmd.NewCmd(
				"Install Java 11",
				"dnf install -y java-11-openjdk-src java-11-openjdk-javadoc java-11-openjdk-devel",
			),
			cmd.NewCmd(
				"Install Java 17",
				"dnf install -y java-17-openjdk-src java-17-openjdk-javadoc java-17-openjdk-devel",
			),
		)
	}
	if opts.go_ {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install Go", "dnf install -y golang"))
	}
	if opts.math {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install Math utilities", "dnf install -y octave rstudio texstudio"))
	}
	if opts.c {
		pendingCmds = append(pendingCmds, cmd.NewCmd("Install C Dev tools", "dnf group install -y \"C Development Tools and Libraries\""))
	}
	if opts.librewolf {
		pendingCmds = append(
			pendingCmds,
			cmd.NewCmd(
				"Import LibeWolf GPG key",
				"rpm --import https://keys.openpgp.org/vks/v1/by-fingerprint/034F7776EF5E0C613D2F7934D29FBD5F93C0CFC3",
			),
			cmd.NewCmd(
				"Add LibreWold repo",
				"dnf config-manager --add-repo https://rpm.librewolf.net",
			),
			cmd.NewCmd(
				"Install LibreWolf",
				"dnf install -y librewolf",
			),
		)
	}
	if opts.ohMyZsh {
		pendingCmds = append(
			pendingCmds,
			cmd.NewCmd("Install zsh", "dnf install -y zsh"),
			cmd.NewCmd("Install oh-my-zsh", "echo Installation of oh my zsh is interactive, so it is not run by this program"),
		)
	}

	for _, c := range pendingCmds {
		c.Run()
	}

	fmt.Println("\n------------------------------")
	fmt.Println("All commands were run")
	fmt.Println("-- Final recommendations --")
	fmt.Println("1. Don't use firefox. Use librewolf instead")
	fmt.Println("2. Set static IP")
	fmt.Println("3. Configure the DNS to block adds. See https://github.com/Ultimate-Hosts-Blacklist/Ultimate.Hosts.Blacklist")
	fmt.Println("4. Disable file search")
	fmt.Println("5. Import SSH & GPG keys if needed")

	fmt.Println("\nHave a nice experience with your new Fedora installation 🙃")
}
