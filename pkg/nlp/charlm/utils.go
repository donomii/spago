// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package charlm

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/mat/matutils"
	"github.com/nlpodyssey/spago/pkg/mat/rand"
	"github.com/nlpodyssey/spago/pkg/nlp/vocabulary"
)

func targetsIds(sequence []string, vocab *vocabulary.Vocabulary, unknownToken string) []int {
	targetsIds := make([]int, len(sequence)-1) // skip last character
	for i, target := range sequence[1:] {      // the target is always the next character
		id, ok := vocab.ID(target)
		if !ok {
			targetsIds[i] = vocab.MustID(unknownToken)
			continue
		}
		targetsIds[i] = id
	}
	return targetsIds
}

// sample extracts the next character from the probability multinomial distribution.
// Note that the softmax must NOT have been applied to the prediction values.
func sample[T mat.DType](prediction []T, temperature T) int {
	for i := range prediction {
		prediction[i] *= 1.0 / temperature
	}
	prediction = matutils.SoftMax(prediction)
	p := rand.Float[T]() // TODO: use a local random generator?
	for i, x := range prediction {
		p -= x
		if p < 0 {
			return i
		}
	}
	return 0 // TODO: should panic here?
}
