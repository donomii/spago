// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package adanorm

import (
	mat "github.com/nlpodyssey/spago/pkg/mat32"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel_Forward(t *testing.T) {
	model := New(0.8)
	g := ag.NewGraph()

	// == Forward
	x1 := g.NewVariable(mat.NewVecDense([]mat.Float{1.0, 2.0, 0.0, 4.0}), true)
	x2 := g.NewVariable(mat.NewVecDense([]mat.Float{3.0, 2.0, 1.0, 6.0}), true)
	x3 := g.NewVariable(mat.NewVecDense([]mat.Float{6.0, 2.0, 5.0, 1.0}), true)

	y := nn.Reify(model, g, nn.Training).(*Model).Forward(x1, x2, x3)

	assert.InDeltaSlice(t, []mat.Float{-0.4262454708, 0.1329389665, -1.0585727653, 1.0318792697}, y[0].Value().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []mat.Float{0, -0.4504751299, -0.9466645455, 1.0771396755}, y[1].Value().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []mat.Float{0.8524954413, -0.6244384413, 0.5397325589, -1.087789559}, y[2].Value().Data(), 1.0e-06)

	// == Backward
	y[0].PropagateGrad(mat.NewVecDense([]mat.Float{-1.0, -0.2, 0.4, 0.6}))
	y[1].PropagateGrad(mat.NewVecDense([]mat.Float{-0.3, 0.1, 0.7, 0.9}))
	y[2].PropagateGrad(mat.NewVecDense([]mat.Float{0.3, -0.4, 0.7, -0.8}))
	g.BackwardAll()

	assert.InDeltaSlice(t, []mat.Float{-0.4779089755, -0.0839735551, 0.4004185091, 0.1614640214}, x1.Grad().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []mat.Float{-0.2710945487, -0.0790678529, 0.2259110116, 0.12425139}, x2.Grad().Data(), 1.0e-06)
	assert.InDeltaSlice(t, []mat.Float{-0.1154695275, 0.0283184423, 0.1372573, -0.050106214}, x3.Grad().Data(), 1.0e-06)
}
