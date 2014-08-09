package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	Version         = "0.1.0"
	ExitCodeSuccess = 0
	ExitCodeError   = 1
)

type Options struct {
	OutDir string
}

type Mock struct {
	Name    string
	Meta    Meta
	Screens []Screen
	Launch  Launch
	Colors  []Color
	Strings []String
}

type Meta struct {
	Android Android
}

type Android struct {
	Package             string
	GradlePluginVersion string `json:"gradle_plugin_version"`
	BuildToolsVersion   string `json:"build_tools_version"`
	MinSdkVersion       int    `json:"min_sdk_version"`
	TargetSdkVersion    int    `json:"target_sdk_version"`
	CompileSdkVersion   string `json:"compile_sdk_version"`
	VersionCode         int    `json:"version_code"`
	VersionName         string `json:"version_name"`
}

type Screen struct {
	Id        string
	Name      string
	Layout    []View
	Behaviors []Behavior
}

type LayoutOptions struct {
	Width  string
	Height string
}

type View struct {
	Id      string
	Type    string
	Sub     []View
	Label   string
	Gravity string
	Below   string
	SizeW   string `json:"size_w"`
	SizeH   string `json:"size_h"`
	AlignH  string `json:"align_h"`
	AlignV  string `json:"align_v"`
}

type Behavior struct {
	Trigger Trigger
	Action  Action
}

type Trigger struct {
	Type   string
	Widget string
}

type Action struct {
	Type    string
	Transit string
}

type Launch struct {
	Screen string
}

type Color struct {
	Id    string
	Value string
}

type String struct {
	Lang string
	Defs []Def
}

type Def struct {
	Id    string
	Value string
}

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
	opt := Options{
		OutDir: *outDir,
	}
	gen(&opt, &mock, genId)
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

func parseConfigs() (mock Mock) {
	filename := "Mockerfile"
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file", err)
		return
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)
	err = json.Unmarshal(b, &mock)
	if err != nil {
		fmt.Println("Error unmarshaling Mockerfile", err)
		return
	}

	return
}
