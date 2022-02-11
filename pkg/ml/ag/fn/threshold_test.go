// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThresholdForward(t *testing.T) {
	x := &variable[mat.Float]{
		value:        mat.NewVecDense([]mat.Float{0.1, -0.2, 3.3, 0.0}),
		grad:         nil,
		requiresGrad: true,
	}
	threshold := &variable[mat.Float]{
		value:        mat.NewScalar[mat.Float](2.0),
		grad:         nil,
		requiresGrad: false,
	}
	k := &variable[mat.Float]{
		value:        mat.NewScalar[mat.Float](1.6),
		grad:         nil,
		requiresGrad: false,
	}

	f := NewThreshold[mat.Float](x, threshold, k)
	y := f.Forward()

	assert.InDeltaSlice(t, []mat.Float{1.6, 1.6, 3.3, 1.6}, y.Data(), 1.0e-6)

	f.Backward(mat.NewVecDense([]mat.Float{-1.0, 0.5, 0.8, 0.0}))

	assert.InDeltaSlice(t, []mat.Float{0.0, 0.0, 0.8, 0}, x.grad.Data(), 1.0e-6)
}
