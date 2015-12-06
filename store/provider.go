package store

// Defs is the interface for collection
// of Source and Provider definitions
type Defs interface {

	// SetSource sets a Source to the key
	SetSource(srcKey interface{}, src Source)

	// GetSource gets a Source with the given key
	GetSource(srcKey interface{}) Source

	// Set associates a source and a store provider to the key
	Set(key, srcKey interface{}, provider Provider)

	// Get retrieve a source and a store provider
	// associated with the given key
	Get(key interface{}) (srcKey interface{}, provider Provider)
}

// NewDefs returns an empty Defs implementation
func NewDefs() Defs {
	return &defs{
		make(map[interface{}]Source),
		make(map[interface{}]storeDef),
	}
}

// storeDef contains definition of how to get
// a Store
type storeDef struct {
	srcKey   interface{}
	provider Provider
}

// defs contain all definitions of
// Source and Store
type defs struct {
	sources map[interface{}]Source
	stores  map[interface{}]storeDef
}

// SetSource implements Defs
func (d *defs) SetSource(key interface{}, src Source) {
	d.sources[key] = src
}

// GetSource implements Defs
func (d defs) GetSource(key interface{}) Source {
	if item, ok := d.sources[key]; ok {
		return item
	}
	return nil
}

// Set implements Defs
func (d *defs) Set(key, srcKey interface{}, provider Provider) {
	d.stores[key] = storeDef{srcKey, provider}
}

// Get implements Defs
func (d *defs) Get(key interface{}) (interface{}, Provider) {
	if def, ok := d.stores[key]; ok {
		return def.srcKey, def.provider
	}
	return nil, nil
}

// Conn is the interface to handle
// database connections session to Source
type Conn interface {

	// Raw returns the raw session object
	Raw() interface{}

	// Close disconnect the session to the data source
	Close()
}

// Source provides connection and, if any, connection error
type Source func() (Conn, error)

// Provider takes a connection and return Store
// for that session
type Provider func(sess interface{}) (Store, error)
