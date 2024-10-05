package qdrant

import (
	"fmt"
	"strconv"

	"github.com/dshills/wiggle/vector"
)

type VectorDB struct {
	client     *Client
	collection string
}

// NewVectorDB creates a new instance of VectorDB.
func NewVectorDB(client *Client, collection string) *VectorDB {
	return &VectorDB{
		client:     client,
		collection: collection,
	}
}

// InsertVector adds a vector with the given ID and metadata into the Qdrant collection.
func (q *VectorDB) InsertVector(id string, vec []float64, metadata map[string]interface{}) error {
	pointID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %s", id)
	}

	point := Point{
		ID:      pointID,
		Vector:  vec,
		Payload: mapToStringMap(metadata),
	}

	req := UpsertPointsRequest{
		Points: []Point{point},
	}

	_, err = q.client.UpsertPoints(q.collection, req)
	return err
}

// QueryVector searches for vectors similar to the provided vector and returns the top K results.
func (q *VectorDB) QueryVector(vec []float64, topK int) ([]vector.QueryResult, error) {
	req := SearchPointsRequest{
		Vector: vec,
		Limit:  topK,
	}

	res, err := q.client.SearchPoints(q.collection, req)
	if err != nil {
		return nil, err
	}

	// Convert Qdrant search results to vector.QueryResult
	queryResults := make([]vector.QueryResult, len(res.Result))
	for i, point := range res.Result {
		queryResults[i] = vector.QueryResult{
			ID:    fmt.Sprintf("%d", point.ID),
			Score: 0, // Score not directly provided, but can be inferred later if needed
		}
	}

	return queryResults, nil
}

// UpdateVector updates an existing vector or inserts it if it doesn't exist.
func (q *VectorDB) UpdateVector(id string, vec []float64, metadata map[string]interface{}) error {
	return q.InsertVector(id, vec, metadata)
}

// DeleteVector removes the vector with the specified ID from the Qdrant collection.
func (q *VectorDB) DeleteVector(id string) error {
	pointID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %s", id)
	}

	req := DeletePointsRequest{
		IDs: []int{pointID},
	}

	return q.client.DeletePoints(q.collection, req)
}

// CreateIndex initializes a new collection in Qdrant with the specified dimension size and adds a field index.
func (q *VectorDB) CreateIndex(_ int, fieldName, indexType string, indexParams map[string]interface{}) error {
	collectionReq := CollectionRequest{
		Name: q.collection,
	}

	res, err := q.client.CreateCollection(collectionReq)
	if err != nil {
		return err
	}

	if res.Status != "ok" {
		return fmt.Errorf("failed to create collection: %s", res.Status)
	}

	fieldIndexReq := CreateFieldIndexRequest{
		FieldName: fieldName,
		IndexType: indexType,
		Params:    indexParams,
	}

	fieldIndexRes, err := q.client.CreateFieldIndex(q.collection, fieldIndexReq)
	if err != nil {
		return err
	}

	if fieldIndexRes.Status != "ok" {
		return fmt.Errorf("failed to create field index: %s", fieldIndexRes.Status)
	}

	return nil
}

// DeleteIndex deletes the collection and its field index from the Qdrant database.
func (q *VectorDB) DeleteIndex(fieldName string) error {
	fieldIndexReq := DeleteFieldIndexRequest{
		FieldName: fieldName,
	}

	fieldIndexRes, err := q.client.DeleteFieldIndex(q.collection, fieldIndexReq)
	if err != nil {
		return err
	}

	if fieldIndexRes.Status != "ok" {
		return fmt.Errorf("failed to delete field index: %s", fieldIndexRes.Status)
	}

	return q.client.DeleteCollection(q.collection)
}

// CheckIndexExists verifies if a collection exists in Qdrant.
func (q *VectorDB) CheckIndexExists() (bool, error) {
	exists, err := q.client.CollectionExists(q.collection)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Utility function to convert metadata map to map[string]string
func mapToStringMap(metadata map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range metadata {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
