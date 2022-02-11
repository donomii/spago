// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRotateR_Forward(t *testing.T) {
	x := &variable[mat.Float]{
		value:        mat.NewVecDense([]mat.Float{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}),
		grad:         nil,
		requiresGrad: true,
	}
	f := NewRotateR[mat.Float](x, 1)
	y := f.Forward()

	assert.InDeltaSlice(t, []mat.Float{0.8, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7}, y.Data(), 1.0e-6)

	f.Backward(mat.NewVecDense([]mat.Float{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}))

	assert.InDeltaSlice(t, []mat.Float{
		0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.1,
	}, x.grad.Data(), 1.0e-6)
}
