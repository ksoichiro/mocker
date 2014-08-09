package main

import (
	"fmt"
	"os"
)

func gen(opt *Options, mock *Mock, genId string) {
	switch genId {
	case "ios":
		genIos(opt, mock)
	case "android":
		genAndroid(opt, mock)
	default:
		fmt.Printf("Invalid gen ID: %s\n", genId)
		printUsage()
		os.Exit(ExitCodeError)
	}
}
