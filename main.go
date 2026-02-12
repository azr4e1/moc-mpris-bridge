package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const VERSION = "v0.2.0"

func main() {
	flag.Usage = func() {
		fmt.Printf("%s is a small bridge that implements\n", os.Args[0])
		fmt.Println("the MediaPlayer2 DBus interface for MOC")
		fmt.Println("\nIt's preferrable to run this utility as a systemd service.")
		fmt.Fprintf(os.Stderr, "\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -h, -help\n")
		fmt.Fprintf(os.Stderr, "        Show this help message\n")
		fmt.Fprintf(os.Stderr, "  -v, -version\n")
		fmt.Fprintf(os.Stderr, "        Print current version\n")
		fmt.Fprintf(os.Stderr, "  -n, -name NAME\n")
		fmt.Fprintf(os.Stderr, "        Register interface with NAME. Default: mocp-mpris-bridge\n")
	}

	var version bool
	var name string
	flag.BoolVar(&version, "v", false, "print current version")
	flag.BoolVar(&version, "version", false, "print current version")
	flag.StringVar(&name, "n", "mocp-mpris-bridge", "register service with this name")
	flag.StringVar(&name, "name", "mocp-mpris-bridge", "register service with this name")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "argument not valid")
		os.Exit(1)
	}

	if version {
		fmt.Printf("%s version %s", os.Args[0], VERSION)
		return
	}
	err := MPRISLoop(name)
	if err != nil {
		log.Fatal(err)
	}
}
