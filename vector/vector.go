package vector

// VectorDB defines the operations that a Wiggle node needs to interact with a vector database.
type Vector interface {
	// InsertVector inserts a vector with the given ID and metadata into the database.
	InsertVector(id string, vector []float64, metadata map[string]interface{}) error

	// QueryVector searches for vectors in the database similar to the provided query vector.
	// It returns a slice of IDs of the most similar vectors and their respective scores.
	QueryVector(vector []float64, topK int) ([]QueryResult, error)

	// UpdateVector updates the vector or metadata of an existing entry by its ID.
	UpdateVector(id string, vector []float64, metadata map[string]interface{}) error

	// DeleteVector removes the vector entry with the specified ID from the database.
	DeleteVector(id string) error

	// CreateIndex initializes an index on the vector database with the given dimension size.
	CreateIndex(name string, dim int) error

	// DeleteIndex deletes the index in the vector database.
	DeleteIndex(name string) error

	// CheckIndexExists verifies if an index is present in the vector database.
	CheckIndexExists(name string) (bool, error)
}

// QueryResult represents a search result from querying the vector database.
type QueryResult struct {
	ID    string
	Score float64
}
