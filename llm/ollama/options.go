package ollama

type Options struct {
	F16Kv              bool     `json:"f16_kv,omitempty"`
	FrequencyPenalty   float64  `json:"frequency_penalty,omitempty"`
	LowVram            bool     `json:"low_vram,omitempty"`
	MainGpu            int      `json:"main_gpu,omitempty"`
	Mirostat           int      `json:"mirostat,omitempty"`
	MirostatEta        float64  `json:"mirostat_eta,omitempty"`
	MirostatTau        float64  `json:"mirostat_tau,omitempty"`
	NumBatch           int      `json:"num_batch,omitempty"`
	NumCtx             int      `json:"num_ctx,omitempty"`
	NumGpu             int      `json:"num_gpu,omitempty"`
	NumGqa             int      `json:"num_gqa,omitempty"`
	NumKeep            int      `json:"num_keep,omitempty"`
	NumPredict         int      `json:"num_predict,omitempty"`
	NumThread          int      `json:"num_thread,omitempty"`
	Numa               bool     `json:"numa,omitempty"`
	PenalizeNewline    bool     `json:"penalize_newline,omitempty"`
	PresencePenalty    float64  `json:"presence_penalty,omitempty"`
	RepeatLastN        int      `json:"repeat_last_n,omitempty"`
	RepeatPenalty      float64  `json:"repeat_penalty,omitempty"`
	RopeFrequencyBase  float64  `json:"rope_frequency_base,omitempty"`
	RopeFrequencyScale float64  `json:"rope_frequency_scale,omitempty"`
	Seed               int      `json:"seed,omitempty"`
	Stop               []string `json:"stop,omitempty"`
	Temperature        float64  `json:"temperature,omitempty"`
	TfsZ               float64  `json:"tfs_z,omitempty"`
	TopK               int      `json:"top_k,omitempty"`
	TopP               float64  `json:"top_p,omitempty"`
	TypicalP           float64  `json:"typical_p,omitempty"`
	UseMlock           bool     `json:"use_mlock,omitempty"`
	UseMmap            bool     `json:"use_mmap,omitempty"`
	VocabOnly          bool     `json:"vocab_only,omitempty"`
}
