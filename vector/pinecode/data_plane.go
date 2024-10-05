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
