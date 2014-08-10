package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ksoichiro/mocker/encoding/mockerfile"
	"github.com/ksoichiro/mocker/gen"
)

const (
	Version         = "0.1.0"
	ExitCodeSuccess = 0
	ExitCodeError   = 1
)

func main() {
	// Parse command
	flag.Usage = printUsage
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(ExitCodeError)
	}
	switch os.Args[1] {
	case "gen", "g":
	case "version":
		printVersion()
		os.Exit(ExitCodeSuccess)
	case "help":
		fallthrough
	default:
		printUsage()
		os.Exit(ExitCodeError)
	}

	// Gen needs platform ID
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(ExitCodeError)
	}
	genId := os.Args[2]

	// Options for gen subcommand
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var (
		outDir = fs.String("out", "out", "Output directory for generated codes.")
	)
	fs.Parse(os.Args[3:])

	mock := parseConfigs()
	opt := gen.Options{
		OutDir: *outDir,
	}
	//gen(&opt, &mock, genId)
	g := gen.NewGenerator(&opt, &mock, genId)
	if g == nil {
		fmt.Printf("Invalid gen ID: %s\n", genId)
		printUsage()
		os.Exit(ExitCodeError)
	}
	g.Generate()
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `mocker is a mock up framework for mobile apps.
Usage: %s command
Command:
  g[en]    generate source code (see 'Generator')
  help     show this help
  version  show version of mocker

Generator:
  mocker g[en] ID [options]

  ID:
    android  Java and XML code for Android app
    ios      Objective-C code for iOS app

  options:
    -out="out": Output directory for generated codes
`, os.Args[0])
}

func printVersion() {
	fmt.Println("mocker version \"" + Version + "\"")
}

func parseConfigs() (mock gen.Mock) {
	filename := "Mockerfile"
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)
	err = mockerfile.Unmarshal(b, &mock)
	if err != nil {
		fmt.Println("Error unmarshaling Mockerfile", err)
		return
	}

	return
}
