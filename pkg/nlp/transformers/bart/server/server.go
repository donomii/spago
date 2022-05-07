// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/nlpodyssey/spago/pkg/nlp/tokenizers/bpetokenizer"
	"github.com/nlpodyssey/spago/pkg/nlp/tokenizers/sentencepiece"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/head/conditionalgeneration"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/head/sequenceclassification"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/server/grpcapi"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/tasks"
	"github.com/nlpodyssey/spago/pkg/utils/grpcutils"
	"github.com/nlpodyssey/spago/pkg/utils/httputils"
	"github.com/nlpodyssey/spago/pkg/webui/bartnli"
	"net/http"
)

// Server contains everything needed to run a BART server.
type Server struct {
	model           nn.Model
	bpeTokenizer    *bpetokenizer.BPETokenizer
	spTokenizer     *sentencepiece.Tokenizer
	TimeoutSeconds  int
	MaxRequestBytes int

	// UnimplementedBARTServer must be embedded to have forward compatible implementations for gRPC.
	grpcapi.UnimplementedBARTServer
}

// NewServer returns a new Server.
func NewServer(
	model nn.Model,
	bpeTokenizer *bpetokenizer.BPETokenizer,
	spTokenizer *sentencepiece.Tokenizer,
) *Server {
	switch model.(type) {
	case *sequenceclassification.Model, *conditionalgeneration.Model: // ok
	default:
		panic("bart: invalid model type")
	}
	return &Server{
		model:        model,
		bpeTokenizer: bpeTokenizer,
		spTokenizer:  spTokenizer,
	}
}

// StartDefaultServer is used to start a basic BART gRPC server.
func (s *Server) StartDefaultServer(grpcAddress string) {
	grpcServer := grpcutils.NewGRPCServer(grpcutils.GRPCServerConfig{
		TimeoutSeconds:  s.TimeoutSeconds,
		MaxRequestBytes: s.MaxRequestBytes,
	})
	grpcapi.RegisterBARTServer(grpcServer, s)
	grpcutils.RunGRPCServer(grpcAddress, grpcServer)
}

// StartDefaultHTTPServer is used to start a basic BERT HTTP server.
// If you want more control of the HTTP server you can run your own
// HTTP router using the public handler functions
func (s *Server) StartDefaultHTTPServer(address string) {
	mux := http.NewServeMux()
	switch s.model.(type) {
	case *sequenceclassification.Model:
		mux.HandleFunc("/classify-nli-ui", bartnli.Handler)
		mux.HandleFunc("/classify", s.ClassifyHandler)
		mux.HandleFunc("/classify-nli", s.ClassifyNLIHandler)
	case *conditionalgeneration.Model:
		mux.HandleFunc("/generate", s.GenerateHandler)
	default:
		panic("bart: invalid model type")
	}

	go httputils.RunHTTPServer(httputils.HTTPServerConfig{
		Address:         address,
		TimeoutSeconds:  s.TimeoutSeconds,
		MaxRequestBytes: s.MaxRequestBytes,
	}, mux)
}

// Classify handles a classification request over gRPC.
func (s *Server) Classify(_ context.Context, req *grpcapi.ClassifyRequest) (*grpcapi.ClassifyReply, error) {
	result := s.classify(req.GetText(), req.GetText2())
	return classificationFrom(result), nil
}

// ClassifyNLI handles a zero-shot classification request over gRPC.
func (s *Server) ClassifyNLI(_ context.Context, req *grpcapi.ClassifyNLIRequest) (*grpcapi.ClassifyReply, error) {
	result, err := s.classifyNLI(
		req.GetText(),
		req.GetHypothesisTemplate(),
		req.GetPossibleLabels(),
		req.MultiClass,
	)
	if err != nil {
		return nil, err
	}
	return classificationFrom(result), nil
}

// Generate handles a conditional generation request over gRPC.
func (s *Server) Generate(_ context.Context, req *grpcapi.GenerateRequest) (*grpcapi.GenerateReply, error) {
	result, err := s.generate(req.GetText())
	if err != nil {
		return nil, err
	}
	return &grpcapi.GenerateReply{
		Text: result.Text,
		Took: result.Took,
	}, nil
}

type body struct {
	Text  string `json:"text"`
	Text2 string `json:"text2"`
	// Following fields used by ClassifyNLI
	HypothesisTemplate string   `json:"hypothesis_template"`
	PossibleLabels     []string `json:"possible_labels"`
	MultiClass         bool     `json:"multi_class"`
}

// ClassifyHandler handles a classify request over HTTP.
func (s *Server) ClassifyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // that's intended for testing purposes only
	w.Header().Set("Content-Type", "application/json")

	var content body
	err := json.NewDecoder(req.Body).Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := s.classify(content.Text, content.Text2)
	_, pretty := req.URL.Query()["pretty"]
	response, err := Dump(result, pretty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ClassifyNLIHandler handles a classify request over HTTP.
func (s *Server) ClassifyNLIHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // that's intended for testing purposes only
	w.Header().Set("Content-Type", "application/json")

	var content body
	err := json.NewDecoder(req.Body).Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.classifyNLI(
		content.Text,
		content.HypothesisTemplate,
		content.PossibleLabels,
		content.MultiClass,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, pretty := req.URL.Query()["pretty"]
	response, err := Dump(result, pretty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GenerateHandler handles a conditional generation request over HTTP.
func (s *Server) GenerateHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // that's intended for testing purposes only
	w.Header().Set("Content-Type", "application/json")

	var content body
	err := json.NewDecoder(req.Body).Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.generate(content.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, pretty := req.URL.Query()["pretty"]
	response, err := Dump(result, pretty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GenerateResponse is a JSON-serializable structure which holds server
// generation response data.
type GenerateResponse struct {
	Text string `json:"text"`
	// Took is the number of milliseconds it took the server to execute the request.
	Took int64 `json:"took"`
}

func classificationFrom(resp *tasks.ClassifyResponse) *grpcapi.ClassifyReply {
	distribution := make([]*grpcapi.ClassConfidencePair, len(resp.Distribution))
	for i, t := range resp.Distribution {
		distribution[i] = &grpcapi.ClassConfidencePair{
			Class:      t.Class,
			Confidence: float64(t.Confidence),
		}
	}
	return &grpcapi.ClassifyReply{
		Class:        resp.Class,
		Confidence:   float64(resp.Confidence),
		Distribution: distribution,
		Took:         resp.Took,
	}
}

const (
	defaultStartSequenceTokenID = 0
	defaultEndSequenceTokenID   = 2
)

func getInputIDs(tokenizer *bpetokenizer.BPETokenizer, text, text2 string) []int {
	encoded, _ := tokenizer.Encode(text) // TODO: error handling
	inputIds := append(append([]int{defaultStartSequenceTokenID}, encoded.IDs...), defaultEndSequenceTokenID)
	if text2 != "" {
		encoded2, _ := tokenizer.Encode(text2) // TODO: error handling
		inputIds2 := append(append([]int{defaultEndSequenceTokenID}, encoded2.IDs...), defaultEndSequenceTokenID)
		inputIds = append(inputIds, inputIds2...)
	}
	return inputIds
}

// Dump serializes the given value to JSON.
func Dump(value interface{}, pretty bool) ([]byte, error) {
	buf := bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if pretty {
		enc.SetIndent("", "    ")
	}
	enc.SetEscapeHTML(true)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
