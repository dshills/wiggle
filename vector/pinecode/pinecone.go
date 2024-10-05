package pinecone

import (
	"github.com/dshills/wiggle/vector" // Import the vector package
)

// Ensure VectorDB satisfies the vector.Vector interface.
var _ vector.Vector = (*VectorDB)(nil)

// VectorDB is an implementation of vector.Vector using Pinecone.
type VectorDB struct {
	client    *Client
	namespace string
}

// NewVectorDB creates a new instance of VectorDB.
func NewVectorDB(client *Client, namespace string) *VectorDB {
	return &VectorDB{
		client:    client,
		namespace: namespace,
	}
}

// InsertVector adds a vector with the given ID and metadata into the Pinecone index.
func (p *VectorDB) InsertVector(id string, vec []float64, _ map[string]interface{}) error {
	v := []Vector{{ID: id, Values: vec}}
	_, err := p.client.UpsertVectors(v, p.namespace)
	return err
}

// QueryVector searches for vectors similar to the provided vector and returns the top K results.
func (p *VectorDB) QueryVector(vec []float64, topK int) ([]vector.QueryResult, error) {
	results, err := p.client.QueryVectors(vec, topK, p.namespace)
	if err != nil {
		return nil, err
	}

	// Convert Pinecone query results to vector.QueryResult
	queryResults := make([]vector.QueryResult, len(results.Matches))
	for i, match := range results.Matches {
		queryResults[i] = vector.QueryResult{
			ID:    match.ID,
			Score: match.Score,
		}
	}
	return queryResults, nil
}

// UpdateVector updates an existing vector or inserts it if it doesn't exist.
func (p *VectorDB) UpdateVector(id string, vec []float64, metadata map[string]interface{}) error {
	// Pinecone upserts both new and existing vectors.
	return p.InsertVector(id, vec, metadata)
}

// DeleteVector removes the vector with the specified ID from the Pinecone index.
func (p *VectorDB) DeleteVector(id string) error {
	return p.client.DeleteVectors([]string{id}, p.namespace)
}

// CreateIndex initializes a new index in Pinecone with the specified dimension size.
func (p *VectorDB) CreateIndex(name string, dim int) error {
	return p.client.CreateIndex(CreateIndexRequest{Name: name, Dimension: dim, Namespace: p.namespace})
}

// DeleteIndex deletes the index from the Pinecone database.
func (p *VectorDB) DeleteIndex(name string) error {
	return p.client.DeleteIndex(name)
}

// CheckIndexExists verifies if an index already exists in Pinecone.
func (p *VectorDB) CheckIndexExists(name string) (bool, error) {
	index, err := p.client.DescribeIndex(name)
	if err != nil {
		return false, err
	}
	return index != nil, nil
}
