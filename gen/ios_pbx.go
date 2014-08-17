package gen

import (
	"fmt"
	"path/filepath"
	"strings"
)

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
	FileEncoding           string
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
	BuildSettings          string
}

func genCodeIosProjectPbxproj(mock *Mock, buf *CodeBuffer) {
	cp := mock.Meta.Ios.ClassPrefix
	pj := mock.Meta.Ios.Project
	fileId := 0xDE5E0B8D
	pbxBuildFiles := map[string]pbxObject{}
	pbxFileReferences := map[string]pbxObject{}
	pbxFrameworksBuildPhases := map[string]pbxObject{}
	pbxGroups := map[string]pbxObject{}
	pbxNativeTargets := map[string]pbxObject{}
	pbxProjects := map[string]pbxObject{}
	pbxResourcesBuildPhases := map[string]pbxObject{}
	pbxSourcesBuildPhases := map[string]pbxObject{}
	pbxVariantGroups := map[string]pbxObject{}
	xcProjectBuildConfigurations := map[string]pbxObject{}
	xcNativeTargetBuildConfigurations := map[string]pbxObject{}
	xcConfigurationLists := map[string]pbxObject{}
	// PBXFileReference
	pbxFileReferences[pj+".app"] = pbxObject{
		Name:             pj + ".app",
		Id:               genIosFileId(&fileId),
		ExplicitFileType: "wrapper.application",
		IncludeInIndex:   "0",
		Path:             pj + ".app",
		SourceTree:       "BUILT_PRODUCTS_DIR",
	}
	pbxFileReferences["Foundation.framework"] = pbxObject{
		Name:              "Foundation.framework",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "wrapper.framework",
		ShowNameInFileRef: true,
		Path:              "System/Library/Frameworks/Foundation.framework",
		SourceTree:        "SDKROOT",
	}
	pbxFileReferences["CoreGraphics.framework"] = pbxObject{
		Name:              "CoreGraphics.framework",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "wrapper.framework",
		ShowNameInFileRef: true,
		Path:              "System/Library/Frameworks/CoreGraphics.framework",
		SourceTree:        "SDKROOT",
	}
	pbxFileReferences["UIKit.framework"] = pbxObject{
		Name:              "UIKit.framework",
		Id:                genIosFileId(&fileId),
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
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "sourcecode.c.objc",
		Path:              cp + "AppDelegate.m",
		SourceTree:        "<group>",
	}
	pbxFileReferences["Images.xcassets"] = pbxObject{
		Name:              "Images.xcassets",
		Id:                genIosFileId(&fileId),
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
		ShowNameInFileRef: true,
		Path:              "en.lproj/InfoPlist.strings",
		SourceTree:        "<group>",
	}
	pbxFileReferences["main.m"] = pbxObject{
		Name:              "main.m",
		Id:                genIosFileId(&fileId),
		LastKnownFileType: "sourcecode.c.objc",
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
	// ViewControllers for each Screens
	for _, screen := range mock.Screens {
		hname := cp + strings.Title(screen.Id) + "ViewController.h"
		pbxFileReferences[hname] = pbxObject{
			Name:              hname,
			Id:                genIosFileId(&fileId),
			FileEncoding:      "4",
			LastKnownFileType: "sourcecode.c.h",
			Path:              hname,
			SourceTree:        "<group>",
		}

		mname := cp + strings.Title(screen.Id) + "ViewController.m"
		pbxFileReferences[mname] = pbxObject{
			Name:              mname,
			Id:                genIosFileId(&fileId),
			FileEncoding:      "4",
			LastKnownFileType: "sourcecode.c.objc",
			Path:              mname,
			SourceTree:        "<group>",
		}
	}
	// PBXVariantGroup
	pbxVariantGroups["InfoPlist.strings"] = pbxObject{
		Name: "InfoPlist.strings",
		Id:   genIosFileId(&fileId),
		Children: []pbxObject{
			pbxFileReferences["en"],
		},
	}
	// PBXBuildFile
	pbxBuildFiles["Foundation.framework"] = pbxObject{
		Name:     "Foundation.framework",
		Id:       genIosFileId(&fileId),
		Location: "Frameworks",
		FileRef:  pbxFileReferences["Foundation.framework"].Id,
	}
	pbxBuildFiles["CoreGraphics.framework"] = pbxObject{
		Name:     "CoreGraphics.framework",
		Id:       genIosFileId(&fileId),
		Location: "Frameworks",
		FileRef:  pbxFileReferences["CoreGraphics.framework"].Id,
	}
	pbxBuildFiles["UIKit.framework"] = pbxObject{
		Name:     "UIKit.framework",
		Id:       genIosFileId(&fileId),
		Location: "Frameworks",
		FileRef:  pbxFileReferences["UIKit.framework"].Id,
	}
	pbxBuildFiles["InfoPlist.strings"] = pbxObject{
		Name:     "InfoPlist.strings",
		Id:       genIosFileId(&fileId),
		Location: "Resources",
		FileRef:  pbxVariantGroups["InfoPlist.strings"].Id,
	}
	pbxBuildFiles["main.m"] = pbxObject{
		Name:     "main.m",
		Id:       genIosFileId(&fileId),
		Location: "Sources",
		FileRef:  pbxFileReferences["main.m"].Id,
	}
	pbxBuildFiles[cp+"AppDelegate.m"] = pbxObject{
		Name:     cp + "AppDelegate.m",
		Id:       genIosFileId(&fileId),
		Location: "Sources",
		FileRef:  pbxFileReferences[cp+"AppDelegate.m"].Id,
	}
	pbxBuildFiles["Images.xcassets"] = pbxObject{
		Name:     "Images.xcassets",
		Id:       genIosFileId(&fileId),
		Location: "Resources",
		FileRef:  pbxFileReferences["Images.xcassets"].Id,
	}
	// ViewControllers for each Screens
	for _, screen := range mock.Screens {
		name := cp + strings.Title(screen.Id) + "ViewController.m"
		pbxBuildFiles[name] = pbxObject{
			Name:     name,
			Id:       genIosFileId(&fileId),
			Location: "Sources",
			FileRef:  pbxFileReferences[name].Id,
		}
	}
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
	vcFileRefs := []pbxObject{
		pbxFileReferences[cp+"AppDelegate.h"],
		pbxFileReferences[cp+"AppDelegate.m"],
		pbxFileReferences["Images.xcassets"],
	}
	for _, screen := range mock.Screens {
		vcFileRefs = append(vcFileRefs, pbxFileReferences[cp+strings.Title(screen.Id)+"ViewController.h"])
		vcFileRefs = append(vcFileRefs, pbxFileReferences[cp+strings.Title(screen.Id)+"ViewController.m"])
	}
	vcFileRefs = append(vcFileRefs, pbxGroups["Supporting Files"])
	pbxGroups[pj] = pbxObject{Name: pj, Id: genIosFileId(&fileId), Path: pj, Children: vcFileRefs}
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
	vcBuildFiles := []pbxObject{
		pbxBuildFiles["main.m"],
	}
	for _, screen := range mock.Screens {
		vcBuildFiles = append(vcBuildFiles, pbxBuildFiles[cp+strings.Title(screen.Id)+"ViewController.m"])
	}
	vcBuildFiles = append(vcBuildFiles, pbxBuildFiles[cp+"AppDelegate.m"])
	pbxSourcesBuildPhases["Sources"] = pbxObject{
		Name:     "Sources",
		Id:       genIosFileId(&fileId),
		Children: vcBuildFiles,
	}
	// PBXResourcesBuildPhase
	pbxResourcesBuildPhases["Resources"] = pbxObject{
		Name: "Resources",
		Id:   genIosFileId(&fileId),
		Children: []pbxObject{
			pbxBuildFiles["InfoPlist.strings"],
			pbxBuildFiles["Images.xcassets"],
		},
	}
	// XCConfiguration
	xcProjectBuildConfigurations["Debug"] = pbxObject{
		Name: "Debug",
		Id:   genIosFileId(&fileId),
		BuildSettings: fmt.Sprintf(`				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
				CLANG_CXX_LIBRARY = "libc++";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CLANG_WARN_BOOL_CONVERSION = YES;
				CLANG_WARN_CONSTANT_CONVERSION = YES;
				CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
				CLANG_WARN_EMPTY_BODY = YES;
				CLANG_WARN_ENUM_CONVERSION = YES;
				CLANG_WARN_INT_CONVERSION = YES;
				CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
				COPY_PHASE_STRIP = NO;
				GCC_C_LANGUAGE_STANDARD = gnu99;
				GCC_DYNAMIC_NO_PIC = NO;
				GCC_OPTIMIZATION_LEVEL = 0;
				GCC_PREPROCESSOR_DEFINITIONS = (
					"DEBUG=1",
					"$(inherited)",
				);
				GCC_SYMBOLS_PRIVATE_EXTERN = NO;
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				IPHONEOS_DEPLOYMENT_TARGET = %s;
				ONLY_ACTIVE_ARCH = YES;
				SDKROOT = iphoneos;`,
			mock.Meta.Ios.DeploymentTarget),
	}
	xcProjectBuildConfigurations["Release"] = pbxObject{
		Name: "Release",
		Id:   genIosFileId(&fileId),
		BuildSettings: fmt.Sprintf(`				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
				CLANG_CXX_LIBRARY = "libc++";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CLANG_WARN_BOOL_CONVERSION = YES;
				CLANG_WARN_CONSTANT_CONVERSION = YES;
				CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
				CLANG_WARN_EMPTY_BODY = YES;
				CLANG_WARN_ENUM_CONVERSION = YES;
				CLANG_WARN_INT_CONVERSION = YES;
				CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
				COPY_PHASE_STRIP = YES;
				ENABLE_NS_ASSERTIONS = NO;
				GCC_C_LANGUAGE_STANDARD = gnu99;
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				IPHONEOS_DEPLOYMENT_TARGET = %s;
				SDKROOT = iphoneos;
				VALIDATE_PRODUCT = YES;`,
			mock.Meta.Ios.DeploymentTarget),
	}
	xcNativeTargetBuildConfigurations["Debug"] = pbxObject{
		Name: "Debug",
		Id:   genIosFileId(&fileId),
		BuildSettings: `				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				ASSETCATALOG_COMPILER_LAUNCHIMAGE_NAME = LaunchImage;
				GCC_PRECOMPILE_PREFIX_HEADER = YES;
				GCC_PREFIX_HEADER = "MockerDemo/MockerDemo-Prefix.pch";
				INFOPLIST_FILE = "MockerDemo/MockerDemo-Info.plist";
				PRODUCT_NAME = "$(TARGET_NAME)";
				WRAPPER_EXTENSION = app;`,
	}
	xcNativeTargetBuildConfigurations["Release"] = pbxObject{
		Name: "Release",
		Id:   genIosFileId(&fileId),
		BuildSettings: `				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				ASSETCATALOG_COMPILER_LAUNCHIMAGE_NAME = LaunchImage;
				GCC_PRECOMPILE_PREFIX_HEADER = YES;
				GCC_PREFIX_HEADER = "MockerDemo/MockerDemo-Prefix.pch";
				INFOPLIST_FILE = "MockerDemo/MockerDemo-Info.plist";
				PRODUCT_NAME = "$(TARGET_NAME)";
				WRAPPER_EXTENSION = app;`,
	}
	// XCConfigurationList
	xcConfigurationLists["PBXProject \""+pj+"\""] = pbxObject{
		Name: "PBXProject \"" + pj + "\"",
		Id:   genIosFileId(&fileId),
		Children: []pbxObject{
			xcProjectBuildConfigurations["Debug"],
			xcProjectBuildConfigurations["Release"],
		},
	}
	xcConfigurationLists["PBXNativeTarget \""+pj+"\""] = pbxObject{
		Name: "PBXNativeTarget \"" + pj + "\"",
		Id:   genIosFileId(&fileId),
		Children: []pbxObject{
			xcNativeTargetBuildConfigurations["Debug"],
			xcNativeTargetBuildConfigurations["Release"],
		},
	}
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
		BuildConfigurationList: "PBXProject \"" + pj + "\"",
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
		s += fmt.Sprintf(` path = %s;`, path)
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
		%s /* %s */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (`,
		pbxFrameworksBuildPhases["Frameworks"].Id,
		pbxFrameworksBuildPhases["Frameworks"].Name,
	)
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

	// PBXProject section
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
				LastUpgradeCheck = 0510;
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

	// PBXResourcesBuildPhase
	buf.add(`
/* Begin PBXResourcesBuildPhase section */`)
	for _, resourcesBuildPhase := range pbxResourcesBuildPhases {
		buf.add(`		%s /* %s */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (`,
			resourcesBuildPhase.Id,
			resourcesBuildPhase.Name,
		)
		for _, child := range resourcesBuildPhase.Children {
			buf.add(`				%s /* %s in %s */,`,
				child.Id,
				child.Name,
				child.Location,
			)
		}
		buf.add(`			);
			runOnlyForDeploymentPostprocessing = 0;
		};`)
	}
	buf.add(`/* End PBXResourcesBuildPhase section */`)

	// PBXSourcesBuildPhase
	buf.add(`
/* Begin PBXSourcesBuildPhase section */`)
	for _, sourcesBuildPhase := range pbxSourcesBuildPhases {
		buf.add(`		%s /* %s */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (`,
			sourcesBuildPhase.Id,
			sourcesBuildPhase.Name,
		)
		for _, child := range sourcesBuildPhase.Children {
			buf.add(`				%s /* %s in %s */,`,
				child.Id,
				child.Name,
				child.Location,
			)
		}
		buf.add(`			);
			runOnlyForDeploymentPostprocessing = 0;
		};`)
	}
	buf.add(`/* End PBXSourcesBuildPhase section */`)

	// PBXVariantGroup
	buf.add(`
/* Begin PBXVariantGroup section */`)
	for _, variantGroup := range pbxVariantGroups {
		buf.add(`		%s /* %s */ = {
			isa = PBXVariantGroup;
			children = (`,
			variantGroup.Id,
			variantGroup.Name,
		)
		for _, child := range variantGroup.Children {
			buf.add(`				%s /* %s */,`,
				child.Id,
				child.Name,
			)
		}
		buf.add(`			);
			name = InfoPlist.strings;
			sourceTree = "<group>";
		};`)
	}
	buf.add(`/* End PBXVariantGroup section */`)

	// XCBuildConfiguration section
	buf.add(`
/* Begin XCBuildConfiguration section */`)
	for _, xcbc := range xcProjectBuildConfigurations {
		buf.add(`		%s /* %s */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
%s
			};
			name = %s;
		};`,
			xcbc.Id,
			xcbc.Name,
			xcbc.BuildSettings,
			xcbc.Name,
		)
	}
	for _, xcbc := range xcNativeTargetBuildConfigurations {
		buf.add(`		%s /* %s */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
%s
			};
			name = %s;
		};`,
			xcbc.Id,
			xcbc.Name,
			xcbc.BuildSettings,
			xcbc.Name,
		)
	}
	buf.add(`/* End XCBuildConfiguration section */`)

	buf.add(`
/* Begin XCConfigurationList section */`)
	for _, c := range xcConfigurationLists {
		buf.add(`		%s /* Build configuration list for %s */ = {
			isa = XCConfigurationList;
			buildConfigurations = (`,
			c.Id,
			c.Name,
		)
		for _, child := range c.Children {
			buf.add(`				%s /* %s */,`,
				child.Id,
				child.Name,
			)
		}
		buf.add(`			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};`)
	}
	buf.add(`/* End XCConfigurationList section */`)

	// Footer
	buf.add(`	};
	rootObject = %s /* Project object */;
}`, pbxProjects["Project object"].Id)
}

func genIosFileId(i *int) string {
	*i++
	return fmt.Sprintf("%08X199F4B790030B30B", *i)
}
