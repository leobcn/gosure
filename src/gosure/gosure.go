// File integrity testing.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var surefileArg = flag.String("file", "2sure", "base name of surefile, will have .dat.gz appended")
var helpArg = flag.Bool("help", false, "Ask for help")

func usage(message string) {
	fmt.Printf("error: %s\n", message)
	fmt.Printf("usage: gosure [{-f|--file} name] {scan|update|check|signoff|show}\n\n")
	os.Exit(1)
}

func sureName(suffix string) string {
	return *surefileArg + "." + suffix + ".gz"
}

func main() {
	flag.Parse()
	if *helpArg {
		usage("Help")
	}
	if flag.NArg() != 1 {
		usage("Expecting a single command")
	}

	switch flag.Arg(0) {
	case "scan":
		dir, err := WalkRoot(".")
		if err != nil {
			log.Fatalf("Unable to walk root directory: %s", err)
		}
		defer dir.Close()

		err = writeSure(sureName("0"), dir)
		if err != nil {
			log.Fatalf("Error writing surefile: %s", err)
		}
		// TODO: Handle these?
		_ = os.Rename(sureName("dat"), sureName("bak"))
		_ = os.Rename(sureName("0"), sureName("dat"))
	case "check":
		dir1, err := ReadSure(sureName("dat"))
		if err != nil {
			log.Fatalf("Unable to read surefile: %s", err)
		}
		defer dir1.Close()

		dir2, err := WalkRoot(".")
		if err != nil {
			log.Fatalf("Error walking root dir: %s", err)
		}
		defer dir2.Close()

		Compare(dir1, dir2)

	case "signoff":
		left, err := ReadSure(sureName("bak"))
		if err != nil {
			log.Fatalf("Unable to read backup surefile: %s", err)
		}
		right, err := ReadSure(sureName("dat"))
		if err != nil {
			log.Fatalf("Unable to read surefile: %s", err)
		}
		Compare(left, right)
	case "tmp":
		in, err := ReadSure(sureName("0"))
		if err != nil {
			log.Fatalf("Unable to read surefile: %s", err)
		}
		defer in.Close()

		writeSure("tmp", in)
	default:
		usage("Unexpected command")
	}
}

// TODO: These really belong somewhere else.

const magic = "asure-2.0\n-----\n"

type Node struct {
	name   string
	atts   map[string]string
	costly func() map[string]string // Get the atts that are costly to make.
}

// TODO: costly should return an error.
func getAllAtts(node *Node) (atts map[string]string) {
	atts = make(map[string]string)

	for k, v := range node.atts {
		atts[k] = v
	}
	for k, v := range node.costly() {
		atts[k] = v
	}
	return
}
