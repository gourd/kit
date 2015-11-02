package upperio

import (
	"github.com/gourd/kit/store"
)

// Sort take a store query and returns upperio Sort usable parameter
func Sort(q store.Query) (res []interface{}) {
	ss := q.GetSorts().GetAll()
	res = make([]interface{}, 0, len(ss))
	for _, s := range ss {
		res = append(res, s.String())
	}
	return
}
