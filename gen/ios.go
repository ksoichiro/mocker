package gen

import (
	"fmt"
	"path/filepath"
	"strings"
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

type pbxObject struct {
	Name                   string
	Id                     string
	Location               string
	FileRef                string
	ExplicitFileType       string
	LastKnownFileType      string
	IncludeInIndex         string
	ShowNameInFileRef      bool
	Path                   string
	SourceTree             string
	BuildConfigurationList string
	MainGroup              string
	ProductRefGroup        string
	Children               []pbxObject
	ProductReference       string
}

func genCodeIosProjectPbxproj(mock *Mock, buf *CodeBuffer) {
	cp := mock.Meta.Ios.ClassPrefix
	pj := mock.Meta.Ios.Project
	fileId := 0
	pbxBuildFiles := map[string]pbxObject{}
	pbxFileReferences := map[string]pbxObject{}
	pbxFrameworksBuildPhases := map[string]pbxObject{}
	pbxGroups := map[string]pbxObject{}
	pbxNativeTargets := map[string]pbxObject{}
	pbxProjects := map[string]pbxObject{}
	pbxResourcesBuildPhases := map[string]pbxObject{}
	pbxSourcesBuildPhases := map[string]pbxObject{}
	pbxVariantGroups := map[string]pbxObject{}
	xcConfigurationLists := map[string]pbxObject{}
	// PBXBuildFile
	pbxBuildFiles["Foundation.framework"] = pbxObject{
		Name:     "Foundation.framework",
		Id:       genIosFileId(&fileId),
		Location: "Frameworks",
		FileRef:  genIosFileId(&fileId),
	}
	pbxBuildFiles["CoreGraphics.framework"] = pbxObject{
		Name:     "CoreGraphics.framework",
		Id:       genIosFileId(&fileId),
		Location: "Frameworks",
		FileRef:  genIosFileId(&fileId),
	}
	pbxBuildFiles["UIKit.framework"] = pbxObject{
		Name:     "UIKit.framework",
		Id:       genIosFileId(&fileId),
		Location: "Frameworks",
		FileRef:  genIosFileId(&fileId),
	}
	pbxBuildFiles["InfoPlist.strings"] = pbxObject{
		Name:     "InfoPlist.strings",
		Id:       genIosFileId(&fileId),
		Location: "Resources",
		FileRef:  genIosFileId(&fileId),
	}
	pbxBuildFiles["main.m"] = pbxObject{
		Name:     "main.m",
		Id:       genIosFileId(&fileId),
		Location: "Sources",
		FileRef:  genIosFileId(&fileId),
	}
	pbxBuildFiles[cp+"AppDelegate.m"] = pbxObject{
		Name:     cp + "AppDelegate.m",
		Id:       genIosFileId(&fileId),
		Location: "Sources",
		FileRef:  genIosFileId(&fileId),
	}
	pbxBuildFiles["Images.xcassets"] = pbxObject{
		Name:     "Images.xcassets",
		Id:       genIosFileId(&fileId),
		Location: "Resources",
		FileRef:  genIosFileId(&fileId),
	}
	// PBXFileReference
	pbxFileReferences[pj+".app"] = pbxObject{
		Name:             pj + ".app",
		Id:               genIosFileId(&fileId),
		ExplicitFileType: "wrapper.framework",
		IncludeInIndex:   "0",
		Path:             pj + ".app",
		SourceTree:       "BUILT_PRODUCTS_DIR",
	}
	pbxFileReferences["Foundation.framework"] = pbxObject{
		Name:              "Foundation.framework",
		Id:                pbxBuildFiles["Foundation.framework"].FileRef,
		LastKnownFileType: "wrapper.framework",
		ShowNameInFileRef: true,
		Path:              "System/Library/Frameworks/Foundation.framework",
		SourceTree:        "SDKROOT",
	}
	pbxFileReferences["CoreGraphics.framework"] = pbxObject{
		Name:              "CoreGraphics.framework",
		Id:                pbxBuildFiles["CoreGraphics.framework"].FileRef,
		LastKnownFileType: "wrapper.framework",
		ShowNameInFileRef: true,
		Path:              "System/Library/Frameworks/CoreGraphics.framework",
		SourceTree:        "SDKROOT",
	}
	pbxFileReferences["UIKit.framework"] = pbxObject{
		Name:              "UIKit.framework",
		Id:                pbxBuildFiles["UIKit.framework"].FileRef,
		LastKnownFileType: "wrapper.framework",
		ShowNameInFileRef: true,
		Path:              "System/Library/Frameworks/UIKit.framework",
		SourceTree:        "SDKROOT",
	}
	pbxFileReferences[cp+"AppDelegate.h"] = pbxObject{
		Name:              cp + "AppDelegate.h",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "sourcecode.c.h",
		Path:              cp + "AppDelegate.h",
		SourceTree:        "<group>",
	}
	pbxFileReferences[cp+"AppDelegate.m"] = pbxObject{
		Name:              cp + "AppDelegate.m",
		Id:                pbxBuildFiles[cp+"AppDelegate.m"].Id,
		LastKnownFileType: "sourcecode.c.objc",
		Path:              cp + "AppDelegate.m",
		SourceTree:        "<group>",
	}
	pbxFileReferences["Images.xcassets"] = pbxObject{
		Name:              "Images.xcassets",
		Id:                pbxBuildFiles["Images.xcassets"].FileRef,
		LastKnownFileType: "folder.assetcatalog",
		Path:              "Images.xcassets",
		SourceTree:        "<group>",
	}
	pbxFileReferences[pj+"-Info.plist"] = pbxObject{
		Name:              pj + "-Info.plist",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "text.plist.xml",
		Path:              pj + "-Info.plist",
		SourceTree:        "<group>",
	}
	pbxFileReferences["en"] = pbxObject{
		Name:              "en",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "text.plist.strings",
		Path:              "en.lproj/InfoPlist.strings",
		SourceTree:        "<group>",
	}
	pbxFileReferences["main.m"] = pbxObject{
		Name:              "main.m",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "sourcecode.c.h",
		Path:              "main.m",
		SourceTree:        "<group>",
	}
	pbxFileReferences[pj+"-Prefix.pch"] = pbxObject{
		Name:              pj + "-Prefix.pch",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "sourcecode.c.h",
		Path:              pj + "-Prefix.pch",
		SourceTree:        "<group>",
	}
	// PBXVariantGroup
	pbxVariantGroups["InfoPlist.strings"] = pbxObject{Name: "InfoPlist.strings", Id: genIosFileId(&fileId)}
	// PBXFrameworksBuildPhase
	pbxFrameworksBuildPhases["Frameworks"] = pbxObject{Name: "Frameworks", Id: genIosFileId(&fileId), Children: []pbxObject{
		pbxBuildFiles["Foundation.framework"],
		pbxBuildFiles["CoreGraphics.framework"],
		pbxBuildFiles["UIKit.framework"],
	}}
	// PBXGroup
	pbxGroups["Supporting Files"] = pbxObject{Name: "Supporting Files", Id: genIosFileId(&fileId), Children: []pbxObject{
		pbxFileReferences[pj+"-Info.plist"],
		pbxVariantGroups["InfoPlist.strings"],
		pbxFileReferences["main.m"],
		pbxFileReferences[pj+"-Prefix.pch"],
	}}
	pbxGroups[pj] = pbxObject{Name: pj, Id: genIosFileId(&fileId), Path: pj, Children: []pbxObject{
		pbxFileReferences[cp+"AppDelegate.h"],
		pbxFileReferences[cp+"AppDelegate.m"],
		pbxFileReferences["Images.xcassets"],
		pbxGroups["Supporting Files"],
	}}
	pbxGroups["Frameworks"] = pbxObject{Name: "Frameworks", Id: genIosFileId(&fileId), Children: []pbxObject{
		pbxFileReferences["Foundation.framework"],
		pbxFileReferences["CoreGraphics.framework"],
		pbxFileReferences["UIKit.framework"],
	}}
	pbxGroups["Products"] = pbxObject{Name: "Products", Id: genIosFileId(&fileId), Children: []pbxObject{
		pbxFileReferences[pj+".app"],
	}}
	pbxGroups["mainGroup"] = pbxObject{Id: genIosFileId(&fileId), Children: []pbxObject{
		pbxGroups[pj],
		pbxGroups["Frameworks"],
		pbxGroups["Products"],
	}}
	// PBXSourcesBuildPhase
	pbxSourcesBuildPhases["Sources"] = pbxObject{Name: "Sources", Id: genIosFileId(&fileId)}
	// PBXResourcesBuildPhase
	pbxResourcesBuildPhases["Resources"] = pbxObject{Name: "Resources", Id: genIosFileId(&fileId)}
	// XCConfigurationList
	xcConfigurationLists["PBXProject \""+pj+"\""] = pbxObject{Name: "PBXProject \"" + pj + "\"", Id: genIosFileId(&fileId)}
	xcConfigurationLists["PBXNativeTarget \""+pj+"\""] = pbxObject{Name: "PBXNativeTarget \"" + pj + "\"", Id: genIosFileId(&fileId)}
	// PBXNativeTarget
	pbxNativeTargets[pj] = pbxObject{
		Name: pj,
		Id:   genIosFileId(&fileId),
		BuildConfigurationList: "PBXNativeTarget \"" + pj + "\"",
		Children: []pbxObject{
			pbxSourcesBuildPhases["Sources"],
			pbxFrameworksBuildPhases["Frameworks"],
			pbxResourcesBuildPhases["Resources"],
		},
		ProductReference: pj + ".app",
	}
	// PBXProject
	pbxProjects["Project object"] = pbxObject{
		Name: "Project object",
		Id:   genIosFileId(&fileId),
		BuildConfigurationList: "PBXProjectTarget \"" + pj + "\"",
		MainGroup:              "mainGroup",
		ProductRefGroup:        "Products",
		Children: []pbxObject{
			pbxNativeTargets[pj],
		},
	}

	// Header
	buf.add(`// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 46;
	objects = {`)

	// PBXBuildFile section
	buf.add(`
/* Begin PBXBuildFile section */`)
	for _, key := range []string{
		"Foundation.framework",
		"CoreGraphics.framework",
		"UIKit.framework",
		"InfoPlist.strings",
		"main.m",
		cp + "AppDelegate.m",
		"Images.xcassets",
	} {
		buf.add(`		%s /* %s in %s */ = {isa = PBXBuildFile; fileRef = %s /* %s */; };`,
			pbxBuildFiles[key].Id,
			pbxBuildFiles[key].Name,
			pbxBuildFiles[key].Location,
			pbxBuildFiles[key].FileRef,
			pbxBuildFiles[key].Name)
	}
	buf.add(`/* End PBXBuildFile section */`)

	// PBXFileReference section
	buf.add(`
/* Begin PBXFileReference section */`)
	for _, fileRef := range pbxFileReferences {
		s := fmt.Sprintf(`		%s /* %s */ = {isa = PBXFileReference;`,
			fileRef.Id,
			fileRef.Name,
		)
		if fileRef.ExplicitFileType != "" {
			s += fmt.Sprintf(` explicitFileType = %s;`, fileRef.ExplicitFileType)
		} else if fileRef.LastKnownFileType != "" {
			s += fmt.Sprintf(` lastKnownFileType = %s;`, fileRef.LastKnownFileType)
		}
		if fileRef.IncludeInIndex != "" {
			s += fmt.Sprintf(` includeInIndex = %s;`, fileRef.IncludeInIndex)
		}
		if fileRef.ShowNameInFileRef {
			s += fmt.Sprintf(` name = %s;`, fileRef.Name)
		}
		path := fileRef.Path
		if strings.Contains(path, "-") {
			path = "\"" + path + "\""
		}
		s += fmt.Sprintf(` path = %s;`, fileRef.Path)
		sourceTree := fileRef.SourceTree
		if strings.Contains(sourceTree, "<") {
			sourceTree = "\"" + sourceTree + "\""
		}
		s += fmt.Sprintf(` sourceTree = %s; };`, sourceTree)
		buf.add(s)
	}
	buf.add(`/* End PBXFileReference section */`)

	// PBXFrameworksBuildPhase section
	buf.add(`
/* Begin PBXFrameworksBuildPhase section */
		%s /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (`, pbxFrameworksBuildPhases["Frameworks"].Name)
	for _, child := range pbxFrameworksBuildPhases["Frameworks"].Children {
		buf.add(`				%s /* %s in %s */,`,
			child.Id,
			child.Name,
			child.Location)
	}
	buf.add(`			);
			runOnlyForDeploymentPostprocessing = 0;
		};`)
	buf.add(`/* End PBXFrameworksBuildPhase section */`)

	// PBXGroup section
	buf.add(`
/* Begin PBXGroup section */`)
	for _, group := range pbxGroups {
		groupComment := ""
		if group.Name != "" {
			groupComment = "/* " + group.Name + " */ "
		}
		buf.add(`		%s %s= {
			isa = PBXGroup;
			children = (`,
			group.Id,
			groupComment)
		for _, child := range group.Children {
			buf.add(`				%s /* %s */,`,
				child.Id,
				child.Name)
		}
		buf.add(`			);`)

		if group.Path != "" {
			buf.add(`			path = %s;`, group.Path)
		} else if group.Name != "" {
			name := group.Name
			if strings.Contains(name, " ") {
				name = "\"" + name + "\""
			}
			buf.add(`			name = %s;`, name)
		}
		buf.add(`			sourceTree = "<group>";
		};`)
	}
	buf.add(`/* End PBXGroup section */`)

	// PBXNativetarget section
	buf.add(`
/* Begin PBXNativeTarget section */`)
	for _, nativeTarget := range pbxNativeTargets {
		buf.add(`		%s /* %s */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = %s /* Build configuration list for %s */;
			buildPhases = (`,
			nativeTarget.Id,
			nativeTarget.Name,
			xcConfigurationLists[nativeTarget.BuildConfigurationList].Id,
			xcConfigurationLists[nativeTarget.BuildConfigurationList].Name,
		)
		for _, child := range nativeTarget.Children {
			buf.add(`				%s /* %s */,`,
				child.Id,
				child.Name,
			)
		}
		buf.add(`			);
			buildRules = (
			);
			dependencies = (
			);
			name = %s;
			productName = %s;
			productReference = %s /* %s */;
			productType = "com.apple.product-type.application";
		};`,
			nativeTarget.Name,
			nativeTarget.Name,
			pbxFileReferences[nativeTarget.ProductReference].Id,
			pbxFileReferences[nativeTarget.ProductReference].Name,
		)
	}
	buf.add(`/* End PBXNativeTarget section */`)

	// FIXME PBXProject section
	buf.add(`
/* Begin PBXProject section */`)
	for _, project := range pbxProjects {
		buf.add(`		%s /* %s */ = {
			isa = PBXProject;
			attributes = {`,
			project.Id,
			project.Name,
		)
		buf.add(`				CLASSPREFIX = %s;
				LastUpgradeCheck = 510;
				ORGANIZATIONNAME = %s;
			};
			buildConfigurationList = %s /* Build configuration list for %s */;
			compatibilityVersion = "Xcode 3.2";
			developmentRegion = English;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
			);
			mainGroup = %s;
			productRefGroup = %s /* %s */;
			projectDirPath = "";
			projectRoot = "";
			targets = (`,
			cp,
			mock.Meta.Ios.OrganizationName,
			xcConfigurationLists[project.BuildConfigurationList].Id,
			xcConfigurationLists[project.BuildConfigurationList].Name,
			pbxGroups[project.MainGroup].Id,
			pbxGroups[project.ProductRefGroup].Id,
			pbxGroups[project.ProductRefGroup].Name,
		)
		for _, child := range project.Children {
			buf.add(`				%s /* %s */,`,
				child.Id,
				child.Name,
			)
			buf.add(`			);`)
		}
		buf.add(`		};`)
	}
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
}`, pbxProjects["Project object"].Id)
}

func genIosFileId(i *int) string {
	*i++
	return fmt.Sprintf("%024X", *i)
}
