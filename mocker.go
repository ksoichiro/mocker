package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	Version         = "0.1.0"
	ExitCodeSuccess = 0
	ExitCodeError   = 1
)

type Mock struct {
	Screens []Screen
	Launch  Launch
	Colors  []Color
	Strings []String
}

type Screen struct {
	Id        string
	Layout    Layout
	Behaviors []Behavior
}

type Layout struct {
	Views []View
}

type View struct {
	Id     string
	Type   string
	Label  string
	SizeW  string `json:"size_w"`
	SizeH  string `json:"size_h"`
	AlignH string `json:"align_h"`
	AlignV string `json:"align_v"`
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
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(ExitCodeError)
	}
	switch os.Args[1] {
	case "gen":
	case "version":
		printVersion()
		os.Exit(ExitCodeSuccess)
	case "help":
		fallthrough
	default:
		printUsage()
		os.Exit(ExitCodeError)
	}

	if len(os.Args) < 3 {
		printUsage()
		os.Exit(ExitCodeError)
	}
	genId := os.Args[2]

	mock := parseConfigs()
	gen(&mock, genId)
}

func printUsage() {
	fmt.Println(`mocker is a mock up framework for mobile apps.

Usage:
	mocker command [options]

Command:
	gen      generate source code (see 'Generator')
	help     show this help
	version  show version of mocker

Generator:
	mocker gen ID

	ID:
		ios  Objective-C code for iOS app
`)
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

func gen(mock *Mock, genId string) {
	ids := []string{
		"ios",
		"android",
	}

	var validId string
	for i := range ids {
		if genId == ids[i] {
			validId = genId
			break
		}
	}
	if validId == "" {
		fmt.Println("Invalid gen ID")
		printUsage()
		os.Exit(ExitCodeError)
	}

	// TODO
	fmt.Printf("%+v\n", mock)
}
