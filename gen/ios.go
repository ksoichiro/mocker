package gen

import (
	"fmt"
	"path/filepath"
	"sync"
)

type IosGenerator struct {
	opt  *Options
	mock *Mock
}

var iwd WidgetsDef

func defineIosWidgets() {
	iwd = WidgetsDef{}
	iwd.Add("button", Widget{
		Name:     "UIButton",
		Textable: true,
		Gravity:  GravityCenter,
		SizeW:    SizeFill,
		SizeH:    SizeWrap,
	})
	iwd.Add("label", Widget{
		Name:     "UILabel",
		Textable: true,
		Gravity:  GravityCenter,
		SizeW:    SizeFill,
		SizeH:    SizeWrap,
	})
	iwd.Add("linear", Widget{
		Textable:    false,
		Orientation: OrientationVertical,
		SizeW:       SizeFill,
		SizeH:       SizeFill,
	})
	iwd.Add("relative", Widget{
		Textable: false,
		SizeW:    SizeFill,
		SizeH:    SizeFill,
	})
}

func (g *IosGenerator) Generate() {
	defineIosWidgets()

	outDir := g.opt.OutDir
	projectDir := filepath.Join(outDir, g.mock.Meta.Ios.Project)

	var wg sync.WaitGroup

	// Generate contents.xcworkspacedata
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosContentsXcWorkspaceData(mock, dir)
	}(g.mock, outDir)

	// Generate .gitignore
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosGitignore(mock, dir)
	}(g.mock, outDir)

	// Generate main.m
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosMain(mock, dir)
	}(g.mock, outDir)

	// Generate Info.plist
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosInfoPlist(mock, dir)
	}(g.mock, outDir)

	// Generate InfoPlist.strings
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosInfoPlistStrings(mock, dir)
	}(g.mock, outDir)

	// Generate Prefix.pch
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosPch(mock, dir)
	}(g.mock, outDir)

	// Generate Images.xcassets
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosImagesXcAssetsAppIcon(mock, dir)
	}(g.mock, outDir)
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosImagesXcAssetsLaunchImage(mock, dir)
	}(g.mock, outDir)

	// Generate AppDelegate
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosAppDelegateHeader(mock, dir)
	}(g.mock, outDir)
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosAppDelegateImplementation(mock, dir)
	}(g.mock, outDir)

	// Generate ViewControllers
	for _, screen := range g.mock.Screens {
		wg.Add(1)
		go func(mock *Mock, dir string, screen Screen) {
			defer wg.Done()
			genIosViewController(mock, dir, screen)
			genIosViewControllerLayout(mock, dir, screen)
		}(g.mock, projectDir, screen)
	}

	// Generate resources
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosLocalizedStrings(mock, dir)
	}(g.mock, projectDir)
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosColors(mock, dir)
	}(g.mock, projectDir)

	// Generate project.pbxproj
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosProjectPbxproj(mock, dir)
	}(g.mock, outDir)

	wg.Wait()
}

func genIosContentsXcWorkspaceData(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosContentsXcWorkspaceData(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project+".xcodeproj", "project.xcworkspace", "contents.xcworkspacedata"))
}

func genCodeIosContentsXcWorkspaceData(mock *Mock, buf *CodeBuffer) {
	buf.add(`<?xml version="1.0" encoding="UTF-8"?>
<Workspace
   version = "1.0">
   <FileRef
      location = "self:%s.xcodeproj">
   </FileRef>
</Workspace>
`,
		mock.Meta.Ios.Project)
}

func genIosGitignore(mock *Mock, outDir string) {
	var buf CodeBuffer
	genCodeIosGitignore(mock, &buf)
	genFile(&buf, filepath.Join(outDir, ".gitignore"))
}

func genCodeIosGitignore(mock *Mock, buf *CodeBuffer) {
	buf.add(`*.xcodeproj/*
!*.xcodeproj/project.pbxproj
!*.xcworkspace/contents.xcworkspacedata
.DS_Store`)
}

func genIosMain(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosMain(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, "main.m"))
}

func genCodeIosMain(mock *Mock, buf *CodeBuffer) {
	buf.add(`#import <UIKit/UIKit.h>

#import "%sAppDelegate.h"

int main(int argc, char * argv[])
{
    @autoreleasepool {
        return UIApplicationMain(argc, argv, nil, NSStringFromClass([%sAppDelegate class]));
    }
}`,
		mock.Meta.Ios.ClassPrefix,
		mock.Meta.Ios.ClassPrefix)
}

func genIosInfoPlist(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosInfoPlist(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, mock.Meta.Ios.Project+"-Info.plist"))
}

func genCodeIosInfoPlist(mock *Mock, buf *CodeBuffer) {
	buf.add(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDevelopmentRegion</key>
	<string>en</string>
	<key>CFBundleDisplayName</key>
	<string>${PRODUCT_NAME}</string>
	<key>CFBundleExecutable</key>
	<string>${EXECUTABLE_NAME}</string>
	<key>CFBundleIdentifier</key>
	<string>%s.${PRODUCT_NAME:rfc1034identifier}</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>${PRODUCT_NAME}</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundleSignature</key>
	<string>????</string>
	<key>CFBundleVersion</key>
	<string>1.0</string>
	<key>LSRequiresIPhoneOS</key>
	<true/>
	<key>UIRequiredDeviceCapabilities</key>
	<array>
		<string>armv7</string>
	</array>
	<key>UISupportedInterfaceOrientations</key>
	<array>
		<string>UIInterfaceOrientationPortrait</string>
		<string>UIInterfaceOrientationLandscapeLeft</string>
		<string>UIInterfaceOrientationLandscapeRight</string>
	</array>
</dict>
</plist>`,
		mock.Meta.Ios.CompanyIdentifier)
}

func genIosInfoPlistStrings(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosInfoPlistStrings(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, "en.lproj", "InfoPlist.strings"))
}

func genCodeIosInfoPlistStrings(mock *Mock, buf *CodeBuffer) {
	buf.add(`/* Localized versions of Info.plist keys */`)
}

func genIosPch(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosPch(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, mock.Meta.Ios.Project+"-Prefix.pch"))
}

func genCodeIosPch(mock *Mock, buf *CodeBuffer) {
	buf.add(`//
//  Prefix header
//
//  The contents of this file are implicitly included at the beginning of every source file.
//

#import <Availability.h>

#ifndef __IPHONE_3_0
#warning "This project uses features only available in iOS SDK 3.0 and later."
#endif

#ifdef __OBJC__
    #import <UIKit/UIKit.h>
    #import <Foundation/Foundation.h>
#endif`)
}

func genIosImagesXcAssetsAppIcon(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosImagesXcAssetsAppIcon(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, "Images.xcassets", "AppIcon.appiconset", "Contents.json"))
}

func genCodeIosImagesXcAssetsAppIcon(mock *Mock, buf *CodeBuffer) {
	buf.add(`{
  "images" : [
    {
      "idiom" : "iphone",
      "size" : "29x29",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "40x40",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "60x60",
      "scale" : "2x"
    }
  ],
  "info" : {
    "version" : 1,
    "author" : "xcode"
  }
}`)
}

func genIosImagesXcAssetsLaunchImage(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosImagesXcAssetsLaunchImage(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, "Images.xcassets", "LaunchImage.launchimage", "Contents.json"))
}

func genCodeIosImagesXcAssetsLaunchImage(mock *Mock, buf *CodeBuffer) {
	buf.add(`{
  "images" : [
    {
      "orientation" : "portrait",
      "idiom" : "iphone",
      "extent" : "full-screen",
      "minimum-system-version" : "%s",
      "scale" : "2x"
    },
    {
      "orientation" : "portrait",
      "idiom" : "iphone",
      "subtype" : "retina4",
      "extent" : "full-screen",
      "minimum-system-version" : "%s",
      "scale" : "2x"
    }
  ],
  "info" : {
    "version" : 1,
    "author" : "xcode"
  }
}`,
		mock.Meta.Ios.DeploymentTarget,
		mock.Meta.Ios.DeploymentTarget)
}

func genIosAppDelegateHeader(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosAppDelegateHeader(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, mock.Meta.Ios.ClassPrefix+"AppDelegate.h"))
}

func genCodeIosAppDelegateHeader(mock *Mock, buf *CodeBuffer) {
	buf.add(`#import <UIKit/UIKit.h>

@interface %sAppDelegate : UIResponder <UIApplicationDelegate>

@property (strong, nonatomic) UIWindow *window;

@end`, mock.Meta.Ios.ClassPrefix)
}

func genIosAppDelegateImplementation(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosAppDelegateImplementation(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, mock.Meta.Ios.ClassPrefix+"AppDelegate.m"))
}

func genCodeIosAppDelegateImplementation(mock *Mock, buf *CodeBuffer) {
	buf.add(`#import "%sAppDelegate.h"

@implementation %sAppDelegate

- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions
{
    self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    self.window.backgroundColor = [UIColor whiteColor];
    [self.window makeKeyAndVisible];
    return YES;
}

@end`,
		mock.Meta.Ios.ClassPrefix,
		mock.Meta.Ios.ClassPrefix)
}

func genIosViewController(mock *Mock, dir string, screen Screen) {
	// TODO
	fmt.Println("iOS: ViewController generator: Not implemented...")
}

func genIosViewControllerLayout(mock *Mock, dir string, screen Screen) {
	// TODO
	fmt.Println("iOS: Layout generator: Not implemented...")
}

func genIosLocalizedStrings(mock *Mock, dir string) {
	// TODO
	fmt.Println("iOS: LocalizedString generator: Not implemented...")
}

func genIosColors(mock *Mock, dir string) {
	// TODO
	fmt.Println("iOS: Colors generator: Not implemented...")
}

func genIosProjectPbxproj(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosProjectPbxproj(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project+".xcodeproj", "project.pbxproj"))
}

func genCodeIosProjectPbxproj(mock *Mock, buf *CodeBuffer) {
	fileId := 0

	// Header
	buf.add(`// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 46;
	objects = {`)

	buf.add(`
/* Begin PBXBuildFile section */`)
	buf.add(`/* End PBXBuildFile section */`)

	buf.add(`
/* Begin PBXContainerItemProxy section */`)
	buf.add(`/* End PBXContainerItemProxy section */`)

	buf.add(`
/* Begin PBXFileReference section */`)
	buf.add(`/* End PBXFileReference section */`)

	buf.add(`
/* Begin PBXFrameworksBuildPhase section */`)
	buf.add(`/* End PBXFrameworksBuildPhase section */`)

	buf.add(`
/* Begin PBXGroup section */`)
	buf.add(`/* End PBXGroup section */`)

	buf.add(`
/* Begin PBXNativeTarget section */`)
	buf.add(`/* End PBXNativeTarget section */`)

	buf.add(`
/* Begin PBXProject section */`)
	buf.add(`/* End PBXProject section */`)

	buf.add(`
/* Begin PBXResourcesBuildPhase section */`)
	buf.add(`/* End PBXResourcesBuildPhase section */`)

	buf.add(`
/* Begin PBXSourcesBuildPhase section */`)
	buf.add(`/* End PBXSourcesBuildPhase section */`)

	buf.add(`
/* Begin PBXTargetDependency section */`)
	buf.add(`/* End PBXTargetDependency section */`)

	buf.add(`
/* Begin PBXVariantGroup section */`)
	buf.add(`/* End PBXVariantGroup section */`)

	buf.add(`
/* Begin XCBuildConfiguration section */`)
	buf.add(`/* End XCBuildConfiguration section */`)

	buf.add(`
/* Begin XCConfigurationList section */`)
	buf.add(`/* End XCConfigurationList section */`)

	// Footer
	buf.add(`	};
	rootObject = %s /* Project object */;
}`, genIosFileId(&fileId))
}

func genIosFileId(i *int) string {
	*i++
	return fmt.Sprintf("%024d", *i)
}
