# Vector Database Interface for Wiggle Nodes

This package provides a generic interface for interacting with vector databases within the Wiggle framework. The Vector interface defines essential methods needed by a Wiggle node to insert, query, update, and manage vector data and indexes.

## Overview

The Vector interface abstracts the underlying operations of a vector database, allowing the Wiggle node to interact seamlessly with various vector database implementations such as Pinecone, Qdrant, or others. The focus is on supporting operations like similarity search, updating vectors, and managing indexes.

## Interface Methods

- InsertVector
    - Inserts a vector into the database, associating it with a unique ID and optional metadata.
	- Parameters:
		- id: A unique identifier for the vector.
		- vector: A slice of float64 values representing the vector.
		- metadata: An optional map of key-value pairs to attach to the vector.
- QueryVector
	- Searches for vectors similar to the provided query vector, returning the top K most similar vectors.
	- Parameters:
		- vector: A query vector as a slice of float64 values.
		- topK: The number of most similar vectors to return.
	- Returns:
		- A slice of QueryResult containing the vector IDs and their similarity scores.
- UpdateVector
	- Updates the vector or metadata of an existing entry by its unique ID.
	- Parameters:
		- id: The unique identifier of the vector to update.
		- vector: A new vector (optional).
		- metadata: Updated metadata (optional).
- DeleteVector
	- Removes the vector associated with the given ID from the database.
	- Parameters:
		- id: The unique identifier of the vector to delete.
- CreateIndex
	- Initializes an index in the vector database, optimized for vectors with a specified dimension size.
	- Parameters:
		- name: The name of the index.
		- dim: The dimensionality of vectors in the index.
- DeleteIndex
	- Deletes the index associated with the provided name.
	- Parameters:
		- name: The name of the index to delete.
- CheckIndexExists
	- Verifies whether an index exists in the vector database.
	- Parameters:
		- name: The name of the index to check.
	- Returns:
		- A boolean indicating if the index exists.

## Example Usage

To implement this interface for a specific vector database (e.g., Pinecone or Qdrant), you would create a struct that satisfies each of the method signatures.

```go
type PineconeDB struct {
    client *pinecone.Client
}

func (db *PineconeDB) InsertVector(id string, vector []float64, metadata map[string]interface{}) error {
    // Implementation for inserting a vector into Pinecone
    return nil
}

// Implement other methods as required...
```

## QueryResult

The QueryResult struct represents a search result from querying the vector database. It contains:

- ID: The identifier of the similar vector.
- Score: The similarity score between the query vector and the result vector.
