package gen

const (
	GravityCenter       = "center"
	GravityCenterV      = "center_v"
	SizeFill            = "fill"
	SizeWrap            = "wrap"
	OrientationVertical = "vertical"
)

// Default layout params for widgets
type Widget struct {
	Name        string
	Textable    bool
	Gravity     string
	Orientation string
	SizeW       string
	SizeH       string
}

type WidgetsDef struct {
	widgets map[string]Widget
}

func (d *WidgetsDef) Add(name string, w Widget) {
	if d.widgets == nil {
		d.widgets = make(map[string]Widget)
	}
	d.widgets[name] = w
}

func (d *WidgetsDef) Has(name string) (ret bool) {
	ret = false
	if _, ok := d.widgets[name]; ok {
		ret = true
	}
	return
}

func (d *WidgetsDef) Get(name string) Widget {
	return d.widgets[name]
}
