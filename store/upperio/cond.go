package upperio

import (
	"github.com/gourd/kit/store"

	"fmt"
	"upper.io/db"
)

// Conds Translate the store.Conds interface into
// upperio flavor conditions representation
func Conds(cs store.Conds) interface{} {
	conds := cs.GetAll()
	out := make([]interface{}, 0)

	// accumulate relations
	for _, cond := range conds {
		if cond.Prop == "" {
			// if no prop, assume to be prop
			if v, ok := cond.Value.(string); ok {
				out = append(out, db.Raw{v})
			} else if v, ok := cond.Value.(store.Conds); ok {
				leaf := Conds(v)
				out = append(out, leaf)
			} else if v, ok := cond.Value.(db.Raw); ok {
				out = append(out, v)
			} else if v, ok := cond.Value.(db.And); ok {
				out = append(out, v)
			} else if v, ok := cond.Value.(db.Or); ok {
				out = append(out, v)
			}
		} else {
			out = append(out, db.Cond{cond.Prop: cond.Value})
		}
	}

	if len(out) == 0 {
		return nil // nil for empty query, searchs everything
	}

	// determine relations
	if cs.GetRel() == store.And {
		return db.And(out)
	} else if cs.GetRel() == store.Or {
		return db.Or(out)
	}

	panic(fmt.Sprintf("Incorrect value of Rel in %#v", cs))
	return nil

}
