package pinecone

import (
	"fmt"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

type ListIndexesResponse struct {
	Indexes []string `json:"indexes"`
}

// ListIndexes retrieves the list of all indexes from the control plane.
func (pc *Client) ListIndexes() (*ListIndexesResponse, error) {
	urlStr := fmt.Sprintf("%s/databases", pc.controlPlaneURL)

	res, err := pc.sendRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	listResp, ok := res.(*ListIndexesResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return listResp, nil
}

// Struct for Create Index Request
type CreateIndexRequest struct {
	Name      string `json:"name"`
	Dimension int    `json:"dimension"`
	Metric    string `json:"metric"`
	Replicas  int    `json:"replicas,omitempty"`
	PodType   string `json:"pod_type,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func (pc *Client) CreateIndex(index CreateIndexRequest) error {
	urlStr := fmt.Sprintf("%s/databases", pc.controlPlaneURL)

	_, err := pc.sendRequest("POST", urlStr, index)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	fmt.Println("Index created successfully.")
	return nil
}

// Struct for Describe Index Response
type DescribeIndexResponse struct {
	Name      string `json:"name"`
	Metric    string `json:"metric"`
	Dimension int    `json:"dimension"`
	Replicas  int    `json:"replicas"`
	PodType   string `json:"pod_type"`
	Status    struct {
		State   string `json:"state"`
		Message string `json:"message,omitempty"`
	} `json:"status"`
}

// DescribeIndex retrieves details about a specific index from the control plane.
func (pc *Client) DescribeIndex(indexName string) (*DescribeIndexResponse, error) {
	urlStr := fmt.Sprintf("%s/databases/%s", pc.controlPlaneURL, indexName)

	res, err := pc.sendRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	describeResp, ok := res.(*DescribeIndexResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return describeResp, nil
}

func (pc *Client) DeleteIndex(indexName string) error {
	urlStr := fmt.Sprintf("%s/databases/%s", pc.controlPlaneURL, indexName)

	_, err := pc.sendRequest("DELETE", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to delete index: %w", err)
	}

	fmt.Println("Index deleted successfully.")
	return nil
}

// Struct for Configure Index Request
type ConfigureIndexRequest struct {
	Replicas int    `json:"replicas,omitempty"`
	PodType  string `json:"pod_type,omitempty"`
}

// ConfigureIndex modifies the configuration of an index in the Pinecone project.
func (pc *Client) ConfigureIndex(indexName string, config ConfigureIndexRequest) error {
	urlStr := fmt.Sprintf("%s/databases/%s", pc.controlPlaneURL, indexName)

	_, err := pc.sendRequest("PATCH", urlStr, config)
	if err != nil {
		return fmt.Errorf("failed to configure index: %w", err)
	}

	fmt.Println("Index configuration updated successfully.")
	return nil
}

// Struct for List Collections Response
type ListCollectionsResponse struct {
	Collections []string `json:"collections"`
}

// ListCollections retrieves the list of all collections in the Pinecone project.
func (pc *Client) ListCollections() (*ListCollectionsResponse, error) {
	urlStr := fmt.Sprintf("%s/collections", pc.controlPlaneURL)

	res, err := pc.sendRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	listResp, ok := res.(*ListCollectionsResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return listResp, nil
}

// Struct for Create Collection Request
type CreateCollectionRequest struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// CreateCollection creates a new collection in the Pinecone project.
func (pc *Client) CreateCollection(collection CreateCollectionRequest) error {
	urlStr := fmt.Sprintf("%s/collections", pc.controlPlaneURL)

	_, err := pc.sendRequest("POST", urlStr, collection)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	fmt.Println("Collection created successfully.")
	return nil
}

// Struct for Describe Collection Response
type DescribeCollectionResponse struct {
	Name   string `json:"name"`
	Size   int    `json:"size"`
	Status struct {
		State   string `json:"state"`
		Message string `json:"message,omitempty"`
	} `json:"status"`
	Source  string `json:"source"`
	Created string `json:"created"`
}

// DescribeCollection retrieves details about a specific collection from the control plane.
func (pc *Client) DescribeCollection(collectionName string) (*DescribeCollectionResponse, error) {
	urlStr := fmt.Sprintf("%s/collections/%s", pc.controlPlaneURL, collectionName)

	res, err := pc.sendRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	describeResp, ok := res.(*DescribeCollectionResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return describeResp, nil
}

// DeleteCollection deletes a specific collection from the Pinecone project.
func (pc *Client) DeleteCollection(collectionName string) error {
	urlStr := fmt.Sprintf("%s/collections/%s", pc.controlPlaneURL, collectionName)

	_, err := pc.sendRequest("DELETE", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	fmt.Println("Collection deleted successfully.")
	return nil
}
