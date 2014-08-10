package gen

import "fmt"

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

type CodeBuffer []string

func (b *CodeBuffer) add(format string, a ...interface{}) {
	*b = append(*b, fmt.Sprintf(format, a...))
}
