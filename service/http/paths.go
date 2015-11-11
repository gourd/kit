package http

import "path"

// Paths return the URL base and the RESTful
// service noun
type Paths interface {
	Base() string
	Noun() Noun
	Path(name string) string
	Method(name string) string
}

// NounToPath is the type to translate name and noun into
// basic path (path after base)
type NounToPath func(name string, noun Noun) (path, method string)

// NewPaths return the default Paths implementation with given
// information
func NewPaths(base string, noun Noun, toPath NounToPath) Paths {
	return &paths{
		base,
		noun,
		toPath,
	}
}

// paths is the default implementation of Paths
type paths struct {
	base   string
	noun   Noun
	toPath NounToPath
}

// Base implements Paths
func (p paths) Base() string {
	return p.base
}

// Noun implements Paths
func (p paths) Noun() Noun {
	return p.noun
}

// Path implements Paths
func (p paths) Path(name string) string {
	namedPath, _ := p.toPath(name, p.Noun())
	return path.Join(p.Base(), namedPath)
}

// Method implements Paths
func (p paths) Method(name string) string {
	_, method := p.toPath(name, p.Noun())
	return method
}
