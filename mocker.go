package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	Name    string
	Screens []Screen
	Launch  Launch
	Colors  []Color
	Strings []String
}

type Screen struct {
	Id        string
	Name      string
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
	srcDir := filepath.Join(outDir, "src")
	mainDir := filepath.Join(srcDir, "main")
	javaDir := filepath.Join(mainDir, "java")
	packageDir := filepath.Join(javaDir, strings.Replace(mock.Package, ".", string(os.PathSeparator), -1))
	resDir := filepath.Join(mainDir, "res")
	layoutDir := filepath.Join(resDir, "layout")
	valuesDir := filepath.Join(resDir, "values")

	// Generate base file set using android command
	cmd := exec.Command("android", "create", "project", "-n", "mock", "-v", "0.12.+", "-g", "-k", mock.Package, "-a", "DummyActivity", "-t", "android-19", "-p", outDir)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error while generating project")
	}

	// Remove unecessery directories and files
	os.RemoveAll(filepath.Join(srcDir, "androidTest"))
	os.Remove(filepath.Join(packageDir, "DummyActivity.java"))
	os.Remove(filepath.Join(layoutDir, "main.xml"))

	// Generate Manifest
	genAndroidManifest(mock, mainDir)

	// Generate Activities
	for i := range mock.Screens {
		screen := mock.Screens[i]
		genAndroidActivity(mock, packageDir, screen)
		genAndroidActivityLayout(mock, layoutDir, screen)
	}

	// Generate resources
	genAndroidStrings(mock, valuesDir)
	genAndroidLocalizedStrings(mock, resDir)
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

func genAndroidActivityLayout(mock *Mock, layoutDir string, screen Screen) {
	filename := filepath.Join(layoutDir, "activity_"+screen.Id+".xml")
	f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	xmlns := `xmlns:android="http://schemas.android.com/apk/res/android"`
	f.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<LinearLayout %s
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    android:background="@android:color/white"
    android:gravity="center"
    android:orientation="vertical"
    android:padding="16dp" >
`, xmlns))

	for i := range screen.Layout.Views {
		view := screen.Layout.Views[i]
		switch view.Type {
		case "button":
			f.WriteString(fmt.Sprintf(`
    <Button
        android:id="@+id/%s"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="@string/%s"
        android:textColor="#FF8800"
        android:textSize="30dp" />
`, view.Id, view.Label))
		default:
		}
	}
	f.WriteString(`</LinearLayout>
`)

	f.Close()
}

func genAndroidStrings(mock *Mock, valuesDir string) {
	filename := filepath.Join(valuesDir, "strings_app.xml")
	f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	// App name
	f.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">%s</app_name>
`, mock.Name))

	// Activity title
	for i := range mock.Screens {
		screen := mock.Screens[i]
		f.WriteString(fmt.Sprintf(`    <string name="activity_title_%s">%s</string>
`, screen.Id, screen.Name))
	}

	f.WriteString(`</resources>
`)
	f.Close()
}

func genAndroidLocalizedStrings(mock *Mock, resDir string) {
	for i := range mock.Strings {
		s := mock.Strings[i]
		lang := s.Lang
		suffix := "-" + lang
		if strings.ToLower(lang) == "base" {
			suffix = ""
		}
		valuesDir := filepath.Join(resDir, "values"+suffix)
		os.MkdirAll(valuesDir, 0777)
		filename := filepath.Join(valuesDir, "strings.xml")
		f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
		defer f.Close()
		f.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<resources>
`)
		for j := range s.Defs {
			def := s.Defs[j]
			f.WriteString(fmt.Sprintf(`    <string name="%s">%s</string>
`, def.Id, def.Value))
		}
		f.WriteString(`</resources>
`)
		f.Close()
	}
}
