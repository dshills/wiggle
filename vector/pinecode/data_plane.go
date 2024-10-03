package pinecone

import (
	"fmt"
	"net/url"
	"strings"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

// Structs for various API operations
type Vector struct {
	ID     string    `json:"id"`
	Values []float64 `json:"values"`
}

type UpsertRequest struct {
	Vectors   []Vector `json:"vectors"`
	Namespace string   `json:"namespace,omitempty"`
}

type UpsertResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}

type QueryRequest struct {
	Namespace string    `json:"namespace,omitempty"`
	TopK      int       `json:"topK"`
	Include   []string  `json:"include,omitempty"`
	Vector    []float64 `json:"vector"`
}

type QueryResponse struct {
	Matches []Match `json:"matches"`
}

type Match struct {
	ID     string    `json:"id"`
	Score  float64   `json:"score"`
	Values []float64 `json:"values,omitempty"`
}

type FetchResponse struct {
	Vectors map[string]Vector `json:"vectors"`
}

type UpdateRequest struct {
	ID        string    `json:"id"`
	Values    []float64 `json:"values"`
	Namespace string    `json:"namespace,omitempty"`
}

type DeleteRequest struct {
	IDs       []string `json:"ids"`
	Namespace string   `json:"namespace,omitempty"`
}

type ListResponse struct {
	Indexes []string `json:"indexes"`
}

type IndexStatsResponse struct {
	Dimension  int                       `json:"dimension"`
	Namespaces map[string]NamespaceStats `json:"namespaces"`
}

type NamespaceStats struct {
	VectorCount int `json:"vectorCount"`
}

// UpsertVectors inserts or updates vectors in the index.
func (pc *Client) UpsertVectors(vectors []Vector, namespace string) (*UpsertResponse, error) {
	urlStr := fmt.Sprintf("%s/vectors/upsert", pc.dataPlaneURL)
	body := UpsertRequest{Vectors: vectors, Namespace: namespace}

	res, err := pc.sendRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}

	upsertResp, ok := res.(*UpsertResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return upsertResp, nil
}

// QueryVectors queries the index with a vector and retrieves the top K matches.
func (pc *Client) QueryVectors(vector []float64, topK int, namespace string) (*QueryResponse, error) {
	urlStr := fmt.Sprintf("%s/vectors/query", pc.dataPlaneURL)
	body := QueryRequest{Namespace: namespace, TopK: topK, Vector: vector, Include: []string{"values"}}

	res, err := pc.sendRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}

	queryResp, ok := res.(*QueryResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return queryResp, nil
}

// FetchVectors retrieves vectors by their IDs from the index.
func (pc *Client) FetchVectors(ids []string, namespace string) (*FetchResponse, error) {
	urlStr := fmt.Sprintf("%s/vectors/fetch", pc.dataPlaneURL)

	queryURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	queryParams := queryURL.Query()
	queryParams.Set("ids", strings.Join(ids, ","))
	if namespace != "" {
		queryParams.Set("namespace", namespace)
	}
	queryURL.RawQuery = queryParams.Encode()

	res, err := pc.sendRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	fetchResp, ok := res.(*FetchResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return fetchResp, nil
}

// UpdateVector updates a vector in the index.
func (pc *Client) UpdateVector(id string, values []float64, namespace string) error {
	urlStr := fmt.Sprintf("%s/vectors/update", pc.dataPlaneURL)
	body := UpdateRequest{ID: id, Values: values, Namespace: namespace}

	_, err := pc.sendRequest("POST", urlStr, body)
	return err
}

// DeleteVectors removes vectors from the index by their IDs.
func (pc *Client) DeleteVectors(ids []string, namespace string) error {
	urlStr := fmt.Sprintf("%s/vectors/delete", pc.dataPlaneURL)
	body := DeleteRequest{IDs: ids, Namespace: namespace}

	_, err := pc.sendRequest("POST", urlStr, body)
	return err
}

// ListIndexes retrieves the list of all indexes in the Pinecone project.
func (pc *Client) ListIndexes() (*ListResponse, error) {
	urlStr := fmt.Sprintf("%s/databases", pc.controlPlaneURL)
	res, err := pc.sendRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	listResp, ok := res.(*ListResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return listResp, nil
}

// DescribeIndexStats retrieves stats about the index, such as dimension and vector count.
func (pc *Client) DescribeIndexStats(namespace string) (*IndexStatsResponse, error) {
	urlStr := fmt.Sprintf("%s/stats", pc.dataPlaneURL)

	queryURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if namespace != "" {
		queryParams := queryURL.Query()
		queryParams.Set("namespace", namespace)
		queryURL.RawQuery = queryParams.Encode()
	}

	res, err := pc.sendRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	statsResp, ok := res.(*IndexStatsResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return statsResp, nil
}
