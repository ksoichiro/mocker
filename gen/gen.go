package gen

type Generator interface {
	Generate()
}

func NewGenerator(opt *Options, mock *Mock, genId string) Generator {
	var g Generator
	switch genId {
	case "ios":
		g = &IosGenerator{opt, mock}
	case "android":
		g = &AndroidGenerator{opt, mock}
	}
	return g
}
