// Generated by gourd (version 0.1dev)
// Generated at 2015/11/02 17:15:59 (+0800)
// Note: If you want to re-generate this file in the future,
//       do not change it.

package oauth2

import (
	"github.com/gourd/kit/store"
	"github.com/gourd/kit/store/upperio"
	"net/http"

	"encoding/base64"
	"github.com/satori/go.uuid"
	"strings"

	"log"
	"upper.io/db"
)

func init() {
	// define store provider with proxy
	store.Providers.DefineFunc("AccessData", func(r *http.Request) (s store.Store, err error) {
		return GetAccessDataStore(r)
	})
}

// GetAccessDataStore provides raw AccessDataStore
func GetAccessDataStore(r *http.Request) (s *AccessDataStore, err error) {

	// obtain database
	db, err := upperio.Open(r, "default")
	if err != nil {
		return
	}

	// define store and return
	s = &AccessDataStore{db}
	return
}

// AccessDataStore serves generic CURD for type AccessData
// Generated by gourd CLI tool
type AccessDataStore struct {
	Db db.Database
}

// Create a AccessData in the database, of the parent
func (s *AccessDataStore) Create(
	cond store.Conds, ep store.EntityPtr) (err error) {

	// get collection
	coll, err := s.Coll()
	if err != nil {
		return
	}

	// apply random uuid string to string id

	uid := uuid.NewV4()
	e := ep.(*AccessData)
	e.Id = strings.TrimRight(base64.URLEncoding.EncodeToString(uid[:]), "=")

	// Marshal the item, if possible
	// (quick fix for upperio problem with db.Marshaler)
	if me, ok := ep.(db.Marshaler); ok {
		ep, err = me.MarshalDB()
		if err != nil {
			return
		}
	}

	// add the entity to collection

	_, err = coll.Append(ep)

	if err != nil {
		log.Printf("Error creating AccessData: %s", err.Error())
		err = store.ErrorInternal
		return
	}

	return
}

// Search a AccessData by its condition(s)
func (s *AccessDataStore) Search(
	q store.Query) store.Result {

	return upperio.NewResult(func() (res db.Result, err error) {
		// get collection
		coll, err := s.Coll()
		if err != nil {
			return
		}

		// retrieve entities by given query conditions
		conds := upperio.Conds(q.GetConds())
		if conds == nil {
			res = coll.Find()
		} else {
			res = coll.Find(conds)
		}

		// handle paging
		if q.GetOffset() != 0 {
			res = res.Skip(uint(q.GetOffset()))
		}
		if q.GetLimit() != 0 {
			res = res.Limit(uint(q.GetLimit()))
		}

		return
	})

}

// One returns the first AccessData matches condition(s)
func (s *AccessDataStore) One(
	c store.Conds, ep store.EntityPtr) (err error) {

	// retrieve results from database
	l := &[]AccessData{}
	q := store.NewQuery().SetConds(c)

	// dump results into pointer of map / struct
	err = s.Search(q).All(l)
	if err != nil {
		return
	}

	// if not found, report
	if len(*l) == 0 {
		err = store.ErrorNotFound
		return
	}

	// assign the value of given point
	// to the first retrieved value
	(*ep.(*AccessData)) = (*l)[0]
	return nil
}

// Update AccessData on condition(s)
func (s *AccessDataStore) Update(
	c store.Conds, ep store.EntityPtr) (err error) {

	// get collection
	coll, err := s.Coll()
	if err != nil {
		return
	}

	// get by condition and ignore the error
	cond, _ := c.GetMap()
	res := coll.Find(db.Cond(cond))

	// Marshal the item, if possible
	// (quick fix for upperio problem with db.Marshaler)
	if me, ok := ep.(db.Marshaler); ok {
		ep, err = me.MarshalDB()
		if err != nil {
			return
		}
	}

	// update the matched entities
	err = res.Update(ep)
	if err != nil {
		log.Printf("Error updating AccessData: %s", err.Error())
		err = store.ErrorInternal
	}
	return
}

// Delete AccessData on condition(s)
func (s *AccessDataStore) Delete(
	c store.Conds) (err error) {

	// get collection
	coll, err := s.Coll()
	if err != nil {
		return
	}

	// get by condition and ignore the error
	cond, _ := c.GetMap()
	res := coll.Find(db.Cond(cond))

	// remove the matched entities
	err = res.Remove()
	if err != nil {
		log.Printf("Error deleting AccessData: %s", err.Error())
		err = store.ErrorInternal
	}
	return nil
}

// AllocEntity allocate memory for an entity
func (s *AccessDataStore) AllocEntity() store.EntityPtr {
	return &AccessData{}
}

// AllocEntityList allocate memory for an entity list
func (s *AccessDataStore) AllocEntityList() store.EntityListPtr {
	return &[]AccessData{}
}

// Len inspect the length of an entity list
func (s *AccessDataStore) Len(pl store.EntityListPtr) int64 {
	el := pl.(*[]AccessData)
	return int64(len(*el))
}

// Coll return the raw upper.io collection
func (s *AccessDataStore) Coll() (coll db.Collection, err error) {
	// get raw collection
	coll, err = s.Db.Collection("oauth2_access")
	if err != nil {
		log.Printf("Error connecting collection oauth2_access: %s",
			err.Error())
		err = store.ErrorInternal
	}
	return
}

// Close the database session that AccessData is using
func (s *AccessDataStore) Close() error {
	return s.Db.Close()
}
