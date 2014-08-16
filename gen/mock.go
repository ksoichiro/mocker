package gen

type Options struct {
	InDir  string
	OutDir string
}

type Mock struct {
	Name    string
	Meta    Meta
	Screens []Screen
	Launch  Launch
	Colors  []Color
	Strings []String
}

type Meta struct {
	Android Android
	Ios     Ios
}

type Android struct {
	Package             string
	GradlePluginVersion string `json:"gradle_plugin_version"`
	BuildToolsVersion   string `json:"build_tools_version"`
	MinSdkVersion       int    `json:"min_sdk_version"`
	TargetSdkVersion    int    `json:"target_sdk_version"`
	CompileSdkVersion   string `json:"compile_sdk_version"`
	VersionCode         int    `json:"version_code"`
	VersionName         string `json:"version_name"`
}

type Ios struct {
	Project           string
	ClassPrefix       string `json:"class_prefix"`
	CompanyIdentifier string `json:"company_identifier"`
	DeploymentTarget  string `json:"deployment_target"`
}

type Screen struct {
	Id        string
	Name      string
	Layout    []View
	Behaviors []Behavior
}

type LayoutOptions struct {
	Width  string
	Height string
}

type View struct {
	Id      string
	Type    string
	Sub     []View
	Label   string
	Gravity string
	Below   string
	SizeW   string `json:"size_w"`
	SizeH   string `json:"size_h"`
	AlignH  string `json:"align_h"`
	AlignV  string `json:"align_v"`
	Margin  string
	Padding string
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
