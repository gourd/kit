package store

const (
	Asc = iota
	Desc
)

// Sorts describes an interface to a collection of Sort
type Sorts interface {
	Add(string) Sorts
	GetAll() []*Sort
}

// SortBy returns an Sort (ascending) of the given property name
func SortBy(n string) *Sort {
	return &Sort{
		Name: n,
	}
}

// BasicSorts implements Sorts
type BasicSorts []*Sort

// Add adds a Sort to the BasicSort collection
func (ss *BasicSorts) Add(sstr string) Sorts {
	*ss = append(*ss, SortStr(sstr))
	return ss
}

// GetAll returns the whole collection of Sort as slice
func (ss *BasicSorts) GetAll() []*Sort {
	return *ss
}

// SortStr return *Sort described by a given string
func SortStr(str string) *Sort {
	if str[0] == '-' {
		return &Sort{
			Name:  str[1:],
			Order: Desc,
		}
	}
	return &Sort{
		Name: str,
	}
}

// Sort is the generic description of a sorting
// aims to be interatable with upper.io and google datastore
type Sort struct {
	Name  string
	Order int
}

// Asc sets the Sort to order by asc order
func (s *Sort) Asc() *Sort {
	s.Order = Asc
	return s
}

// Desc sets the Sort to order by desc order
func (s *Sort) Desc() *Sort {
	s.Order = Desc
	return s
}

// String returns a string represetation to sorting
// which is compatible with upperio and Google Datastore
func (s *Sort) String() string {
	if s.Order == Desc {
		return "-" + s.Name
	}
	return s.Name
}
