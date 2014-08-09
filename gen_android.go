package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func genAndroid(mock *Mock) {
	// TODO
	fmt.Printf("%+v\n", mock)

	outDir := "out"
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
	f.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="%s" >

    <application
        android:allowBackup="true"
        android:icon="@drawable/ic_launcher"
        android:label="@string/app_name"
        android:theme="@style/AppTheme" >
`, mock.Meta.Android.Package))

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

func genAndroidGradle(mock *Mock, outDir string) {
	filename := filepath.Join(outDir, "build.gradle")
	f := createFile(filename)
	defer f.Close()
	f.WriteString(fmt.Sprintf(`buildscript {
    repositories {
        mavenCentral()
    }
    dependencies {
        classpath 'com.android.tools.build:gradle:%s'
    }
}
`, mock.Meta.Android.GradlePluginVersion))
	f.WriteString(fmt.Sprintf(`apply plugin: 'com.android.application'

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
}
`,
		mock.Meta.Android.CompileSdkVersion,
		mock.Meta.Android.BuildToolsVersion,
		mock.Meta.Android.Package,
		mock.Meta.Android.MinSdkVersion,
		mock.Meta.Android.TargetSdkVersion,
		mock.Meta.Android.VersionCode,
		mock.Meta.Android.VersionName))

	f.Close()
}

func genAndroidActivity(mock *Mock, packageDir string, screen Screen) {
	activityId := strings.Title(screen.Id)
	filename := filepath.Join(packageDir, activityId+"Activity.java")
	f := createFile(filename)
	defer f.Close()
	f.WriteString(fmt.Sprintf(`package %s;

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

    private void init() {
`, mock.Meta.Android.Package, activityId, screen.Id))

	for i := range screen.Behaviors {
		b := screen.Behaviors[i]
		if b.Trigger.Type != "click" {
			// Not support other than click currently
			continue
		}
		f.WriteString(fmt.Sprintf(`        findViewById(R.id.%s).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
`, b.Trigger.Widget))

		if b.Action.Type == "transit_forward" {
			var id string
			for j := range mock.Screens {
				next := mock.Screens[j]
				if next.Id == b.Action.Transit {
					id = next.Id
				}
			}
			f.WriteString(fmt.Sprintf(`                startActivity(new Intent(%sActivity.this, %sActivity.class));
`,
				strings.Title(screen.Id),
				strings.Title(id)))
		}

		f.WriteString(`            }
        });
`)
	}

	f.WriteString(`    }

}
`)
	f.Close()
}

func genAndroidActivityLayout(mock *Mock, layoutDir string, screen Screen) {
	filename := filepath.Join(layoutDir, "activity_"+screen.Id+".xml")
	f := createFile(filename)
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
        android:text="@string/%s" />
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
	f := createFile(filename)
	defer f.Close()

	// App name
	f.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">%s</string>
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
		f := createFile(filename)
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

func genAndroidColors(mock *Mock, valuesDir string) {
	filename := filepath.Join(valuesDir, "colors.xml")
	f := createFile(filename)
	defer f.Close()

	f.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<resources>
`)

	for i := range mock.Colors {
		c := mock.Colors[i]
		f.WriteString(fmt.Sprintf(`    <color name="%s">%s</color>
`, c.Id, c.Value))
	}

	f.WriteString(`</resources>
`)
	f.Close()
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
