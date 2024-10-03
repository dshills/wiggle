package qdrant

import (
	"fmt"
	"net/http"
)

type Point struct {
	ID      int               `json:"id"`
	Vector  []float64         `json:"vector"`
	Payload map[string]string `json:"payload,omitempty"`
}

type GetPointsRequest struct {
	IDs []int `json:"ids"`
}

type PointResponse struct {
	Result []Point `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) GetPoints(collectionName string, req GetPointsRequest) (*PointResponse, error) {
	var res PointResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type UpsertPointsRequest struct {
	Points []Point `json:"points"`
}

type UpsertPointsResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) UpsertPoints(collectionName string, req UpsertPointsRequest) (*UpsertPointsResponse, error) {
	var res UpsertPointsResponse
	err := c.sendRequest(http.MethodPut, fmt.Sprintf("/collections/%s/points", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetPoint(collectionName string, pointID int) (*PointResponse, error) {
	var res PointResponse
	err := c.sendRequest(http.MethodGet, fmt.Sprintf("/collections/%s/points/%d", collectionName, pointID), nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DeletePointsRequest struct {
	IDs []int `json:"ids"`
}

func (c *Client) DeletePoints(collectionName string, req DeletePointsRequest) error {
	var res CollectionResponse
	err := c.sendRequest(http.MethodDelete, fmt.Sprintf("/collections/%s/points/delete", collectionName), req, &res)
	if err != nil {
		return err
	}
	if res.Status != "ok" {
		return fmt.Errorf("failed to delete points: %s", res.Status)
	}
	return nil
}

type UpdateVectorsRequest struct {
	Points []Point `json:"points"`
}

type UpdateVectorsResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) UpdateVectors(collectionName string, req UpdateVectorsRequest) (*UpdateVectorsResponse, error) {
	var res UpdateVectorsResponse
	err := c.sendRequest(http.MethodPatch, fmt.Sprintf("/collections/%s/points/vectors", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DeleteVectorsRequest struct {
	IDs []int `json:"ids"`
}

func (c *Client) DeleteVectors(collectionName string, req DeleteVectorsRequest) error {
	var res CollectionResponse
	err := c.sendRequest(http.MethodDelete, fmt.Sprintf("/collections/%s/points/vectors/delete", collectionName), req, &res)
	if err != nil {
		return err
	}
	if res.Status != "ok" {
		return fmt.Errorf("failed to delete vectors: %s", res.Status)
	}
	return nil
}

type SetPayloadRequest struct {
	Payload map[string]interface{} `json:"payload"`
	IDs     []int                  `json:"points"`
}

type SetPayloadResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) SetPayload(collectionName string, req SetPayloadRequest) (*SetPayloadResponse, error) {
	var res SetPayloadResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/payload", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type OverwritePayloadRequest struct {
	Payload map[string]interface{} `json:"payload"`
	IDs     []int                  `json:"points"`
}

type OverwritePayloadResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) OverwritePayload(collectionName string, req OverwritePayloadRequest) (*OverwritePayloadResponse, error) {
	var res OverwritePayloadResponse
	err := c.sendRequest(http.MethodPut, fmt.Sprintf("/collections/%s/points/payload", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DeletePayloadRequest struct {
	Keys []string `json:"keys"`
	IDs  []int    `json:"points"`
}

type DeletePayloadResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) DeletePayload(collectionName string, req DeletePayloadRequest) (*DeletePayloadResponse, error) {
	var res DeletePayloadResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/payload/delete", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ClearPayloadRequest struct {
	IDs []int `json:"points"`
}

type ClearPayloadResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) ClearPayload(collectionName string, req ClearPayloadRequest) (*ClearPayloadResponse, error) {
	var res ClearPayloadResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/payload/clear", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type BatchUpdateRequest struct {
	UpsertPoints UpsertPointsRequest `json:"upsert_points,omitempty"`
	DeletePoints DeletePointsRequest `json:"delete_points,omitempty"`
	SetPayload   SetPayloadRequest   `json:"set_payload,omitempty"`
	ClearPayload ClearPayloadRequest `json:"clear_payload,omitempty"`
}

type BatchUpdateResponse struct {
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) BatchUpdatePoints(collectionName string, req BatchUpdateRequest) (*BatchUpdateResponse, error) {
	var res BatchUpdateResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/batch", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type CountPointsRequest struct {
	Filter map[string]interface{} `json:"filter,omitempty"`
	Exact  bool                   `json:"exact,omitempty"`
}

type CountPointsResponse struct {
	Count  int     `json:"count"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) CountPoints(collectionName string, req CountPointsRequest) (*CountPointsResponse, error) {
	var res CountPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/count", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ScrollPointsRequest struct {
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
	Payload map[string]string `json:"payload,omitempty"`
}

type ScrollPointsResponse struct {
	Result []Point `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

func (c *Client) ScrollPoints(collectionName string, req ScrollPointsRequest) (*ScrollPointsResponse, error) {
	var res ScrollPointsResponse
	err := c.sendRequest(http.MethodPost, fmt.Sprintf("/collections/%s/points/scroll", collectionName), req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
