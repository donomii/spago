// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package losses

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
)

// MAE measures the mean absolute error (a.k.a. L1 Loss) between each element in the input x and target y.
func MAE[T mat.DType](g *ag.Graph[T], x ag.Node[T], y ag.Node[T], reduceMean bool) ag.Node[T] {
	loss := g.Abs(g.Sub(x, y))
	if reduceMean {
		return g.ReduceMean(loss)
	}
	return g.ReduceSum(loss)
}

// MSE measures the mean squared error (squared L2 norm) between each element in the input x and target y.
func MSE[T mat.DType](g *ag.Graph[T], x ag.Node[T], y ag.Node[T], reduceMean bool) ag.Node[T] {
	loss := g.ProdScalar(g.Square(g.Sub(x, y)), g.Constant(0.5))
	if reduceMean {
		return g.ReduceMean(loss)
	}
	return g.ReduceSum(loss)
}

// NLL returns the loss of the input x respect to the target y.
// The target is expected to be a one-hot vector.
func NLL[T mat.DType](g *ag.Graph[T], x ag.Node[T], y ag.Node[T]) ag.Node[T] {
	return g.Neg(g.ReduceSum(g.Prod(y, g.Log(x))))
}

// CrossEntropy implements a cross-entropy loss function.
// x is the raw scores for each class (logits).
// c is the index of the gold class.
func CrossEntropy[T mat.DType](g *ag.Graph[T], x ag.Node[T], c int) ag.Node[T] {
	return g.Add(g.Neg(g.AtVec(x, c)), g.Log(g.ReduceSum(g.Exp(x))))
}

// WeightedCrossEntropy implements a weighted cross-entropy loss function.
// x is the raw scores for each class (logits).
// c is the index of the gold class.
// This function is scaled by a weighting factor weights[class] ∈ [0,1]
func WeightedCrossEntropy[T mat.DType](weights []T) func(g *ag.Graph[T], x ag.Node[T], c int) ag.Node[T] {
	return func(g *ag.Graph[T], x ag.Node[T], c int) ag.Node[T] {
		return g.ProdScalar(CrossEntropy(g, x, c), g.NewScalar(weights[c]))
	}
}

// FocalLoss implements a variant of the CrossEntropy loss that reduces
// the loss contribution from "easy" examples and increases the importance
// of correcting misclassified examples.
// x is the raw scores for each class (logits).
// c is the index of the gold class.
// gamma is the focusing parameter (gamma ≥ 0).
func FocalLoss[T mat.DType](g *ag.Graph[T], x ag.Node[T], c int, gamma T) ag.Node[T] {
	ce := CrossEntropy(g, x, c)
	p := g.Exp(g.Neg(ce))
	sub := g.ReverseSub(p, g.NewScalar(1.0))
	a := g.Pow(sub, gamma)
	return g.Prod(a, ce)
}

// WeightedFocalLoss implements a variant of the CrossEntropy loss that reduces
// the loss contribution from "easy" examples and increases the importance
// of correcting misclassified examples.
// x is the raw scores for each class (logits).
// c is the index of the gold class.
// gamma is the focusing parameter (gamma ≥ 0).
// This function is scaled by a weighting factor weights[class] ∈ [0,1].
func WeightedFocalLoss[T mat.DType](weights []T) func(g *ag.Graph[T], x ag.Node[T], c int, gamma T) ag.Node[T] {
	return func(g *ag.Graph[T], x ag.Node[T], c int, gamma T) ag.Node[T] {
		ce := CrossEntropy(g, x, c)
		p := g.Exp(g.Neg(ce))
		sub := g.ReverseSub(p, g.NewScalar(1.0))
		b := g.Pow(sub, gamma)
		fl := g.Prod(b, ce)
		return g.ProdScalar(fl, g.NewScalar(weights[c]))
	}
}

// Perplexity computes the perplexity, implemented as exp over the cross-entropy.
func Perplexity[T mat.DType](g *ag.Graph[T], x ag.Node[T], c int) ag.Node[T] {
	return g.Exp(CrossEntropy(g, x, c))
}

// ZeroOneQuantization is a loss function that is minimized when each component
// of x satisfies x(i) ≡ [x]i ∈ {0, 1}.
func ZeroOneQuantization[T mat.DType](g *ag.Graph[T], x ag.Node[T]) ag.Node[T] {
	return g.ReduceSum(g.Prod(g.Square(x), g.Square(g.ReverseSub(x, g.NewScalar(1.0)))))
}

// Norm2Quantization is a loss function that is minimized when norm2(x) = 1.
func Norm2Quantization[T mat.DType](g *ag.Graph[T], x ag.Node[T]) ag.Node[T] {
	return g.Square(g.SubScalar(g.ReduceSum(g.Square(x)), g.NewScalar(1.0)))
}

// OneHotQuantization is a loss function that pushes towards the x vector to be 1-hot.
// q is the quantization regularizer weight (suggested  0.00001).
func OneHotQuantization[T mat.DType](g *ag.Graph[T], x ag.Node[T], q T) ag.Node[T] {
	return g.ProdScalar(g.Add(ZeroOneQuantization(g, x), Norm2Quantization(g, x)), g.NewScalar(q))
}

// Distance is a loss function that calculates the distance between target and x.
func Distance[T mat.DType](g *ag.Graph[T], x ag.Node[T], target T) ag.Node[T] {
	return g.Abs(g.Sub(g.NewScalar(target), x))
}

// MSESeq calculates the MSE loss on the given sequence.
func MSESeq[T mat.DType](g *ag.Graph[T], predicted []ag.Node[T], target []ag.Node[T], reduceMean bool) ag.Node[T] {
	loss := MSE(g, predicted[0], target[0], false)
	for i := 1; i < len(predicted); i++ {
		loss = g.Add(loss, MSE(g, predicted[i], target[i], false))
	}
	if reduceMean {
		return g.DivScalar(loss, g.NewScalar(T(len(predicted))))
	}
	return loss
}

// MAESeq calculates the MAE loss on the given sequence.
func MAESeq[T mat.DType](g *ag.Graph[T], predicted []ag.Node[T], target []ag.Node[T], reduceMean bool) ag.Node[T] {
	loss := MAE(g, predicted[0], target[0], false)
	for i := 1; i < len(predicted); i++ {
		loss = g.Add(loss, MAE(g, predicted[i], target[i], false))
	}
	if reduceMean {
		return g.DivScalar(loss, g.NewScalar(T(len(predicted))))
	}
	return loss
}

// CrossEntropySeq calculates the CrossEntropy loss on the given sequence.
func CrossEntropySeq[T mat.DType](g *ag.Graph[T], predicted []ag.Node[T], target []int, reduceMean bool) ag.Node[T] {
	loss := CrossEntropy(g, predicted[0], target[0])
	for i := 1; i < len(predicted); i++ {
		loss = g.Add(loss, CrossEntropy(g, predicted[i], target[i]))
	}
	if reduceMean {
		return g.DivScalar(loss, g.NewScalar(T(len(predicted))))
	}
	return loss
}

// SPG (Softmax Policy Gradient) is a Gradient Policy used in Reinforcement Learning.
// logPropActions are the log-probability of the chosen action by the Agent at each time;
// logProbTargets are results of the reward function i.e. the predicted log-likelihood of the ground truth at each time;
func SPG[T mat.DType](g *ag.Graph[T], logPropActions []ag.Node[T], logProbTargets []ag.Node[T]) ag.Node[T] {
	var loss ag.Node[T]
	for t := 0; t < len(logPropActions); t++ {
		loss = g.Add(loss, g.Prod(logPropActions[t], logProbTargets[t]))
	}
	return g.Neg(loss)
}
