package store

// Store defines interface of an entity service
type Store interface {
	// Basic entity operations
	Create(Conds, EntityPtr) error
	Search(Query) Result
	One(Conds, EntityPtr) error
	Update(Conds, EntityPtr) error
	Delete(Conds) error

	// Memory allocation
	AllocEntity() EntityPtr
	AllocEntityList() EntityListPtr

	// Helper
	Len(EntityListPtr) int64

	// Close the database session this
	// service is using
	Close() error
}
