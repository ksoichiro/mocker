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

	// TODO Generate base file set

	var wg sync.WaitGroup

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

	wg.Wait()
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
