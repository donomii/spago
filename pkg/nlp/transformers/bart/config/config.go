// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"github.com/nlpodyssey/spago/pkg/mat"
	"os"
)

const (
	// DefaultConfigurationFile is the default BART JSON configuration filename.
	DefaultConfigurationFile = "config.json"
	// DefaultModelFile is the default BART spaGO model filename.
	DefaultModelFile = "spago_model.bin"
	// DefaultEmbeddingsStorage is the default directory name for BART model's embedding storage.
	DefaultEmbeddingsStorage = "embeddings_storage"
)

// Config contains the global configuration of the BART model and the heads of fine-tuning tasks.
// The configuration coincides with that of Hugging Face to facilitate compatibility between the two architectures.
type Config[T mat.DType] struct {
	NumLabels                  int               `json:"_num_labels"`
	ActivationDropout          T                 `json:"activation_dropout"`
	ActivationFunction         string            `json:"activation_function"`
	BiasLogits                 bool              `json:"add_bias_logits"`
	FinalLayerNorm             bool              `json:"add_final_layer_norm"`
	Architecture               []string          `json:"architectures"`
	AttentionDropout           T                 `json:"attention_dropout"`
	BosTokenID                 int               `json:"bos_token_id"`
	ClassifierDropout          T                 `json:"classif_dropout"`
	DModel                     int               `json:"d_model"`
	DecoderAttentionHeads      int               `json:"decoder_attention_heads"`
	DecoderFFNDim              int               `json:"decoder_ffn_dim"`
	DecoderLayerDrop           T                 `json:"decoder_layerdrop"`
	DecoderLayers              int               `json:"decoder_layers"`
	DecoderStartTokenID        int               `json:"decoder_start_token_id"`
	Dropout                    T                 `json:"dropout"`
	EncoderAttentionHeads      int               `json:"encoder_attention_heads"`
	EncoderFFNDim              int               `json:"encoder_ffn_dim"`
	EncoderLayerDrop           T                 `json:"encoder_layerdrop"`
	EncoderLayers              int               `json:"encoder_layers"`
	EosTokenID                 int               `json:"eos_token_id"`
	ExtraPosEmbedding          int               `json:"extra_pos_embeddings"`
	FineTuningTask             string            `json:"finetuning_task"`
	ForceBosTokenToBeGenerated bool              `json:"force_bos_token_to_be_generated"`
	ID2Label                   map[string]string `json:"id2label"`
	InitStd                    T                 `json:"init_std"`
	IsEncoderDecoder           bool              `json:"is_encoder_decoder"`
	Label2ID                   map[string]int    `json:"label2id"`
	MaxPositionEmbeddings      int               `json:"max_position_embeddings"`
	ModelType                  string            `json:"model_type"`
	NormalizeBefore            bool              `json:"normalize_before"`
	NormalizeEmbedding         bool              `json:"normalize_embedding"`
	NumHiddenLayers            int               `json:"num_hidden_layers"`
	OutputPast                 bool              `json:"output_past"`
	PadTokenID                 int               `json:"pad_token_id"`
	ScaleEmbedding             bool              `json:"scale_embedding"`
	StaticPositionEmbeddings   bool              `json:"static_position_embeddings"`
	TotalFlos                  T                 `json:"total_flos"`
	VocabSize                  int               `json:"vocab_size"`
	NumBeams                   int               `json:"num_beams"`
	MaxLength                  int               `json:"max_length"`
	BadWordsIDs                [][]int           `json:"bad_words_ids"`
	Training                   bool              `json:"training"` // Custom for spaGO
}

// Load loads a BART model Config from file.
func Load[T mat.DType](file string) (Config[T], error) {
	var config Config[T]
	configFile, err := os.Open(file)
	if err != nil {
		return Config[T]{}, err
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return Config[T]{}, err
	}
	return config, nil
}
