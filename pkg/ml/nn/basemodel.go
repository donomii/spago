// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nn

import (
	"encoding/gob"
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
)

// BaseModel satisfies some methods of the Model interface.
// Don't use it directly; it is meant to be embedded in other processors to reduce the amount of boilerplate code.
type BaseModel[T mat.DType] struct {
	// G is the computational graph on which the (reified) model operates.
	G *ag.Graph[T]
	// ProcessingMode is the processing mode for the model (training or inference).
	ProcessingMode ProcessingMode
}

func init() {
	gob.Register(&BaseModel[float32]{})
	gob.Register(&BaseModel[float64]{})
}

// Mode returns whether the (reified) model is being used for training or inference.
func (m *BaseModel[_]) Mode() ProcessingMode {
	return m.ProcessingMode
}

// Graph returns the computational graph on which the (reified) model operates.
// It panics if the Graph is nil.
func (m *BaseModel[T]) Graph() *ag.Graph[T] {
	if m.G == nil {
		panic("nn: attempting to access Graph on a not reified model. Hint: use nn.Reify().")
	}
	return m.G
}

// IsProcessor returns whether the model has been reified (i.e., contextualized to operate
// on a graph) and can perform the Forward().
func (m *BaseModel[_]) IsProcessor() bool {
	return m.G != nil
}

// InitProcessor is used to initialize structures and data useful for the Forward().
// nn.Reify() automatically invokes InitProcessor() for any sub-models.
func (m *BaseModel[_]) InitProcessor() {}

// Close can be used to close or finalize model structures.
func (m *BaseModel[_]) Close() {}
