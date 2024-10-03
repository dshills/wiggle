package qdrant

import (
	"fmt"
	"net/http"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/
// https://api.qdrant.tech/api-reference/aliases/update-aliases

type UpdateAliasesRequest struct {
	Aliases map[string]string `json:"aliases"`
}

type UpdateAliasesResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) UpdateAliases(req UpdateAliasesRequest) (*UpdateAliasesResponse, error) {
	var res UpdateAliasesResponse
	err := c.sendRequest(http.MethodPost, "/aliases", req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetCollectionAliasesResponse struct {
	Aliases map[string]string `json:"aliases"`
	Status  string            `json:"status"`
	Time    float64           `json:"time"`
}

func (c *Client) GetCollectionAliases(collectionName string) (*GetCollectionAliasesResponse, error) {
	var res GetCollectionAliasesResponse
	err := c.sendRequest(http.MethodGet, fmt.Sprintf("/collections/%s/aliases", collectionName), nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type GetCollectionsAliasesResponse struct {
	Aliases map[string]map[string]string `json:"aliases"`
	Status  string                       `json:"status"`
	Time    float64                      `json:"time"`
}

func (c *Client) GetCollectionsAliases() (*GetCollectionsAliasesResponse, error) {
	var res GetCollectionsAliasesResponse
	err := c.sendRequest(http.MethodGet, "/aliases", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
