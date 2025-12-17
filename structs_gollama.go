package gollama

// Models

type ModelInfo struct {
	Model string `json:"model"`
	Size  int    `json:"size"`
}

// ModelDetails

type ModelDetails struct {
	License    string `json:"license"`
	Modelfile  string `json:"modelfile"`
	Parameters string `json:"parameters"`
	Template   string `json:"template"`
	Details    struct {
		ParentModel       string   `json:"parent_model"`
		Format            string   `json:"format"`
		Family            string   `json:"family"`
		Families          []string `json:"families"`
		ParameterSize     string   `json:"parameter_size"`
		QuantizationLevel string   `json:"quantization_level"`
	} `json:"details"`
	ModelInfo struct {
		GeneralArchitecture               string   `json:"general.architecture"`
		GeneralBasename                   string   `json:"general.basename"`
		GeneralFileType                   int      `json:"general.file_type"`
		GeneralFinetune                   string   `json:"general.finetune"`
		GeneralLanguages                  []string `json:"general.languages"`
		GeneralLicense                    string   `json:"general.license"`
		GeneralParameterCount             int64    `json:"general.parameter_count"`
		GeneralQuantizationVersion        int      `json:"general.quantization_version"`
		GeneralSizeLabel                  string   `json:"general.size_label"`
		GeneralTags                       []string `json:"general.tags"`
		GeneralType                       string   `json:"general.type"`
		LlamaAttentionHeadCount           int      `json:"llama.attention.head_count"`
		LlamaAttentionHeadCountKv         int      `json:"llama.attention.head_count_kv"`
		LlamaAttentionLayerNormRmsEpsilon float64  `json:"llama.attention.layer_norm_rms_epsilon"`
		LlamaBlockCount                   int      `json:"llama.block_count"`
		LlamaContextLength                int      `json:"llama.context_length"`
		LlamaEmbeddingLength              int      `json:"llama.embedding_length"`
		LlamaFeedForwardLength            int      `json:"llama.feed_forward_length"`
		LlamaRopeDimensionCount           int      `json:"llama.rope.dimension_count"`
		LlamaRopeFreqBase                 int      `json:"llama.rope.freq_base"`
		LlamaVocabSize                    int      `json:"llama.vocab_size"`
		TokenizerGgmlBosTokenID           int      `json:"tokenizer.ggml.bos_token_id"`
		TokenizerGgmlEosTokenID           int      `json:"tokenizer.ggml.eos_token_id"`
		TokenizerGgmlMerges               any      `json:"tokenizer.ggml.merges"`
		TokenizerGgmlModel                string   `json:"tokenizer.ggml.model"`
		TokenizerGgmlPre                  string   `json:"tokenizer.ggml.pre"`
		TokenizerGgmlTokenType            any      `json:"tokenizer.ggml.token_type"`
		TokenizerGgmlTokens               any      `json:"tokenizer.ggml.tokens"`
	} `json:"model_info"`
	ModifiedAt string `json:"modified_at"`
}

// Format structs

type PromptImage struct {
	Filename string `json:"filename"`
}

type ItemProperty struct {
	Type       string                    `json:"type"`
	Properties map[string]FormatProperty `json:"properties,omitempty"`
	Enum       []string                  `json:"enum,omitempty"`
	Required   []string                  `json:"required,omitempty"`
	Items      *ItemProperty             `json:"items,omitempty"`
}

type FormatProperty struct {
	Type        string       `json:"type"`
	Description string       `json:"description,omitempty"`
	Enum        []string     `json:"enum,omitempty"`
	Items       ItemProperty `json:"items,omitempty"`
}

type StructuredFormat struct {
	Type       string                    `json:"type"`
	Properties map[string]FormatProperty `json:"properties"`
	Required   []string                  `json:"required,omitempty"`
}

type ToolFunction struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Parameters  StructuredFormat `json:"parameters"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolSource interface {
	ListTools() ([]Tool, error)
}

// Output structs

type ChatOuput struct {
	Role           string     `json:"role"`
	Content        string     `json:"content"`
	ToolCalls      []ToolCall `json:"tool_calls"`
	PromptTokens   int        `json:"prompt_tokens"`
	ResponseTokens int        `json:"response_tokens"`
}
