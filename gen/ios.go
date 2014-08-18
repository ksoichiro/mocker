package gen

import (
	"encoding/hex"
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
		Name:     "button",
		Textable: true,
		Gravity:  GravityCenter,
		SizeW:    SizeFill,
		SizeH:    SizeWrap,
	})
	iwd.Add("label", Widget{
		Name:     "label",
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
			layoutCodeBuf := genIosViewControllerLayout(mock, dir, screen)
			genIosViewController(mock, dir, screen, &layoutCodeBuf)
		}(g.mock, projectDir, screen)
	}

	// Generate view helper
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosViewHelper(mock, dir)
	}(g.mock, projectDir)

	// Generate resources
	wg.Add(1)
	go func(mock *Mock, dir string) {
		defer wg.Done()
		genIosLocalizableStrings(mock, dir)
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
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, "Base.lproj", "InfoPlist.strings"))
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.Project, "ja.lproj", "InfoPlist.strings"))
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
	launchVcPrefix := mock.Meta.Ios.ClassPrefix + strings.Title(mock.Launch.Screen)
	buf.add(`#import "%sAppDelegate.h"
#import "%sViewController.h"

@implementation %sAppDelegate

- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions
{
    self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    self.window.backgroundColor = [UIColor whiteColor];
    self.window.rootViewController = [[UINavigationController alloc] initWithRootViewController:[%sViewController new]];

    [self.window makeKeyAndVisible];
    return YES;
}

@end`,
		mock.Meta.Ios.ClassPrefix,
		launchVcPrefix,
		mock.Meta.Ios.ClassPrefix,
		launchVcPrefix,
	)
}

func genIosViewController(mock *Mock, dir string, screen Screen, layoutCodeBuf *CodeBuffer) {
	var buf CodeBuffer
	genCodeIosViewControllerHeader(mock, screen, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.ClassPrefix+strings.Title(screen.Id)+"ViewController.h"))
	buf = CodeBuffer{}
	genCodeIosViewControllerImplementation(mock, screen, &buf, layoutCodeBuf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, mock.Meta.Ios.ClassPrefix+strings.Title(screen.Id)+"ViewController.m"))
}

func genCodeIosViewControllerHeader(mock *Mock, screen Screen, buf *CodeBuffer) {
	buf.add(`#import <UIKit/UIKit.h>

@interface %s%sViewController : UIViewController
`,
		mock.Meta.Ios.ClassPrefix,
		strings.Title(screen.Id))

	if 0 < len(screen.Layout) {
		views := []View{}
		genCodeIosAggregateWidgets(&screen.Layout[0], &views)
		for _, view := range views {
			widgetName := "UIView *"
			switch view.Type {
			case "button":
				widgetName = "UIButton *"
			case "label":
				widgetName = "UILabel *"
			}
			buf.add(`@property %s%s;`,
				widgetName,
				view.Id,
			)
		}
	}

	buf.add(`
@end`)
}

func genCodeIosAggregateWidgets(current *View, views *[]View) {
	if current != nil && iwd.Has((*current).Type) {
		if (*current).Id != "" {
			*views = append(*views, *current)
		}
		if 0 < len((*current).Sub) {
			for _, sub := range (*current).Sub {
				genCodeIosAggregateWidgets(&sub, views)
			}
		}
	}
}

func genCodeIosViewControllerImplementation(mock *Mock, screen Screen, buf *CodeBuffer, layoutCodeBuf *CodeBuffer) {
	buf.add(`#import "%s%sViewController.h"
#import "UIView+Extension.h"

@interface %s%sViewController ()

@end

@implementation %s%sViewController

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self) {
        UIView *root = [[UIView alloc] initWithFrame:CGRectMake(0, 44/*FIXME*/, self.view.frame.size.width, self.view.frame.size.height)];
        [self.view addSubview:root];
        NSMutableDictionary *views = [NSMutableDictionary new];
        [root createWithViewInfo:[self viewInfo] views:views];`,
		mock.Meta.Ios.ClassPrefix,
		strings.Title(screen.Id),
		mock.Meta.Ios.ClassPrefix,
		strings.Title(screen.Id),
		mock.Meta.Ios.ClassPrefix,
		strings.Title(screen.Id))

	if 0 < len(screen.Layout) {
		views := []View{}
		genCodeIosAggregateWidgets(&screen.Layout[0], &views)
		for _, view := range views {
			switch view.Type {
			case "button":
				buf.add(`
        if ([views.allKeys containsObject:@"%s"]) {
            self.%s = (UIButton *) [views objectForKey:@"%s"];
            [self.%s addTarget:self action:@selector(didPush%s) forControlEvents:UIControlEventTouchUpInside];
        }`, view.Id, view.Id, view.Id, view.Id, strings.Title(view.Id))
			case "label":
				buf.add(`
        if ([views.allKeys containsObject:@"%s"]) {
            self.%s = (UILabel *) [views objectForKey:@"%s"];
        }`, view.Id, view.Id, view.Id)
			default:
				buf.add(`
        if ([views.allKeys containsObject:@"%s"]) {
            self.%s = (UIView *) [views objectForKey:@"%s"];
        }`, view.Id, view.Id, view.Id)
			}
		}
	}

	buf.add(`    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
}

#pragma mark - Generated layout methods

/**
 * Creates view layout information as a dictionary.
 * Layout is not determined by generator, it's up to Objective-C.
 * Generator just passes the structure of the views.
 */
- (NSDictionary *)viewInfo
{
    return`)

	// Insert layout codes
	buf.join(layoutCodeBuf)

	buf.add(`}

@end`)
}

func genIosViewControllerLayout(mock *Mock, dir string, screen Screen) (buf CodeBuffer) {
	genCodeIosViewControllerLayout(mock, screen, &buf)
	return
}

func genCodeIosViewControllerLayout(mock *Mock, screen Screen, buf *CodeBuffer) {
	if 0 < len(screen.Layout) {
		genIosLayoutRecur(&screen.Layout[0], true, buf, 2, ";")
	}
}

func genIosLayoutRecur(view *View, top bool, buf *CodeBuffer, indent int, trail string) {
	if !iwd.Has(view.Type) {
		return
	}
	widget := iwd.Get(view.Type)

	t := tab(indent)
	tt := tab(indent + 1)
	buf.add(t + `@{`)

	matchParentW := "@YES"
	matchParentH := "@YES"
	base := view.SizeW
	if base == "" {
		base = widget.SizeW
	}
	if base == SizeFill {
		matchParentW = "@YES"
	} else {
		matchParentW = "@NO"
	}
	base = view.SizeH
	if base == "" {
		base = widget.SizeH
	}
	if base == SizeFill {
		matchParentH = "@YES"
	} else {
		matchParentH = "@NO"
	}
	buf.add(tt+`@"MatchParentWidth": %s,`, matchParentW)
	buf.add(tt+`@"MatchParentHeight": %s,`, matchParentH)

	hasSub := 0 < len(view.Sub)

	buf.add(tt+`@"Widget": @"%s",`, widget.Name)
	if view.Id != "" {
		buf.add(tt+`@"Id": @"%s",`, view.Id)
	}
	if view.Below != "" {
		buf.add(tt+`@"Below": @"%s",`, view.Below)
	}
	if widget.Textable && view.Label != "" {
		buf.add(tt+`@"Text": @"%s",`, view.Label)
	}
	if widget.Orientation != "" {
		buf.add(tt+`@"Orientation": @"%s",`, widget.Orientation)
	}
	if view.Gravity != "" {
		buf.add(tt+`@"Gravity": @"%s",`, view.Gravity)
	} else if widget.Gravity != "" {
		buf.add(tt+`@"Gravity": @"%s",`, widget.Gravity)
	}
	if view.Margin != "" {
		if view.Margin == "normal" {
			buf.add(tt + `@"Margin": @16,`)
		} else {
			buf.add(tt+`@"Margin": @%s`, view.Margin)
		}
	}
	if view.Padding != "" {
		if view.Padding == "normal" {
			buf.add(tt + `@"Padding": @16,`)
		} else {
			buf.add(tt+`@"Padding": @%s`, view.Padding)
		}
	}
	if hasSub {
		buf.add(tt + `@"Subviews": @[`)
		// Print sub views recursively
		for i, sv := range view.Sub {
			subTrail := ""
			if i < len(view.Sub)-1 {
				subTrail = ","
			}
			genIosLayoutRecur(&sv, false, buf, indent+2, subTrail)
		}
		buf.add(tt + `]`)
	}
	buf.add(t + `}` + trail)
}

func genIosViewHelper(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosViewHelperHeader(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, "UIView+Extension.h"))
	buf = CodeBuffer{}
	genCodeIosViewHelperImplementation(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, "UIView+Extension.m"))
}

func genCodeIosViewHelperHeader(mock *Mock, buf *CodeBuffer) {
	buf.add(`#import <UIKit/UIKit.h>

@interface UIView (Extension)

- (void)createWithViewInfo:(NSDictionary *)viewInfo views:(NSMutableDictionary *)views;

@end
`)
}

func genCodeIosViewHelperImplementation(mock *Mock, buf *CodeBuffer) {
	buf.add(`#import "UIView+Extension.h"

@implementation UIView (Extension)

- (void)createWithViewInfo:(NSDictionary *)viewInfo views:(NSMutableDictionary *)views
{
    // Get the max Y from sibling views
    CGFloat maxY = 0;
    for (UIView *sibling in self.subviews) {
        CGFloat y = CGRectGetMaxY(sibling.frame);
        if (maxY < y) {
            maxY = y;
        }
    }

    CGFloat margin = 0;
    if ([viewInfo.allKeys containsObject:@"Margin"]) {
        margin = [[viewInfo objectForKey:@"Margin"] floatValue];
    }
    CGFloat padding = 0;
    // FIXME This is wrong. Padding should be got from parent's viewInfo
    if ([viewInfo.allKeys containsObject:@"Padding"]) {
        padding = [[viewInfo objectForKey:@"Padding"] floatValue];
    }

    NSString *widget = [viewInfo objectForKey:@"Widget"];
    if ([widget isEqualToString:@"button"]) {
        // UIButton
        UIButton *button = [[UIButton alloc] initWithFrame:CGRectMake(margin + padding, maxY + margin + padding, self.frame.size.width - (margin + padding) * 2, 100/*FIXME This is temporary */)];
        [button setTitle:NSLocalizedString([viewInfo objectForKey:@"Text"], nil) forState:UIControlStateNormal];
        [button setTitleColor:[UIColor blackColor] forState:UIControlStateNormal];
        // TODO adjust button size with button text
        if ([viewInfo.allKeys containsObject:@"Id"]) {
            [views setObject:button forKey:[viewInfo objectForKey:@"Id"]];
        }
        [self addSubview:button];
        return;
    }
    if ([widget isEqualToString:@"label"]) {
        // UILabel
        UILabel *label = [[UILabel alloc] initWithFrame:CGRectMake(margin + padding, maxY + margin + padding, self.frame.size.width - (margin + padding) * 2, 100/*FIXME This is temporary */)];
        label.text = NSLocalizedString([viewInfo objectForKey:@"Text"], nil);
        // TODO adjust label size with label text
        [self addSubview:label];
        if ([viewInfo.allKeys containsObject:@"Id"]) {
            [views setObject:label forKey:[viewInfo objectForKey:@"Id"]];
        }
        return;
    }

    // Determine size with layout_width/layout_height/padding/margin
    CGFloat width = 0;
    BOOL matchParentWidth = YES;
    if ([viewInfo.allKeys containsObject:@"MatchParentWidth"]) {
        matchParentWidth = [[viewInfo objectForKey:@"MatchParentWidth"] boolValue];
    }
    // FIXME Currently, wrap_content is ignored
    //if (matchParentWidth) {
    width = self.frame.size.width - margin * 2 - padding * 2;
    //} else {
    // Wrap content
    // TODO Check gravity. Temporarily, it's left
    //}

    CGFloat height = 0;
    BOOL matchParentHeight = YES;
    if ([viewInfo.allKeys containsObject:@"MatchParentHeight"]) {
        matchParentHeight = [[viewInfo objectForKey:@"MatchParentHeight"] boolValue];
    }
    // FIXME Currently, wrap_content is ignored
    //if (matchParentHeight) {
    height = self.frame.size.height - margin * 2 - padding * 2;
    //} else {
    // Wrap content
    // TODO Check gravity. Temporarily, it's left
    //}

    // if ([widget isEqualToString:@"linear"] || [widget isEqualToString:@"relative"]) {
    // LinearLayout and RelativeLayout
    // TODO Separate each layout with each appropriate algorithms
    UIView *view = [[UIView alloc] initWithFrame:CGRectMake(0, 0, width, height)];
    if ([viewInfo.allKeys containsObject:@"Id"]) {
        [views setObject:view forKey:[viewInfo objectForKey:@"Id"]];
    }

    // Process subviews
    if ([viewInfo.allKeys containsObject:@"Subviews"]) {
        // TODO Issue: cannot determine subview size with parent size because it's not also determined yet.
        for (NSDictionary *subviewInfo in [viewInfo objectForKey:@"Subviews"]) {
            [view createWithViewInfo:subviewInfo views:views];
        }
    }

    // Add myself to parent finally
    [self addSubview:view];
}

@end`)
}

func genIosLocalizableStrings(mock *Mock, dir string) {
	for _, s := range mock.Strings {
		lang := s.Lang
		if strings.ToLower(lang) == "base" {
			// base -> Base
			lang = strings.Title(lang)
		}
		var buf CodeBuffer
		genCodeIosLocalizableStrings(s, &buf)
		genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, lang+".lproj", "Localizable.strings"))
	}
}

func genCodeIosLocalizableStrings(s String, buf *CodeBuffer) {
	for _, def := range s.Defs {
		buf.add(`"%s" = "%s";`, def.Id, def.Value)
	}
}

func genIosColors(mock *Mock, dir string) {
	var buf CodeBuffer
	genCodeIosColorHeader(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, "UIColor+Extension.h"))
	buf = CodeBuffer{}
	genCodeIosColorImplementation(mock, &buf)
	genFile(&buf, filepath.Join(dir, mock.Meta.Ios.Project, "UIColor+Extension.m"))
}

func genCodeIosColorHeader(mock *Mock, buf *CodeBuffer) {
	buf.add(`#import <UIKit/UIKit.h>

@interface UIColor (Extension)
`)

	for _, c := range mock.Colors {
		buf.add(`+ (UIColor *)%sColor;`, c.Id)
	}

	buf.add(`
@end`)
}

func genCodeIosColorImplementation(mock *Mock, buf *CodeBuffer) {
	buf.add(`#import "UIColor+Extension.h"

@implementation UIColor (Extension)
`)

	for _, c := range mock.Colors {
		a, r, g, b := hexToInt(c.Value)
		buf.add(`+ (UIColor *)%sColor { return [UIColor colorWithRed:%d/255.0 green:%d/255.0 blue:%d/255.0 alpha:%d/255.0]; }`, c.Id, r, g, b, a)
	}

	buf.add(`
@end`)
}

func hexToInt(hexString string) (a, r, g, b int) {
	raw := hexString
	// Remove prefix '#'
	if strings.HasPrefix(raw, "#") {
		braw := []byte(raw)
		raw = string(braw[1:])
	}

	// Format hex string
	if len(raw) == 8 {
		// AARRGGBB: Do nothing
	} else if len(raw) == 6 {
		// RRGGBB: Insert alpha(FF)
		raw = "FF" + raw
	} else if len(raw) == 4 {
		// ARGB: Duplicate each hex
		braw := []byte(raw)
		sa := string(braw[0:1])
		sr := string(braw[1:2])
		sg := string(braw[2:3])
		sb := string(braw[3:4])
		raw = sa + sa + sr + sr + sg + sg + sb + sb
	} else if len(raw) == 3 {
		// RGB: Insert alpha(F) and duplicate each hex
		raw = "F" + raw
		braw := []byte(raw)
		sa := string(braw[0:1])
		sr := string(braw[1:2])
		sg := string(braw[2:3])
		sb := string(braw[3:4])
		raw = sa + sa + sr + sr + sg + sg + sb + sb
	}
	bytes, _ := hex.DecodeString(raw)
	a = int(bytes[0])
	r = int(bytes[1])
	g = int(bytes[2])
	b = int(bytes[3])
	return
}

func convertIosLayoutOptions(widget Widget, view *View) (lo LayoutOptions) {
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

func tab(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += "    "
	}
	return s
}
