package pinecone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Ex Data Plane URL: https://<your-index-name>.svc.<your-region>.pinecone.io
// Ex Control Plane URL: https://controller.<your-region>.pinecone.io
const (
	ControlPlaneURLTemplate = "https://controller.<%REGION%>.pinecone.io"
	DataPlaneURLTemplate    = "https://<%INDEXNAME%>.svc.<%REGION%>.pinecone.io"
)

// Client is the primary struct for interacting with both the Data Plane and Control Plane APIs of Pinecone.
// It encapsulates all the necessary configuration for making requests to both planes, providing functionality to
// manage vector data (Data Plane) and manage indexes and collections (Control Plane).
//
// Fields:
//
//   - ApiKey: The API key used for authenticating all requests to Pinecone. This key is required for every
//     request made to either the Data Plane or Control Plane. Ensure this key is kept secure.
//
//   - DataPlane: The base URL for operations on vector data. This URL is tied to a specific index and region.
//     Operations like upserting, querying, fetching, updating, and deleting vectors, as well as generating embeddings,
//     are performed using this URL. It typically has the following format:
//     https://<index-name>.svc.<region>.pinecone.io
//     Each index will have its own Data Plane URL.
//
//   - ControlPlane: The base URL for managing indexes and collections. This URL is tied to the region of your
//     Pinecone project and is used for operations like creating, deleting, or configuring indexes, as well as managing
//     collections. The format for the Control Plane URL is typically:
//     https://controller.<region>.pinecone.io
//     The Control Plane URL is common across all indexes in a given region.
//
//   - Client: An instance of the standard Go http.Client, which is responsible for making the HTTP requests to
//     Pinecone’s API. It includes settings such as request timeouts, and can be configured further if necessary.
//     It’s generally reused across all requests to maintain performance and reduce overhead.
//
// Typical Usage:
// To use the Pinecone Client, you must initialize it with your API key and the appropriate Data Plane and Control Plane URLs.
// Once initialized, you can perform various vector operations on the Data Plane, such as upserting or querying vectors,
// or administrative tasks on the Control Plane, such as creating or deleting indexes.
//
// Example:
//
//	client := pinecone.NewClient("YOUR_API_KEY", "https://your-index.svc.your-region.pinecone.io", "https://controller.your-region.pinecone.io")
//
// This will allow you to use the client for all Pinecone operations, such as vector upsert, query, index creation,
// and collection management.
type Client struct {
	apiKey          string
	dataPlaneURL    string
	controlPlaneURL string // Base URL for Control Plane operations
	httpClient      *http.Client
}

// NewClient initializes a new Pinecone client.
func NewClient(apiKey, dataPlaneURL, controlPlaneURL string) *Client {
	return &Client{
		apiKey:          apiKey,
		controlPlaneURL: controlPlaneURL,
		dataPlaneURL:    dataPlaneURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Generic function for sending HTTP requests
func (pc *Client) sendRequest(method, urlStr string, body interface{}) (interface{}, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = nil
	}

	req, err := http.NewRequest(method, urlStr, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", pc.apiKey)

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
