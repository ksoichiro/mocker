package gen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type AndroidGenerator struct {
	opt  *Options
	mock *Mock
}

var awd WidgetsDef

func defineAndroidWidgets() {
	awd = WidgetsDef{}
	awd.Add("button", Widget{
		Name:     "Button",
		Textable: true,
		Gravity:  GravityCenter,
		SizeW:    SizeFill,
		SizeH:    SizeWrap,
	})
	awd.Add("label", Widget{
		Name:     "TextView",
		Textable: true,
		Gravity:  GravityCenter,
		SizeW:    SizeFill,
		SizeH:    SizeWrap,
	})
	awd.Add("linear", Widget{
		Name:        "LinearLayout",
		Textable:    false,
		Orientation: OrientationVertical,
		SizeW:       SizeFill,
		SizeH:       SizeFill,
	})
	awd.Add("relative", Widget{
		Name:     "RelativeLayout",
		Textable: false,
		SizeW:    SizeFill,
		SizeH:    SizeFill,
	})
}

func (g *AndroidGenerator) Generate() {
	defineAndroidWidgets()

	outDir := g.opt.OutDir
	srcDir := filepath.Join(outDir, "src")
	mainDir := filepath.Join(srcDir, "main")
	javaDir := filepath.Join(mainDir, "java")
	packageDir := filepath.Join(javaDir, strings.Replace(g.mock.Meta.Android.Package, ".", string(os.PathSeparator), -1))
	resDir := filepath.Join(mainDir, "res")
	layoutDir := filepath.Join(resDir, "layout")
	valuesDir := filepath.Join(resDir, "values")

	// Generate base file set using android command
	cmd := exec.Command("android", "create", "project",
		"-n", "mock",
		"-v", g.mock.Meta.Android.GradlePluginVersion,
		"-g",
		"-k", g.mock.Meta.Android.Package,
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

	var wg sync.WaitGroup

	// Generate Manifest
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genAndroidManifest(mock, dir)
	}(g.mock, mainDir)

	// Generate build.gradle
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genAndroidGradle(mock, dir)
	}(g.mock, outDir)

	// Generate Activities
	for _, screen := range g.mock.Screens {
		wg.Add(1)
		go func(mock *Mock, dir1, dir2 string, screen Screen) {
			defer wg.Done()
			genAndroidActivity(mock, dir1, screen)
			genAndroidActivityLayout(mock, dir2, screen)
		}(g.mock, packageDir, layoutDir, screen)
	}

	// Generate resources
	wg.Add(1)
	go func(mock *Mock, dir1, dir2 string) {
		defer wg.Done()
		genAndroidStrings(mock, dir1)
		genAndroidLocalizedStrings(mock, dir2)
	}(g.mock, valuesDir, resDir)
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genAndroidColors(mock, dir)
	}(g.mock, valuesDir)
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genAndroidStyles(mock, dir)
	}(g.mock, valuesDir)

	wg.Wait()
}

func genAndroidManifest(mock *Mock, outDir string) {
	var buf CodeBuffer
	genCodeAndroidManifest(mock, &buf)
	genFile(&buf, filepath.Join(outDir, "AndroidManifest.xml"))
}

func genCodeAndroidManifest(mock *Mock, buf *CodeBuffer) {
	buf.add(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="%s" >

    <application
        android:allowBackup="true"
        android:icon="@drawable/ic_launcher"
        android:label="@string/app_name"
        android:theme="@style/AppTheme" >`, mock.Meta.Android.Package)

	launcherId := mock.Launch.Screen
	for _, screen := range mock.Screens {
		activityId := strings.Title(screen.Id)
		if screen.Id == launcherId {
			// Launcher
			buf.add(`        <activity
            android:label="@string/activity_title_%s"
            android:name=".%sActivity" >
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />

                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>`, screen.Id, activityId)
		} else {
			buf.add(`        <activity
            android:label="@string/activity_title_%s"
            android:name=".%sActivity" />`, screen.Id, activityId)
		}
	}

	buf.add(`    </application>
</manifest>`)
}

func genAndroidGradle(mock *Mock, outDir string) {
	var buf CodeBuffer
	genCodeAndroidGradle(mock, &buf)
	genFile(&buf, filepath.Join(outDir, "build.gradle"))
}

func genCodeAndroidGradle(mock *Mock, buf *CodeBuffer) {
	buf.add(`buildscript {
    repositories {
        mavenCentral()
    }
    dependencies {
        classpath 'com.android.tools.build:gradle:%s'
    }
}`, mock.Meta.Android.GradlePluginVersion)
	buf.add(`apply plugin: 'com.android.application'

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
		mock.Meta.Android.VersionName)
}

func genAndroidActivity(mock *Mock, packageDir string, screen Screen) {
	var buf CodeBuffer
	genCodeAndroidActivity(mock, screen, &buf)
	genFile(&buf, filepath.Join(packageDir, strings.Title(screen.Id)+"Activity.java"))
}

func genCodeAndroidActivity(mock *Mock, screen Screen, buf *CodeBuffer) {
	activityId := strings.Title(screen.Id)
	buf.add(`package %s;

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
		mock.Meta.Android.Package, activityId, screen.Id)

	for _, b := range screen.Behaviors {
		if b.Trigger.Type != "click" {
			// Not support other than click currently
			continue
		}
		buf.add(`        findViewById(R.id.%s).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {`, b.Trigger.Widget)

		if b.Action.Type == "transit_forward" {
			var id string
			for _, next := range mock.Screens {
				if next.Id == b.Action.Transit {
					id = next.Id
				}
			}
			buf.add(`                startActivity(new Intent(%sActivity.this, %sActivity.class));`,
				strings.Title(screen.Id),
				strings.Title(id))
		}

		buf.add(`            }
        });`)
	}

	buf.add(`    }

}`)
}

func genAndroidActivityLayout(mock *Mock, layoutDir string, screen Screen) {
	var buf CodeBuffer
	genCodeAndroidActivityLayout(mock, screen, &buf)
	genFile(&buf, filepath.Join(layoutDir, "activity_"+screen.Id+".xml"))
}

func genCodeAndroidActivityLayout(mock *Mock, screen Screen, buf *CodeBuffer) {
	buf.add(`<?xml version="1.0" encoding="utf-8"?>`)
	if 0 < len(screen.Layout) {
		// Only parse root view
		genAndroidLayoutRecur(&screen.Layout[0], true, buf)
	}
}

func genAndroidLayoutRecur(view *View, top bool, buf *CodeBuffer) {
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

	buf.add(`<%s%s`, widget.Name, xmlns)
	if view.Id != "" {
		buf.add(`    android:id="@+id/%s"`, view.Id)
	}
	if view.Below != "" {
		buf.add(`    android:layout_below="@id/%s"`, view.Below)
	}
	if widget.Textable && view.Label != "" {
		buf.add(`    android:text="@string/%s"`, view.Label)
	}
	if widget.Orientation != "" {
		buf.add(`    android:orientation="%s"`, widget.Orientation)
	}
	if view.Gravity != "" {
		gravity := ""
		switch view.Gravity {
		case GravityCenter:
			gravity = "center"
		case GravityCenterV:
			gravity = "center_vertical"
		}
		buf.add(`    android:gravity="%s"`, gravity)
	} else if widget.Gravity != "" {
		buf.add(`    android:gravity="%s"`, widget.Gravity)
	}
	if view.Margin != "" {
		if view.Margin == "normal" {
			buf.add(`    android:layout_margin="%s"`, "16dp")
		} else {
			buf.add(`    android:layout_margin="%s"`, view.Margin)
		}
	}
	buf.add(`    android:layout_width="%s"
    android:layout_height="%s"`,
		lo.Width,
		lo.Height)

	if hasSub {
		// Print sub views recursively
		buf.add(`    >`)
		for _, sv := range view.Sub {
			genAndroidLayoutRecur(&sv, false, buf)
		}
		buf.add(`</%s>`, widget.Name)
	} else {
		buf.add(`    />`)
	}
}

func genAndroidStrings(mock *Mock, valuesDir string) {
	var buf CodeBuffer
	genCodeAndroidStrings(mock, &buf)
	genFile(&buf, filepath.Join(valuesDir, "strings_app.xml"))
}

func genCodeAndroidStrings(mock *Mock, buf *CodeBuffer) {
	// App name
	buf.add(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">%s</string>`, mock.Name)

	// Activity title
	for _, screen := range mock.Screens {
		buf.add(`    <string name="activity_title_%s">%s</string>`,
			screen.Id, screen.Name)
	}

	buf.add(`</resources>`)
}

func genAndroidLocalizedStrings(mock *Mock, resDir string) {
	for _, s := range mock.Strings {
		lang := s.Lang
		suffix := "-" + lang
		if strings.ToLower(lang) == "base" {
			suffix = ""
		}
		valuesDir := filepath.Join(resDir, "values"+suffix)
		var buf CodeBuffer
		genCodeAndroidLocalizedStrings(s, &buf)
		genFile(&buf, filepath.Join(valuesDir, "strings.xml"))
	}
}

func genCodeAndroidLocalizedStrings(s String, buf *CodeBuffer) {
	buf.add(`<?xml version="1.0" encoding="utf-8"?>
<resources>`)
	for _, def := range s.Defs {
		buf.add(`    <string name="%s">%s</string>`, def.Id, def.Value)
	}
	buf.add(`</resources>`)
}

func genAndroidColors(mock *Mock, valuesDir string) {
	var buf CodeBuffer
	genCodeAndroidColors(mock, &buf)
	genFile(&buf, filepath.Join(valuesDir, "colors.xml"))
}

func genCodeAndroidColors(mock *Mock, buf *CodeBuffer) {
	buf.add(`<?xml version="1.0" encoding="utf-8"?>
<resources>`)

	for _, c := range mock.Colors {
		buf.add(`    <color name="%s">%s</color>`, c.Id, c.Value)
	}

	buf.add(`</resources>`)
}

func genAndroidStyles(mock *Mock, valuesDir string) {
	var buf CodeBuffer
	genCodeAndroidStyles(mock, &buf)
	genFile(&buf, filepath.Join(valuesDir, "styles.xml"))
}

func genCodeAndroidStyles(mock *Mock, buf *CodeBuffer) {
	buf.add(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <style name="AppTheme" parent="android:Theme.Holo.Light.DarkActionBar">
    </style>
</resources>`)
}

func convertAndroidLayoutOptions(widget Widget, view *View) (lo LayoutOptions) {
	base := view.SizeW
	if base == "" {
		base = widget.SizeW
	}
	if base == SizeFill {
		lo.Width = "match_parent"
	} else {
		lo.Width = "wrap_content"
	}
	base = view.SizeH
	if base == "" {
		base = widget.SizeH
	}
	if base == SizeFill {
		lo.Height = "match_parent"
	} else {
		lo.Height = "wrap_content"
	}
	return
}
