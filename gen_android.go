package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Default layout params for Android's widgets
type AndroidWidget struct {
	Name        string
	Textable    bool
	Gravity     string
	Orientation string
	SizeW       string
	SizeH       string
}

type AndroidWidgetsDef struct {
	widgets map[string]AndroidWidget
}

func (d *AndroidWidgetsDef) Add(name string, w AndroidWidget) {
	if d.widgets == nil {
		d.widgets = make(map[string]AndroidWidget)
	}
	d.widgets[name] = w
}

func (d *AndroidWidgetsDef) Has(name string) (ret bool) {
	ret = false
	if _, ok := d.widgets[name]; ok {
		ret = true
	}
	return
}

func (d *AndroidWidgetsDef) Get(name string) AndroidWidget {
	return d.widgets[name]
}

var awd AndroidWidgetsDef

const (
	gravityCenter       = "center"
	gravityCenterV      = "center_v"
	sizeFill            = "fill"
	sizeWrap            = "wrap"
	orientationVertical = "vertical"
)

func defineAndroidWidgets() {
	awd = AndroidWidgetsDef{}
	awd.Add("button", AndroidWidget{
		Name:     "Button",
		Textable: true,
		Gravity:  gravityCenter,
		SizeW:    sizeFill,
		SizeH:    sizeWrap,
	})
	awd.Add("label", AndroidWidget{
		Name:     "TextView",
		Textable: true,
		Gravity:  gravityCenter,
		SizeW:    sizeFill,
		SizeH:    sizeWrap,
	})
	awd.Add("linear", AndroidWidget{
		Name:        "LinearLayout",
		Textable:    false,
		Orientation: orientationVertical,
		SizeW:       sizeFill,
		SizeH:       sizeFill,
	})
	awd.Add("relative", AndroidWidget{
		Name:     "RelativeLayout",
		Textable: false,
		SizeW:    sizeFill,
		SizeH:    sizeFill,
	})
}

func genAndroid(opt *Options, mock *Mock) {
	defineAndroidWidgets()

	outDir := opt.OutDir
	srcDir := filepath.Join(outDir, "src")
	mainDir := filepath.Join(srcDir, "main")
	javaDir := filepath.Join(mainDir, "java")
	packageDir := filepath.Join(javaDir, strings.Replace(mock.Meta.Android.Package, ".", string(os.PathSeparator), -1))
	resDir := filepath.Join(mainDir, "res")
	layoutDir := filepath.Join(resDir, "layout")
	valuesDir := filepath.Join(resDir, "values")

	// Generate base file set using android command
	cmd := exec.Command("android", "create", "project",
		"-n", "mock",
		"-v", mock.Meta.Android.GradlePluginVersion,
		"-g",
		"-k", mock.Meta.Android.Package,
		"-a", "DummyActivity",
		"-t", "android-19",
		"-p", outDir)
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

	// Generate build.gradle
	genAndroidGradle(mock, outDir)

	// Generate Activities
	for i := range mock.Screens {
		screen := mock.Screens[i]
		genAndroidActivity(mock, packageDir, screen)
		genAndroidActivityLayout(mock, layoutDir, screen)
	}

	// Generate resources
	genAndroidStrings(mock, valuesDir)
	genAndroidLocalizedStrings(mock, resDir)
	genAndroidColors(mock, valuesDir)
	genAndroidStyles(mock, valuesDir)
}

func genAndroidManifest(mock *Mock, outDir string) {
	filename := filepath.Join(outDir, "AndroidManifest.xml")
	f := createFile(filename)
	defer f.Close()
	var buf []string
	genCodeAndroidManifest(mock, &buf)
	for _, s := range buf {
		f.WriteString(s + "\n")
	}
	f.Close()
}

func genCodeAndroidManifest(mock *Mock, buf *[]string) {
	*buf = append(*buf, fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="%s" >

    <application
        android:allowBackup="true"
        android:icon="@drawable/ic_launcher"
        android:label="@string/app_name"
        android:theme="@style/AppTheme" >`, mock.Meta.Android.Package))

	launcherId := mock.Launch.Screen
	for i := range mock.Screens {
		screen := mock.Screens[i]
		activityId := strings.Title(screen.Id)
		if screen.Id == launcherId {
			// Launcher
			*buf = append(*buf, fmt.Sprintf(`        <activity
            android:label="@string/activity_title_%s"
            android:name=".%sActivity" >
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />

                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>`, screen.Id, activityId))
		} else {
			*buf = append(*buf, fmt.Sprintf(`        <activity
            android:label="@string/activity_title_%s"
            android:name=".%sActivity" />`, screen.Id, activityId))
		}
	}

	*buf = append(*buf, `    </application>
</manifest>`)
}

func genAndroidGradle(mock *Mock, outDir string) {
	filename := filepath.Join(outDir, "build.gradle")
	f := createFile(filename)
	defer f.Close()
	var buf []string
	genCodeAndroidGradle(mock, &buf)
	for _, s := range buf {
		f.WriteString(s + "\n")
	}
	f.Close()
}

func genCodeAndroidGradle(mock *Mock, buf *[]string) {
	*buf = append(*buf, fmt.Sprintf(`buildscript {
    repositories {
        mavenCentral()
    }
    dependencies {
        classpath 'com.android.tools.build:gradle:%s'
    }
}`, mock.Meta.Android.GradlePluginVersion))
	*buf = append(*buf, fmt.Sprintf(`apply plugin: 'com.android.application'

android {
    compileSdkVersion '%s'
    buildToolsVersion '%s'

    defaultConfig {
        applicationId "%s"
        minSdkVersion %d
        targetSdkVersion %d
        versionCode %d
        versionName "%s"
    }

    buildTypes {
        release {
            runProguard false
            proguardFile getDefaultProguardFile('proguard-android.txt')
        }
    }

    lintOptions {
        checkReleaseBuilds false
        abortOnError false
    }
}`,
		mock.Meta.Android.CompileSdkVersion,
		mock.Meta.Android.BuildToolsVersion,
		mock.Meta.Android.Package,
		mock.Meta.Android.MinSdkVersion,
		mock.Meta.Android.TargetSdkVersion,
		mock.Meta.Android.VersionCode,
		mock.Meta.Android.VersionName))
}

func genAndroidActivity(mock *Mock, packageDir string, screen Screen) {
	activityId := strings.Title(screen.Id)
	filename := filepath.Join(packageDir, activityId+"Activity.java")
	f := createFile(filename)
	defer f.Close()
	var buf []string
	genCodeAndroidActivity(mock, screen, &buf)
	for _, s := range buf {
		f.WriteString(s + "\n")
	}
	f.Close()
}

func genCodeAndroidActivity(mock *Mock, screen Screen, buf *[]string) {
	activityId := strings.Title(screen.Id)
	*buf = append(*buf, fmt.Sprintf(`package %s;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.view.View;

public class %sActivity extends Activity {

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_%s);
        init();
    }

    private void init() {`,
		mock.Meta.Android.Package, activityId, screen.Id))

	for i := range screen.Behaviors {
		b := screen.Behaviors[i]
		if b.Trigger.Type != "click" {
			// Not support other than click currently
			continue
		}
		*buf = append(*buf, fmt.Sprintf(`        findViewById(R.id.%s).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {`, b.Trigger.Widget))

		if b.Action.Type == "transit_forward" {
			var id string
			for j := range mock.Screens {
				next := mock.Screens[j]
				if next.Id == b.Action.Transit {
					id = next.Id
				}
			}
			*buf = append(*buf, fmt.Sprintf(`                startActivity(new Intent(%sActivity.this, %sActivity.class));`,
				strings.Title(screen.Id),
				strings.Title(id)))
		}

		*buf = append(*buf, `            }
        });`)
	}

	*buf = append(*buf, `    }

}`)
}

func genAndroidActivityLayout(mock *Mock, layoutDir string, screen Screen) {
	filename := filepath.Join(layoutDir, "activity_"+screen.Id+".xml")
	f := createFile(filename)
	defer f.Close()
	var buf []string
	genCodeAndroidActivityLayout(mock, screen, &buf)
	for _, s := range buf {
		f.WriteString(s + "\n")
	}
	f.Close()
}

func genCodeAndroidActivityLayout(mock *Mock, screen Screen, buf *[]string) {
	*buf = append(*buf, `<?xml version="1.0" encoding="utf-8"?>`)
	if 0 < len(screen.Layout) {
		// Only parse root view
		genAndroidLayoutRecur(&screen.Layout[0], true, buf)
	}
}

func genAndroidLayoutRecur(view *View, top bool, buf *[]string) {
	if !awd.Has(view.Type) {
		return
	}
	widget := awd.Get(view.Type)

	xmlns := ""
	if top {
		xmlns = ` xmlns:android="http://schemas.android.com/apk/res/android"`
	}

	lo := convertAndroidLayoutOptions(widget, view)
	hasSub := 0 < len(view.Sub)

	*buf = append(*buf, fmt.Sprintf(`<%s%s`, widget.Name, xmlns))
	if view.Id != "" {
		*buf = append(*buf, fmt.Sprintf(`    android:id="@+id/%s"`, view.Id))
	}
	if view.Below != "" {
		*buf = append(*buf, fmt.Sprintf(`    android:layout_below="@id/%s"`, view.Below))
	}
	if widget.Textable && view.Label != "" {
		*buf = append(*buf, fmt.Sprintf(`    android:text="@string/%s"`, view.Label))
	}
	if widget.Orientation != "" {
		*buf = append(*buf, fmt.Sprintf(`    android:orientation="%s"`, widget.Orientation))
	}
	if view.Gravity != "" {
		gravity := ""
		switch view.Gravity {
		case gravityCenter:
			gravity = "center"
		case gravityCenterV:
			gravity = "center_vertical"
		}
		*buf = append(*buf, fmt.Sprintf(`    android:gravity="%s"`, gravity))
	} else if widget.Gravity != "" {
		*buf = append(*buf, fmt.Sprintf(`    android:gravity="%s"`, widget.Gravity))
	}
	*buf = append(*buf, fmt.Sprintf(`    android:layout_width="%s"
    android:layout_height="%s"`,
		lo.Width,
		lo.Height))

	if hasSub {
		// Print sub views recursively
		*buf = append(*buf, `    >`)
		for _, sv := range view.Sub {
			genAndroidLayoutRecur(&sv, false, buf)
		}
		*buf = append(*buf, fmt.Sprintf(`</%s>`, widget.Name))
	} else {
		*buf = append(*buf, `    />`)
	}
}

func genAndroidStrings(mock *Mock, valuesDir string) {
	filename := filepath.Join(valuesDir, "strings_app.xml")
	f := createFile(filename)
	defer f.Close()
	var buf []string
	genCodeAndroidStrings(mock, &buf)
	for _, s := range buf {
		f.WriteString(s + "\n")
	}
	f.Close()
}

func genCodeAndroidStrings(mock *Mock, buf *[]string) {
	// App name
	*buf = append(*buf, fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">%s</string>`, mock.Name))

	// Activity title
	for i := range mock.Screens {
		screen := mock.Screens[i]
		*buf = append(*buf, fmt.Sprintf(`    <string name="activity_title_%s">%s</string>`,
			screen.Id, screen.Name))
	}

	*buf = append(*buf, `</resources>`)
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
		f := createFile(filename)
		defer f.Close()
		var buf []string
		genCodeAndroidLocalizedStrings(s, &buf)
		for _, s := range buf {
			f.WriteString(s + "\n")
		}
		f.Close()
	}
}

func genCodeAndroidLocalizedStrings(s String, buf *[]string) {
	*buf = append(*buf, `<?xml version="1.0" encoding="utf-8"?>
<resources>`)
	for j := range s.Defs {
		def := s.Defs[j]
		*buf = append(*buf, fmt.Sprintf(`    <string name="%s">%s</string>`, def.Id, def.Value))
	}
	*buf = append(*buf, `</resources>`)
}

func genAndroidColors(mock *Mock, valuesDir string) {
	filename := filepath.Join(valuesDir, "colors.xml")
	f := createFile(filename)
	defer f.Close()
	var buf []string
	genCodeAndroidColors(mock, &buf)
	for _, s := range buf {
		f.WriteString(s + "\n")
	}
	f.Close()
}

func genCodeAndroidColors(mock *Mock, buf *[]string) {
	*buf = append(*buf, `<?xml version="1.0" encoding="utf-8"?>
<resources>`)

	for i := range mock.Colors {
		c := mock.Colors[i]
		*buf = append(*buf, fmt.Sprintf(`    <color name="%s">%s</color>`, c.Id, c.Value))
	}

	*buf = append(*buf, `</resources>`)
}

func genAndroidStyles(mock *Mock, valuesDir string) {
	filename := filepath.Join(valuesDir, "styles.xml")
	f := createFile(filename)
	defer f.Close()

	f.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <style name="AppTheme" parent="android:Theme.Holo.Light.DarkActionBar">
    </style>
</resources>
`)

	f.Close()
}

func convertAndroidLayoutOptions(widget AndroidWidget, view *View) (lo LayoutOptions) {
	base := view.SizeW
	if base == "" {
		base = widget.SizeW
	}
	if base == sizeFill {
		lo.Width = "match_parent"
	} else {
		lo.Width = "wrap_content"
	}
	base = view.SizeH
	if base == "" {
		base = widget.SizeH
	}
	if base == sizeFill {
		lo.Height = "match_parent"
	} else {
		lo.Height = "wrap_content"
	}
	return
}
