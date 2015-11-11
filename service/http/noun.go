package http

// Noun contains and return the singular and plural
// variant of a noun
type Noun interface {
	Singular() string
	Plural() string
}

// noun is default implementation of Noun
type noun struct {
	singular string
	plural   string
}

// Singular implements Noun
func (n noun) Singular() string {
	return n.singular
}

// Singular implements Plural
func (n noun) Plural() string {
	return n.plural
}

// NewNoun creates a Noun interface containing given info
func NewNoun(singular, plural string) Noun {
	return &noun{singular, plural}
}
