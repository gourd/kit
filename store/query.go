package store

type Query interface {

	// SetLimit is setter of limit
	SetLimit(uint64) Query

	// GetLimit is getter of limit
	GetLimit() uint64

	// SetOffset is setter of offset
	SetOffset(uint64) Query

	// GetOffset is getter of offset
	GetOffset() uint64

	// SetConds sets the conds interface withing
	SetConds(Conds) Query

	// GetConds gets the conds interface within
	GetConds() Conds

	// AddCond add a Cond to Conds
	AddCond(prop string, val interface{}) Query

	// SetSorts sets the Sorts interface withing
	SetSorts(Sorts) Query

	// GetSorts gets the Sorts interface within
	GetSorts() Sorts

	// AddSort add a Sort to Sorts
	Sort(sstr string) Query
}

// NewQuery constructs a *BasicQuery and return as Query
func NewQuery() Query {
	q := &BasicQuery{
		Conds: &BasicConds{},
		Sorts: &BasicSorts{},
	}
	return q
}

// BasicQuery implements Query interface
type BasicQuery struct {
	Conds  Conds
	Sorts  Sorts
	Limit  uint64
	Offset uint64
}

// SetLimit is setter of limit
func (q *BasicQuery) SetLimit(n uint64) Query {
	q.Limit = n
	return q
}

// GetLimit is getter of limit
func (q *BasicQuery) GetLimit() uint64 {
	return q.Limit
}

// SetOffset is setter of offset
func (q *BasicQuery) SetOffset(n uint64) Query {
	q.Offset = n
	return q
}

// GetOffset is getter of offset
func (q *BasicQuery) GetOffset() uint64 {
	return q.Offset
}

// SetConds set the conds interface within
func (q *BasicQuery) SetConds(cs Conds) Query {
	q.Conds = cs
	return q
}

// GetConds get the conds interface within
func (q *BasicQuery) GetConds() Conds {
	return q.Conds
}

// AddCond add a Cond to Conds
func (q *BasicQuery) AddCond(p string, v interface{}) Query {
	q.Conds.Add(p, v)
	return q
}

// SetSorts set the Sorts interface within
func (q *BasicQuery) SetSorts(cs Sorts) Query {
	q.Sorts = cs
	return q
}

// GetSorts get the Sorts interface within
func (q *BasicQuery) GetSorts() Sorts {
	return q.Sorts
}

// Sort add a Sort to Sorts
func (q *BasicQuery) Sort(sstr string) Query {
	q.Sorts.Add(sstr)
	return q
}
