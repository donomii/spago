// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdmbuilder

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd/adagrad"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd/adam"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd/lamb"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd/radam"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd/rmsprop"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd/sgd"
)

// NewMethod returns a new gd.Method, chosen and initialized according to
// the given config.
// It panics if the config type is unknown or unsupported.
func NewMethod[T mat.DType](config gd.MethodConfig) gd.Method[T] {
	switch config := config.(type) {
	case adagrad.Config[T]:
		return adagrad.New(config)
	case adam.Config[T]:
		return adam.New(config)
	case radam.Config[T]:
		return radam.New(config)
	case rmsprop.Config[T]:
		return rmsprop.New(config)
	case lamb.Config[T]:
		return lamb.New(config)
	case sgd.Config[T]:
		return sgd.New(config)
	default:
		panic("gd: unknown method configuration")
	}
}
