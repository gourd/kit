package store

// Factory is the interface to manufacture Stores.
// It contains definitions of Source and Store
//
// Factory produces and manage Store instances
type Factory interface {

	// SetSource sets a Source to the key
	SetSource(srcKey interface{}, src Source)

	// GetSource gets a Source with the given key
	GetSource(srcKey interface{}) Source

	// Set associates a source and a store provider to a key (store key)
	Set(key, srcKey interface{}, provider Provider)

	// Get retrieve a source and a store provider
	// associated with the given key (store key)
	Get(key interface{}) (srcKey interface{}, provider Provider)
}

// NewFactory returns the default Factory implementation
func NewFactory() Factory {
	return &factoryDef{
		make(map[interface{}]Source),
		make(map[interface{}]storeDef),
	}
}

// storeDef contains definition of how to get a Store
type storeDef struct {
	srcKey   interface{}
	provider Provider
}

// factoryDef implements Factory
type factoryDef struct {
	sources map[interface{}]Source
	stores  map[interface{}]storeDef
}

// SetSource implements Factory.SetSource
func (d *factoryDef) SetSource(key interface{}, src Source) {
	d.sources[key] = src
}

// GetSource implements Factory.GetSource
func (d factoryDef) GetSource(key interface{}) Source {
	if item, ok := d.sources[key]; ok {
		return item
	}
	return nil
}

// Set implements Factory.Set
func (d *factoryDef) Set(key, srcKey interface{}, provider Provider) {
	d.stores[key] = storeDef{srcKey, provider}
}

// Get implements Factory.Get
func (d *factoryDef) Get(key interface{}) (interface{}, Provider) {
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
