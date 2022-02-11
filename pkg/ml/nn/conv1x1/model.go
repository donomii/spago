// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package conv1x1 implements a 1-dimensional 1-kernel convolution model
package conv1x1

import (
	"encoding/gob"
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
)

// Model is a superficial depth-wise 1-dimensional convolution model.
// The following values are fixed: kernel size = 1; stride = 1; padding = 0,
type Model[T mat.DType] struct {
	nn.BaseModel[T]
	Config Config
	W      nn.Param[T] `spago:"type:weights"`
	B      nn.Param[T] `spago:"type:biases"`
}

var _ nn.Model[float32] = &Model[float32]{}

// Config provides configuration parameters for Model.
type Config struct {
	InputChannels  int
	OutputChannels int
}

func init() {
	gob.Register(&Model[float32]{})
	gob.Register(&Model[float64]{})
}

// New returns a new Model.
func New[T mat.DType](config Config) *Model[T] {
	return &Model[T]{
		Config: config,
		W:      nn.NewParam[T](mat.NewEmptyDense[T](config.OutputChannels, config.InputChannels)),
		B:      nn.NewParam[T](mat.NewEmptyVecDense[T](config.OutputChannels)),
	}
}

// Forward performs the forward step. Each "x" is a channel.
func (m *Model[T]) Forward(xs ...ag.Node[T]) []ag.Node[T] {
	g := m.Graph()

	xm := g.Stack(xs...)
	mm := g.Mul(m.W, xm)

	ys := make([]ag.Node[T], m.Config.OutputChannels)
	for outCh := range ys {
		val := g.T(g.RowView(mm, outCh))
		bias := g.AtVec(m.B, outCh)
		ys[outCh] = g.AddScalar(val, bias)
	}
	return ys
}
