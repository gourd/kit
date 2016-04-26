package httpservice

import "path"

// Paths return the URL base and the RESTful
// service noun
type Paths interface {
	Base() string
	Noun() Noun
	Singular() string
	Plural() string
}

// NewPaths return the default Paths implementation with given
// information
func NewPaths(base string, noun Noun, idStr string) Paths {
	return &paths{
		base,
		noun,
		idStr,
	}
}

// paths is the default implementation of Paths
type paths struct {
	base  string
	noun  Noun
	idStr string
}

// Base implements Paths
func (p paths) Base() string {
	return p.base
}

// Noun implements Paths
func (p paths) Noun() Noun {
	return p.noun
}

// Singular implements Paths
func (p paths) Singular() string {
	return path.Join(p.base, p.Noun().Singular(), p.idStr)
}

// Plural implements Paths
func (p paths) Plural() string {
	return path.Join(p.base, p.Noun().Plural())
}
