package pinecone

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
