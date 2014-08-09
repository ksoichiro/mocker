package main

import (
	"fmt"
	"os"
)

func gen(mock *Mock, genId string) {
	switch genId {
	case "ios":
		genIos(mock)
	case "android":
		genAndroid(mock)
	default:
		fmt.Println("Invalid gen ID")
		printUsage()
		os.Exit(ExitCodeError)
	}
}
