// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/saientist/spago/pkg/mat"
	"gonum.org/v1/gonum/floats"
	"testing"
)

func TestView_Forward(t *testing.T) {

	x := &variable{
		value: mat.NewDense(3, 4, []float64{
			0.1, 0.2, 0.3, 0.0,
			0.4, 0.5, -0.6, 0.7,
			-0.5, 0.8, -0.8, -0.1,
		}),
		grad:         nil,
		requiresGrad: true,
	}

	f := NewView(x, 1, 1, 2, 2)
	y := f.Forward()

	if !floats.EqualApprox(y.Data(), []float64{
		0.5, -0.6,
		0.8, -0.8,
	}, 1.0e-6) {
		t.Error("The output doesn't match the expected values")
	}

	if y.Rows() != 2 || y.Columns() != 2 {
		t.Error("The rows and columns of the resulting matrix are not correct")
	}

	f.Backward(mat.NewDense(2, 2, []float64{
		0.1, 0.2,
		-0.8, -0.1,
	}))

	if !floats.EqualApprox(x.grad.Data(), []float64{
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.1, 0.2, 0.0,
		0.0, -0.8, -0.1, 0.0,
	}, 1.0e-6) {
		t.Error("The x-gradients don't match the expected values")
	}
}
