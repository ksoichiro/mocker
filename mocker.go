package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	Version         = "0.1.0"
	ExitCodeSuccess = 0
	ExitCodeError   = 1
)

type Mock struct {
	Package string
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

func genAndroid(mock *Mock) {
	// TODO
	fmt.Printf("%+v\n", mock)
	outDir := "out"
	os.MkdirAll(outDir, 0777)

	genAndroidManifest(mock, outDir)

	javaDir := filepath.Join(outDir, "app", "src", "main", "java")
	os.MkdirAll(javaDir, 0777)
	packageDir := filepath.Join(javaDir, strings.Replace(mock.Package, ".", string(os.PathSeparator), -1))
	os.MkdirAll(packageDir, 0777)

	for i := range mock.Screens {
		screen := mock.Screens[i]
		genAndroidActivity(mock, packageDir, screen)
	}
}

func genIos(mock *Mock) {
	// TODO
	fmt.Printf("%+v\n", mock)
}

func genAndroidManifest(mock *Mock, outDir string) {
	filename := filepath.Join(outDir, "AndroidManifest.xml")
	f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	f.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="%s" >

    <application
        android:allowBackup="true"
        android:icon="@drawable/ic_launcher"
        android:label="@string/app_name"
        android:theme="@style/AppTheme" >
`, mock.Package))

	launcherId := mock.Launch.Screen
	for i := range mock.Screens {
		screen := mock.Screens[i]
		activityId := strings.Title(screen.Id)
		if screen.Id == launcherId {
			// Launcher
			f.WriteString(fmt.Sprintf(`        <activity
            android:label="@string/activity_title_%s"
            android:name=".%sActivity" >
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />

                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
`, screen.Id, activityId))
		} else {
			f.WriteString(fmt.Sprintf(`        <activity
            android:label="@string/activity_title_%s"
            android:name=".%sActivity" />
`, screen.Id, activityId))
		}
	}

	f.WriteString(`    </application>

</manifest>
`)
	f.Close()
}

func genAndroidActivity(mock *Mock, packageDir string, screen Screen) {
	activityId := strings.Title(screen.Id)
	filename := filepath.Join(packageDir, activityId+"Activity.java")
	f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	f.WriteString(fmt.Sprintf(`package %s;

import android.os.Bundle;

public class %sActivity extends Activity {

	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_%s);
		init();
	}

	private void init() {
	}

}
`, mock.Package, activityId, screen.Id))
	f.Close()
}
