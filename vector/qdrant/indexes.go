package qdrant

import (
	"fmt"
	"net/http"
)

type CreateFieldIndexRequest struct {
	FieldName string                 `json:"field_name"`
	IndexType string                 `json:"index_type"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

type CreateFieldIndexResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) CreateFieldIndex(collectionName string, req CreateFieldIndexRequest) (*CreateFieldIndexResponse, error) {
	var res CreateFieldIndexResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/indexes/field", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DeleteFieldIndexRequest struct {
	FieldName string `json:"field_name"`
}

type DeleteFieldIndexResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) DeleteFieldIndex(collectionName string, req DeleteFieldIndexRequest) (*DeleteFieldIndexResponse, error) {
	var res DeleteFieldIndexResponse
	err := c.sendRequest(http.MethodDelete, fmt.Sprintf("/collections/%s/indexes/field", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
