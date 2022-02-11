// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSoftmax_Forward(t *testing.T) {
	x := &variable[mat.Float]{
		value:        mat.NewVecDense([]mat.Float{-0.41, -1.08, 0, 0.87, -0.19, -0.75}),
		grad:         nil,
		requiresGrad: true,
	}
	f := NewSoftmax[mat.Float](x)
	y := f.Forward()

	assert.InDeltaSlice(t, []mat.Float{0.1166451, 0.0596882, 0.1757629, 0.4195304, 0.1453487, 0.083024}, y.Data(), 1.0e-6)

	f.Backward(mat.NewVecDense([]mat.Float{0.0, 0.0, -5.689482, 0.0, 0.0, 0.0}))

	assert.InDeltaSlice(t, []mat.Float{0.1166451, 0.0596882, -0.8242370, 0.4195304, 0.1453487, 0.083024}, x.grad.Data(), 1.0e-6)
}
