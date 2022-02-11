// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gmlp

import (
	"encoding/gob"
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/nlpodyssey/spago/pkg/ml/nn/normalization/layernorm"
)

var _ nn.Model[float32] = &PreNorm[float32]{}

// PreNorm is a helper model to perform pre-normalization.
type PreNorm[T mat.DType] struct {
	nn.BaseModel[T]
	Block *Block[T]
	Norm  *layernorm.Model[T]
}

func init() {
	gob.Register(&PreNorm[float32]{})
	gob.Register(&PreNorm[float64]{})
}

// NewPreNorm returns a new PreNorm.
func NewPreNorm[T mat.DType](dim int, block *Block[T]) *PreNorm[T] {
	return &PreNorm[T]{
		Block: block,
		Norm:  layernorm.New[T](dim),
	}
}

// Forward performs the forward step.
func (m *PreNorm[T]) Forward(xs ...ag.Node[T]) []ag.Node[T] {
	ns := m.Norm.Forward(xs...)
	return m.Block.Forward(ns...)
}
