package vector

// Vector defines the essential operations for interacting with a vector database in the context of a Wiggle node.
// Implementations of this interface should support basic vector database functionalities, including inserting,
// querying, updating, and deleting vectors, as well as managing vector indexes. Each method serves as a
// foundational operation required for vector-based operations like similarity search and vector updates.
//
// Methods:
//  - InsertVector: Adds a vector to the database, associated with a unique ID and optional metadata.
//  - QueryVector: Searches for vectors similar to a given query vector, returning the top K most similar vectors.
//  - UpdateVector: Modifies the vector or its metadata for a specific ID.
//  - DeleteVector: Removes a vector entry from the database using its ID.
//  - CreateIndex: Creates an index for efficient vector queries, based on the specified dimension.
//  - DeleteIndex: Removes an existing vector index from the database.
//  - CheckIndexExists: Checks whether an index with the specified name exists in the database.
//
// The interface is designed to be implementation-agnostic, allowing flexibility for various vector database
// backends such as Pinecone, Qdrant, or other vector search services.

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
