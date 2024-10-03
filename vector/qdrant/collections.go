package qdrant

import (
	"fmt"
	"net/http"
)

type CollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CollectionResponse struct {
	Status string                 `json:"status"`
	Result map[string]interface{} `json:"result"`
	Time   float64                `json:"time"`
}

func (c *Client) CreateCollection(req CollectionRequest) (*CollectionResponse, error) {
	var res CollectionResponse
	err := c.sendRequest(http.MethodPost, "/collections", req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetCollection(collectionName string) (*CollectionResponse, error) {
	var res CollectionResponse
	err := c.sendRequest(http.MethodGet, fmt.Sprintf("/collections/%s", collectionName), nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteCollection(collectionName string) error {
	var res CollectionResponse
	err := c.sendRequest(http.MethodDelete, fmt.Sprintf("/collections/%s", collectionName), nil, &res)
	if err != nil {
		return err
	}
	if res.Status != "ok" {
		return fmt.Errorf("failed to delete collection: %s", res.Status)
	}
	return nil
}

type UpdateCollectionRequest struct {
	OptimizersConfig map[string]interface{} `json:"optimizers_config,omitempty"`
	Params           map[string]interface{} `json:"params,omitempty"`
}

func (c *Client) UpdateCollection(collectionName string, req UpdateCollectionRequest) (*CollectionResponse, error) {
	var res CollectionResponse
	err := c.sendRequest(http.MethodPatch, fmt.Sprintf("/collections/%s", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (c *Client) GetCollections() (*CollectionResponse, error) {
	var res CollectionResponse
	err := c.sendRequest(http.MethodGet, "/collections", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CollectionExists(collectionName string) (bool, error) {
	url := fmt.Sprintf("/collections/%s", collectionName)
	req, err := http.NewRequest(http.MethodHead, c.baseURL+url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
