// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat32

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPow(t *testing.T) {
	assert.InDelta(t, Float(8), Pow(2, 3), 1e-10)
}

func TestCos(t *testing.T) {
	assert.InDelta(t, Float(-1), Cos(Pi), 1e-10)
}

func TestSin(t *testing.T) {
	assert.InDelta(t, Float(1), Sin(Pi/2), 1e-10)
}

func TestCosh(t *testing.T) {
	assert.InDelta(t, 11.59195, Cosh(Pi), 1e-5)
}

func TestSinh(t *testing.T) {
	assert.InDelta(t, 11.54874, Sinh(Pi), 1e-5)
}

func TestExp(t *testing.T) {
	assert.InDelta(t, 2.71828, Exp(1), 1e-5)
}

func TestAbs(t *testing.T) {
	assert.Equal(t, Float(42), Abs(42))
	assert.Equal(t, Float(42), Abs(-42))
}

func TestSqrt(t *testing.T) {
	assert.InDelta(t, Float(3), Sqrt(9), 1e-10)
}

func TestLog(t *testing.T) {
	assert.InDelta(t, 0.69314, Log(2), 1e-5)
}

func TestTan(t *testing.T) {
	assert.InDelta(t, Float(1), Tan(Pi/4), 1e-10)
}

func TestTanh(t *testing.T) {
	assert.InDelta(t, 0.76159, Tanh(1), 1e-5)
}

func TestMax(t *testing.T) {
	assert.Equal(t, Float(2), Max(1, 2))
	assert.Equal(t, Float(2), Max(2, 1))
}

func TestInf(t *testing.T) {
	assert.True(t, math.IsInf(float64(Inf(1)), +1))
	assert.True(t, math.IsInf(float64(Inf(-1)), -1))
}

func TestIsInf(t *testing.T) {
	assert.True(t, IsInf(Inf(1), +1))
	assert.True(t, IsInf(Inf(-1), -1))
	assert.False(t, IsInf(Inf(-1), +1))
	assert.False(t, IsInf(Inf(1), -1))
	assert.False(t, IsInf(Float(0), 1))
	assert.False(t, IsInf(Float(0), -1))
}

func TestNaN(t *testing.T) {
	assert.True(t, math.IsNaN(float64(NaN())))
}

func TestCeil(t *testing.T) {
	assert.Equal(t, Float(2), Ceil(1.2))
}

func TestFloor(t *testing.T) {
	assert.Equal(t, Float(1), Floor(1.9))
}

func TestRound(t *testing.T) {
	assert.Equal(t, Float(1), Round(1.4))
	assert.Equal(t, Float(2), Round(1.5))
	assert.Equal(t, Float(2), Round(1.6))
}
