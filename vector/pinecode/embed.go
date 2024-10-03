package pinecone

import "fmt"

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

// Struct for Generate Embeddings Request
type GenerateEmbeddingsRequest struct {
	Texts []string `json:"texts"`
}

// Struct for Generate Embeddings Response
type GenerateEmbeddingsResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
}

// GenerateEmbeddings sends a request to generate embeddings for the provided texts.
func (pc *Client) GenerateEmbeddings(texts []string) (*GenerateEmbeddingsResponse, error) {
	urlStr := fmt.Sprintf("%s/embeddings/generate", pc.dataPlaneURL)

	requestBody := GenerateEmbeddingsRequest{
		Texts: texts,
	}

	res, err := pc.sendRequest("POST", urlStr, requestBody)
	if err != nil {
		return nil, err
	}

	embeddingsResp, ok := res.(*GenerateEmbeddingsResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return embeddingsResp, nil
}
