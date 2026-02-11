package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const VERSION = "v0.1.2"

func main() {
	flag.Usage = func() {
		fmt.Printf("%s is a small bridge that implements\n", os.Args[0])
		fmt.Println("the MediaPlayer2 DBus interface for MOC")
		fmt.Println("\nIt's preferrable to run this utility as a systemd service.")
	}

	var version bool
	flag.BoolVar(&version, "v", false, "print current version")
	flag.BoolVar(&version, "version", false, "print current version")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "argument not valid")
		os.Exit(1)
	}

	if version {
		fmt.Printf("%s version %s", os.Args[0], VERSION)
		return
	}
	err := MPRISLoop()
	if err != nil {
		log.Fatal(err)
	}
}
