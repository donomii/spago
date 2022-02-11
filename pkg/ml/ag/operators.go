// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ag

import (
	"fmt"
	"github.com/nlpodyssey/spago/pkg/ml/ag/fn"
	"reflect"
	"strings"
)

// OpName is the enumeration-like type used for the set of operators supported
// by spaGO.
type OpName int

const (
	// OpIdentity identifies the Graph.Identity operator.
	OpIdentity OpName = iota
	// OpDropout identifies the Graph.Dropout operator.
	OpDropout
	// OpAtVec identifies the Graph.AtVec operator.
	OpAtVec
	// OpAt identifies the Graph.At operator.
	OpAt
	// OpAdd identifies the Graph.Add operator.
	OpAdd
	// OpSub identifies the Graph.Sub operator.
	OpSub
	// OpSubScalar identifies the Graph.SubScalar operator.
	OpSubScalar
	// OpAddScalar identifies the Graph.AddScalar operator.
	OpAddScalar
	// OpReverseSub identifies the Graph.ReverseSub operator.
	OpReverseSub
	// OpProd identifies the Graph.Prod operator.
	OpProd
	// OpDiv identifies the Graph.Div operator.
	OpDiv
	// OpProdScalar identifies the Graph.ProdScalar operator.
	OpProdScalar
	// OpDivScalar identifies the Graph.DivScalar operator.
	OpDivScalar
	// OpMul identifies the Graph.Mul operator.
	OpMul
	// OpDot identifies the Graph.Dot operator.
	OpDot
	// OpReshape identifies the Graph.Reshape operator.
	OpReshape
	// OpMaxPooling identifies the Graph.MaxPooling operator.
	OpMaxPooling
	// OpView identifies the Graph.View operator.
	OpView
	// OpRowView identifies the Graph.RowView operator.
	OpRowView
	// OpColView identifies the Graph.ColView operator.
	OpColView
	// OpVec identifies the Graph.Vec operator.
	OpVec
	// OpRotateR identifies the Graph.RotateR operator.
	OpRotateR
	// OpT identifies the Graph.T operator.
	OpT
	// OpSquare identifies the Graph.Square operator.
	OpSquare
	// OpPow identifies the Graph.Pow operator.
	OpPow
	// OpSqrt identifies the Graph.Sqrt operator.
	OpSqrt
	// OpTan identifies the Graph.Tan operator.
	OpTan
	// OpTanh identifies the Graph.Tanh operator.
	OpTanh
	// OpSigmoid identifies the Graph.Sigmoid operator.
	OpSigmoid
	// OpHardSigmoid identifies the Graph.HardSigmoid operator.
	OpHardSigmoid
	// OpHardTanh identifies the Graph.HardTanh operator.
	OpHardTanh
	// OpSoftsign identifies the Graph.Softsign operator.
	OpSoftsign
	// OpReLU identifies the Graph.ReLU operator.
	OpReLU
	// OpCELU identifies the Graph.CELU operator.
	OpCELU
	// OpGELU identifies the Graph.GELU operator.
	OpGELU
	// OpELU identifies the Graph.ELU operator.
	OpELU
	// OpPositiveELU identifies the Graph.PositiveELU operator.
	OpPositiveELU
	// OpSwishB identifies the Graph.SwishB operator.
	OpSwishB
	// OpSwish identifies the Graph.Swish operator.
	OpSwish
	// OpSiLU identifies the Graph.SiLU operator.
	OpSiLU
	// OpMish identifies the Graph.Mish operator.
	OpMish
	// OpLeakyReLU identifies the Graph.LeakyReLU operator.
	OpLeakyReLU
	// OpSELU identifies the Graph.SELU operator.
	OpSELU
	// OpSoftPlus identifies the Graph.SoftPlus operator.
	OpSoftPlus
	// OpSoftShrink identifies the Graph.SoftShrink operator.
	OpSoftShrink
	// OpThreshold identifies the Graph.Threshold operator.
	OpThreshold
	// OpSoftmax identifies the Graph.Softmax operator.
	OpSoftmax
	// OpLogSoftmax identifies the Graph.LogSoftmax operator.
	OpLogSoftmax
	// OpSparseMax identifies the Graph.SparseMax operator.
	OpSparseMax
	// OpSparseMaxLoss identifies the Graph.SparseMaxLoss operator.
	OpSparseMaxLoss
	// OpSin identifies the Graph.Sin operator.
	OpSin
	// OpCos identifies the Graph.Cos operator.
	OpCos
	// OpExp identifies the Graph.Exp operator.
	OpExp
	// OpLog identifies the Graph.Log operator.
	OpLog
	// OpAbs identifies the Graph.Abs operator.
	OpAbs
	// OpNeg identifies the Graph.Neg operator.
	OpNeg
	// OpReciprocal identifies the Graph.Reciprocal operator.
	OpReciprocal
	// OpMax identifies the Graph.Max operator.
	OpMax
	// OpMin identifies the Graph.Min operator.
	OpMin
	// OpReduceSum identifies the Graph.ReduceSum operator.
	OpReduceSum
	// OpReduceMean identifies the Graph.ReduceMean operator.
	OpReduceMean
	// OpMean identifies the Graph.Mean operator.
	OpMean
	// OpSum identifies the Graph.Sum operator.
	OpSum
	// OpConcat identifies the Graph.Concat operator.
	OpConcat
	// OpStack identifies the Graph.Stack operator.
	OpStack
)

var opNameToMethodName = map[OpName]string{
	OpIdentity:      "Identity",
	OpDropout:       "Dropout",
	OpAtVec:         "AtVec",
	OpAt:            "At",
	OpAdd:           "Add",
	OpSub:           "Sub",
	OpSubScalar:     "SubScalar",
	OpAddScalar:     "AddScalar",
	OpReverseSub:    "ReverseSub",
	OpProd:          "Prod",
	OpDiv:           "Div",
	OpProdScalar:    "ProdScalar",
	OpDivScalar:     "DivScalar",
	OpMul:           "Mul",
	OpDot:           "Dot",
	OpReshape:       "Reshape",
	OpMaxPooling:    "MaxPooling",
	OpView:          "View",
	OpRowView:       "RowView",
	OpColView:       "ColView",
	OpVec:           "Vec",
	OpRotateR:       "RotateR",
	OpT:             "T",
	OpSquare:        "Square",
	OpPow:           "Pow",
	OpSqrt:          "Sqrt",
	OpTan:           "Tan",
	OpTanh:          "Tanh",
	OpSigmoid:       "Sigmoid",
	OpHardSigmoid:   "HardSigmoid",
	OpHardTanh:      "HardTanh",
	OpSoftsign:      "Softsign",
	OpReLU:          "ReLU",
	OpCELU:          "CELU",
	OpGELU:          "GELU",
	OpELU:           "ELU",
	OpPositiveELU:   "PositiveELU",
	OpSwishB:        "SwishB",
	OpSwish:         "Swish",
	OpSiLU:          "SiLU",
	OpMish:          "Mish",
	OpLeakyReLU:     "LeakyReLU",
	OpSELU:          "SELU",
	OpSoftPlus:      "SoftPlus",
	OpSoftShrink:    "SoftShrink",
	OpThreshold:     "Threshold",
	OpSoftmax:       "Softmax",
	OpLogSoftmax:    "LogSoftmax",
	OpSparseMax:     "SparseMax",
	OpSparseMaxLoss: "SparseMaxLoss",
	OpSin:           "Sin",
	OpCos:           "Cos",
	OpExp:           "Exp",
	OpLog:           "Log",
	OpAbs:           "Abs",
	OpNeg:           "Neg",
	OpReciprocal:    "Reciprocal",
	OpMax:           "Max",
	OpMin:           "Min",
	OpReduceSum:     "ReduceSum",
	OpReduceMean:    "ReduceMean",
	OpMean:          "Mean",
	OpSum:           "Sum",
	OpConcat:        "Concat",
	OpStack:         "Stack",
}

// strToOpName is the inverse map of opNameToMethodName.
var strToOpName = func() map[string]OpName {
	invMap := make(map[string]OpName)
	for k, v := range opNameToMethodName {
		invMap[v] = k
		invMap[strings.ToLower(v)] = k
	}
	return invMap
}()

// GetOpName maps a string to an operator.
// It panics if the string does not match any operator (not even using lowercase).
func GetOpName(str string) (OpName, error) {
	if value, ok := strToOpName[str]; ok {
		return value, nil
	}
	return -1, fmt.Errorf("ag: unknown operator %s", str)
}

// Invoke returns a new node as a result of the application of the input operator.
func (g *Graph[T]) Invoke(operator OpName, xs ...Node[T]) Node[T] {
	v := reflect.ValueOf(g).MethodByName(opNameToMethodName[operator])
	args := make([]reflect.Value, len(xs))
	for i, x := range xs {
		args[i] = reflect.ValueOf(x)
	}
	ret := v.Call(args)
	return ret[0].Interface().(Node[T])
}

// Identity returns a new operator node as a result of the fn.Identity function.
func (g *Graph[T]) Identity(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewIdentity[T](x), x)
}

// Dropout returns a new operator node as a result of the fn.Dropout function.
func (g *Graph[T]) Dropout(x Node[T], p T) Node[T] {
	return g.NewOperator(fn.NewDropout[T](x, p, g.randGen), x)
}

// AtVec returns a new operator node as a result of the fn.AtVec function.
func (g *Graph[T]) AtVec(x Node[T], i int) Node[T] {
	return g.NewOperator(fn.NewAtVec[T](x, i), x)
}

// At returns a new operator node as a result of the fn.At function.
func (g *Graph[T]) At(x Node[T], i int, j int) Node[T] {
	return g.NewOperator(fn.NewAt[T](x, i, j), x)
}

// Add returns a new operator node as a result of the fn.Add function.
// The first node may be null. This help to keep the code as concise as possible e.g. during accumulation.
func (g *Graph[T]) Add(x1 Node[T], x2 Node[T]) Node[T] {
	if x1 != nil {
		return g.NewOperator(fn.NewAdd[T](x1, x2), x1, x2)
	}
	fake := g.NewVariable(nil, false)
	return g.NewOperator(fn.NewAdd[T](fake, x2), fake, x2)
}

// Sub returns a new operator node as a result of the fn.Sub function.
func (g *Graph[T]) Sub(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewSub[T](x1, x2), x1, x2)
}

// SubScalar returns a new operator node as a result of the fn.SubScalar function.
func (g *Graph[T]) SubScalar(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewSubScalar[T](x1, x2), x1, x2)
}

// AddScalar returns a new operator node as a result of the fn.AddScalar function.
func (g *Graph[T]) AddScalar(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewAddScalar[T](x1, x2), x1, x2)
}

// ReverseSub returns a new operator node as a result of the fn.ReverseSub function.
func (g *Graph[T]) ReverseSub(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewReverseSubScalar[T](x1, x2), x1, x2)
}

// Prod returns a new operator node as a result of the fn.Prod function.
func (g *Graph[T]) Prod(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewProd[T](x1, x2), x1, x2)
}

// Div returns a new operator node as a result of the fn.Div function.
func (g *Graph[T]) Div(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewDiv[T](x1, x2), x1, x2)
}

// ProdScalar returns a new operator node as a result of the fn.ProdScalar function.
func (g *Graph[T]) ProdScalar(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewProdScalar[T](x1, x2), x1, x2)
}

// DivScalar returns a new operator node as a result of the fn.DivScalar function.
func (g *Graph[T]) DivScalar(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewDivScalar[T](x1, x2), x1, x2)
}

// Mul returns a new operator node as a result of the fn.Mul function.
func (g *Graph[T]) Mul(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewMul[T](x1, x2), x1, x2)
}

// Dot returns a new operator node as a result of the fn.Dot function.
func (g *Graph[T]) Dot(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewDot[T](x1, x2), x1, x2)
}

// Max returns a new operator node as a result of the fn.Max function.
func (g *Graph[T]) Max(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewMax[T](x1, x2), x1, x2)
}

// Min returns a new operator node as a result of the fn.Min function.
func (g *Graph[T]) Min(x1, x2 Node[T]) Node[T] {
	return g.NewOperator(fn.NewMin[T](x1, x2), x1, x2)
}

// Reshape returns a new operator node as a result of the fn.Reshape function.
func (g *Graph[T]) Reshape(x Node[T], rows, columns int) Node[T] {
	return g.NewOperator(fn.NewReshape[T](x, rows, columns), x)
}

// MaxPooling returns a new operator node as a result of the fn.MaxPooling function.
func (g *Graph[T]) MaxPooling(x Node[T], rows, columns int) Node[T] {
	return g.NewOperator(fn.NewMaxPooling[T](x, rows, columns), x)
}

// View returns a new operator node as a result of the fn.View function.
func (g *Graph[T]) View(x Node[T], row, column, xStride, yStride int) Node[T] {
	return g.NewOperator(fn.NewView[T](x, row, column, xStride, yStride), x)
}

// RowView returns a new operator node as a result of the fn.RowView function.
func (g *Graph[T]) RowView(x Node[T], row int) Node[T] {
	return g.NewOperator(fn.NewRowView[T](x, row), x)
}

// RotateR performs the right circular shift.
// `i` is the number of places by which the elements are shifted.
func (g *Graph[T]) RotateR(x Node[T], i int) Node[T] {
	return g.NewOperator(fn.NewRotateR[T](x, i), x)
}

// ColView returns a new operator node as a result of the fn.ColView function.
func (g *Graph[T]) ColView(x Node[T], column int) Node[T] {
	return g.NewOperator(fn.NewColView[T](x, column), x)
}

// Vec returns a new operator node as a result of the fn.Vec function.
func (g *Graph[T]) Vec(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewVec[T](x), x)
}

// T returns a new operator node as a result of the fn.T function.
func (g *Graph[T]) T(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewTranspose[T](x), x)
}

// Square returns a new operator node as a result of the fn.Prod(x, x) function.
func (g *Graph[T]) Square(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSquare[T](x), x)
}

// Pow returns a new operator node as a result of the fn.Pow function.
func (g *Graph[T]) Pow(x Node[T], power T) Node[T] {
	return g.NewOperator(fn.NewPow[T](x, power), x)
}

// Sqrt returns a new operator node as a result of the `Sqrt` function.
func (g *Graph[T]) Sqrt(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSqrt[T](x), x)
}

// Tan returns a new operator node as a result of the `Tan` function.
func (g *Graph[T]) Tan(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewTan[T](x), x)
}

// Tanh returns a new operator node as a result of the `Tanh` function.
func (g *Graph[T]) Tanh(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewTanh[T](x), x)
}

// Sigmoid returns a new operator node as a result of the `Sigmoid` function.
func (g *Graph[T]) Sigmoid(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSigmoid[T](x), x)
}

// HardSigmoid returns a new operator node as a result of the `HardSigmoid` function.
func (g *Graph[T]) HardSigmoid(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewHardSigmoid[T](x), x)
}

// HardTanh returns a new operator node as a result of the `HardTanh` function.
func (g *Graph[T]) HardTanh(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewHardTanh[T](x), x)
}

// Softsign returns a new operator node as a result of the `SoftSign` function.
func (g *Graph[T]) Softsign(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSoftsign[T](x), x)
}

// ReLU returns a new operator node as a result of the `ReLU` function.
func (g *Graph[T]) ReLU(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewReLU[T](x), x)
}

// CELU returns a new operator node as a result of the fn.CELU function.
func (g *Graph[T]) CELU(x, alpha Node[T]) Node[T] {
	return g.NewOperator(fn.NewCELU[T](x, alpha), x, alpha)
}

// GELU returns a new operator node as a result of the fn.GELU function.
func (g *Graph[T]) GELU(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewGELU[T](x), x)
}

// ELU returns a new operator node as a result of the fn.ELU function.
func (g *Graph[T]) ELU(x, alpha Node[T]) Node[T] {
	return g.NewOperator(fn.NewELU[T](x, alpha), x, alpha)
}

// SwishB returns a new operator node as a result of the fn.SwishB function.
func (g *Graph[T]) SwishB(x, beta Node[T]) Node[T] {
	return g.NewOperator(fn.NewSwishB[T](x, beta), x, beta)
}

// Swish returns a new operator node as a result of the fn.Swish function.
func (g *Graph[T]) Swish(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSwish[T](x), x)
}

// SiLU returns a new operator node as a result of the fn.SiLU function.
func (g *Graph[T]) SiLU(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSiLU[T](x), x)
}

// Mish returns a new operator node as a result of the `Mish` function.
func (g *Graph[T]) Mish(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewMish[T](x), x)
}

// LeakyReLU returns a new operator node as a result of the fn.LeakyReLU function.
func (g *Graph[T]) LeakyReLU(x, alpha Node[T]) Node[T] {
	return g.NewOperator(fn.NewLeakyReLU[T](x, alpha), x, alpha)
}

// SELU returns a new operator node as a result of the fn.SELU function.
func (g *Graph[T]) SELU(x, alpha Node[T], scale Node[T]) Node[T] {
	return g.NewOperator(fn.NewSELU[T](x, alpha, scale), x, alpha, scale)
}

// SoftPlus returns a new operator node as a result of the fn.SoftPlus function.
func (g *Graph[T]) SoftPlus(x, beta, threshold Node[T]) Node[T] {
	return g.NewOperator(fn.NewSoftPlus[T](x, beta, threshold), x, beta, threshold)
}

// SoftShrink returns a new operator node as a result of the fn.SoftShrink function.
func (g *Graph[T]) SoftShrink(x, lambda Node[T]) Node[T] {
	return g.NewOperator(fn.NewSoftShrink[T](x, lambda), x, lambda)
}

// Threshold returns a new operator node as a result of the fn.Threshold function.
func (g *Graph[T]) Threshold(x, threshold, k Node[T]) Node[T] {
	return g.NewOperator(fn.NewThreshold[T](x, threshold, k), x, threshold, k)
}

// Softmax returns a new operator node as a result of the fn.Softmax function.
func (g *Graph[T]) Softmax(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSoftmax[T](x), x)
}

// SparseMax returns a new operator node as a result of the fn.SparseMax function.
func (g *Graph[T]) SparseMax(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSparseMax[T](x), x)
}

// SparseMaxLoss returns a new operator node as a result of the fn.SparseMaxLoss function.
func (g *Graph[T]) SparseMaxLoss(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSparseMaxLoss[T](x), x)
}

// Sin returns a new operator node as a result of the `Sin` function.
func (g *Graph[T]) Sin(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewSin[T](x), x)
}

// Cos returns a new operator node as a result of the `Cos` function.
func (g *Graph[T]) Cos(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewCos[T](x), x)
}

// Exp returns a new operator node as a result of the `Exp` function.
func (g *Graph[T]) Exp(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewExp[T](x), x)
}

// Log returns a new operator node as a result of the `Log` function.
func (g *Graph[T]) Log(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewLog[T](x), x)
}

// Abs returns a new operator node as a result of the `Abs` function.
func (g *Graph[T]) Abs(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewAbs[T](x), x)
}

// Neg returns a new operator node as a result of the `Neg` function.
func (g *Graph[T]) Neg(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewNeg[T](x), x)
}

// Reciprocal returns a new operator node as a result of the `Reciprocal` function.
func (g *Graph[T]) Reciprocal(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewReciprocal[T](x), x)
}

// ReduceSum returns a new operator node as a result of the fn.ReduceSum function.
func (g *Graph[T]) ReduceSum(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewReduceSum[T](x), x)
}

// ReduceMean returns a new operator node as a result of the fn.ReduceMean function.
func (g *Graph[T]) ReduceMean(x Node[T]) Node[T] {
	return g.NewOperator(fn.NewReduceMean[T](x), x)
}

// Concat returns a new operator node as a result of the fn.Concat function.
func (g *Graph[T]) Concat(xs ...Node[T]) Node[T] {
	return g.NewOperator(fn.NewConcat(Operands(xs)), xs...)
}

// Stack returns a new operator node as a result of the fn.Stack function.
func (g *Graph[T]) Stack(xs ...Node[T]) Node[T] {
	return g.NewOperator(fn.NewStack(Operands(xs)), xs...)
}
