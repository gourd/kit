package store

import (
	"fmt"
)

const (
	And = iota
	Or
)

// Conds is the general interface represents conditions
// all setters return itself so the operation can
// cascade
type Conds interface {

	// SetRel is setter of relation flag
	SetRel(int) Conds

	// GetRel is getter of relation flag
	GetRel() int

	// Add adds a condition
	Add(string, interface{}) Conds

	// GetAll gets the list of conditions
	GetAll() []Cond

	// GetMap gets the list of conditions in the
	// form of map[string]interface{}
	//
	// Compatible layer with upper.io
	// the result can be used as db.Cond
	GetMap() (map[string]interface{}, error)
}

// Cond is the generic condition statement
// aims to be interatable with upper.io and google datastore
type Cond struct {
	Prop  string
	Value interface{}
}

// NewConds creates Conds with a BasicConds
func NewConds() Conds {
	c := BasicConds{}
	c.Conds = make([]Cond, 0)
	return &c
}

// BasicConds implements of Conds
type BasicConds struct {
	Conds []Cond
	Rel   int
}

// SetRel is setter of relation flag
func (c *BasicConds) SetRel(r int) Conds {
	c.Rel = r
	return c
}

// GetRel is getter of relation flag
func (c *BasicConds) GetRel() int {
	return c.Rel
}

// Add adds a condition
func (c *BasicConds) Add(prop string, value interface{}) Conds {
	c.Conds = append(c.Conds, Cond{
		Prop:  prop,
		Value: value,
	})
	return c
}

// GetAll gets the list of conditions
func (c *BasicConds) GetAll() []Cond {
	return c.Conds
}

// GetMap gets the list of conditions in the
// form of map[string]interface{}
func (c *BasicConds) GetMap() (m map[string]interface{}, err error) {
	m = make(map[string]interface{})
	for _, cond := range c.Conds {
		if exists, ok := m[cond.Prop]; ok {
			err = fmt.Errorf("\"%s\" is mapped to both \"%#v\" and \"%#v\"",
				cond.Prop, exists, cond.Value)
		}
		m[cond.Prop] = cond.Value
	}
	return
}
