// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequencelabeler

import (
	"net/http"

	"github.com/nlpodyssey/spago/pkg/nlp/sequencelabeler/grpcapi"
	"github.com/nlpodyssey/spago/pkg/utils/grpcutils"
	"github.com/nlpodyssey/spago/pkg/utils/httputils"
	"github.com/nlpodyssey/spago/pkg/webui/ner"
)

// Server is the spaGO built-in implementation of HTTP and gRPC server for
// sequence labeling.
type Server struct {
	model           *Model
	TimeoutSeconds  int
	MaxRequestBytes int

	// UnimplementedSequenceLabelerServer must be embedded to have forward compatible implementations for gRPC.
	grpcapi.UnimplementedSequenceLabelerServer
}

// NewServer returns a new Server.
func NewServer(model *Model) *Server {
	return &Server{
		model: model,
	}
}

// Start starts the HTTP and gRPC servers.
func (s *Server) Start(address, grpcAddress string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ner-ui", ner.Handler)
	mux.HandleFunc("/analyze", s.analyze)

	go httputils.RunHTTPServer(httputils.HTTPServerConfig{
		Address:         address,
		TimeoutSeconds:  s.TimeoutSeconds,
		MaxRequestBytes: s.MaxRequestBytes,
	}, mux)

	grpcServer := grpcutils.NewGRPCServer(grpcutils.GRPCServerConfig{
		TimeoutSeconds:  s.TimeoutSeconds,
		MaxRequestBytes: s.MaxRequestBytes,
	})
	grpcapi.RegisterSequenceLabelerServer(grpcServer, s)
	grpcutils.RunGRPCServer(grpcAddress, grpcServer)
}
