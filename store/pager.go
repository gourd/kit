package store

import (
	"encoding/json"
)

// Pager is an implementation of a pager descriptor
type Pager interface {

	// MarshalJSON implements json Marshaler interface
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON implements json Unarshaler interface
	UnmarshalJSON([]byte) error

	// SetTotal sets the total in the pager descriptor
	SetTotal(total int) Pager

	// GetTotal gets the total in the pager descriptor
	GetTotal() (total int)

	// SetLimit sets the limit and offset in the pager descriptor
	SetLimit(limit, offset int) Pager

	// GetLimit gets the limit and offset in the pager descriptor
	GetLimit() (limit, offset int)
}

// NewPager creates a new pager descriptor
func NewPager() Pager {
	return &pager{
		total:  -1,
		offset: -1,
		limit:  -1,
	}
}

// pager implements Pager
type pager struct {
	total  int
	offset int
	limit  int
}

// MarshalJSON implements json Marshaler interface
func (p pager) MarshalJSON() ([]byte, error) {
	vmap := make(map[string]int)
	if p.total > -1 {
		vmap["total"] = p.total
	}
	if p.limit > -1 {
		vmap["limit"] = p.limit
	}
	if p.offset > -1 {
		vmap["offset"] = p.offset
	}
	return json.Marshal(vmap)
}

// UnmarshalJSON implements json Unarshaler interface
func (p *pager) UnmarshalJSON(data []byte) (err error) {
	vmap := make(map[string]int)
	err = json.Unmarshal(data, &vmap)
	if err != nil {
		return
	}

	if total, ok := vmap["total"]; ok {
		p.total = total
	}
	if limit, ok := vmap["limit"]; ok {
		p.limit = limit
	}
	if offset, ok := vmap["offset"]; ok {
		p.offset = offset
	}

	return
}

// SetLimit implements the Pager interface
func (p *pager) SetLimit(limit, offset int) Pager {
	p.limit = limit
	p.offset = offset
	return p
}

// GetLimit implements the Pager interface
func (p pager) GetLimit() (limit, offset int) {
	return p.limit, p.offset
}

// SetTotal implements the Pager interface
func (p *pager) SetTotal(total int) Pager {
	p.total = total
	return p
}

// GetTotal implements the Pager interface
func (p pager) GetTotal() int {
	return p.total
}
