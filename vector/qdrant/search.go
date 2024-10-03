package qdrant

import (
	"fmt"
	"net/http"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

type SearchPointsRequest struct {
	Vector  []float64         `json:"vector"`
	Limit   int               `json:"limit"`
	Payload map[string]string `json:"payload,omitempty"`
}

type SearchPointsResponse struct {
	Result []Point `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) SearchPoints(collectionName string, req SearchPointsRequest) (*SearchPointsResponse, error) {
	var res SearchPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/search", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type BatchSearchPointsRequest struct {
	Searches []SearchPointsRequest `json:"searches"`
}

type BatchSearchPointsResponse struct {
	Result [][]Point `json:"result"`
	Status string    `json:"status"`
	Time   float64   `json:"time"`
}

func (c *Client) BatchSearchPoints(collectionName string, req BatchSearchPointsRequest) (*BatchSearchPointsResponse, error) {
	var res BatchSearchPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/search/batch", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type PointGroupsSearchRequest struct {
	Filter  map[string]interface{} `json:"filter,omitempty"`
	GroupBy string                 `json:"group_by"`
	Limit   int                    `json:"limit"`
}

type PointGroupsSearchResponse struct {
	Result []map[string]interface{} `json:"result"`
	Status string                   `json:"status"`
	Time   float64                  `json:"time"`
}

func (c *Client) SearchPointGroups(collectionName string, req PointGroupsSearchRequest) (*PointGroupsSearchResponse, error) {
	var res PointGroupsSearchResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/search/groups", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type RecommendPointsRequest struct {
	Positive []int             `json:"positive"`
	Negative []int             `json:"negative,omitempty"`
	Limit    int               `json:"limit"`
	Payload  map[string]string `json:"payload,omitempty"`
}

type RecommendPointsResponse struct {
	Result []Point `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) RecommendPoints(collectionName string, req RecommendPointsRequest) (*RecommendPointsResponse, error) {
	var res RecommendPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/recommend", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type BatchRecommendPointsRequest struct {
	Searches []RecommendPointsRequest `json:"searches"`
}

type BatchRecommendPointsResponse struct {
	Result [][]Point `json:"result"`
	Status string    `json:"status"`
	Time   float64   `json:"time"`
}

func (c *Client) BatchRecommendPoints(collectionName string, req BatchRecommendPointsRequest) (*BatchRecommendPointsResponse, error) {
	var res BatchRecommendPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/recommend/batch", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type RecommendPointGroupsRequest struct {
	Positive []int                  `json:"positive"`
	Negative []int                  `json:"negative,omitempty"`
	GroupBy  string                 `json:"group_by"`
	Limit    int                    `json:"limit"`
	Filter   map[string]interface{} `json:"filter,omitempty"`
}

type RecommendPointGroupsResponse struct {
	Result []map[string]interface{} `json:"result"`
	Status string                   `json:"status"`
	Time   float64                  `json:"time"`
}

func (c *Client) RecommendPointGroups(collectionName string, req RecommendPointGroupsRequest) (*RecommendPointGroupsResponse, error) {
	var res RecommendPointGroupsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/recommend/groups", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DiscoverPointsRequest struct {
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset,omitempty"`
	Filter map[string]interface{} `json:"filter,omitempty"`
}

type DiscoverPointsResponse struct {
	Result []Point `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) DiscoverPoints(collectionName string, req DiscoverPointsRequest) (*DiscoverPointsResponse, error) {
	var res DiscoverPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/discover", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DiscoverBatchPointsRequest struct {
	Queries []DiscoverPointsRequest `json:"queries"`
}

type DiscoverBatchPointsResponse struct {
	Result [][]Point `json:"result"`
	Status string    `json:"status"`
	Time   float64   `json:"time"`
}

func (c *Client) DiscoverBatchPoints(collectionName string, req DiscoverBatchPointsRequest) (*DiscoverBatchPointsResponse, error) {
	var res DiscoverBatchPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/discover/batch", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type QueryPointsRequest struct {
	Queries []map[string]interface{} `json:"queries"`
}

type QueryPointsResponse struct {
	Result []Point `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) QueryPoints(collectionName string, req QueryPointsRequest) (*QueryPointsResponse, error) {
	var res QueryPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/query", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type QueryBatchPointsRequest struct {
	Queries []map[string]interface{} `json:"queries"`
}

type QueryBatchPointsResponse struct {
	Result [][]Point `json:"result"`
	Status string    `json:"status"`
	Time   float64   `json:"time"`
}

func (c *Client) QueryBatchPoints(collectionName string, req QueryBatchPointsRequest) (*QueryBatchPointsResponse, error) {
	var res QueryBatchPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/query/batch", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type QueryPointGroupsRequest struct {
	Queries []map[string]interface{} `json:"queries"`
	GroupBy string                   `json:"group_by"`
}

type QueryPointGroupsResponse struct {
	Result []map[string]interface{} `json:"result"`
	Status string                   `json:"status"`
	Time   float64                  `json:"time"`
}

func (c *Client) QueryPointGroups(collectionName string, req QueryPointGroupsRequest) (*QueryPointGroupsResponse, error) {
	var res QueryPointGroupsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/query/groups", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
