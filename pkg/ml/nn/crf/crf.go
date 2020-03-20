// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crf

import (
	"github.com/saientist/spago/pkg/mat"
	"github.com/saientist/spago/pkg/ml/ag"
	"github.com/saientist/spago/pkg/ml/nn"
	"io"
	"log"
)

type Model struct {
	TransitionScores *nn.Param `type:"weights"`
}

func New(labels int) *Model {
	return &Model{
		TransitionScores: nn.NewParam(mat.NewEmptyDense(labels+1, labels+1)), // +1 for start and end transitions
	}
}

func (m *Model) ForEachParam(callback func(param *nn.Param)) {
	nn.ForEachParam(m, callback)
}

func (m *Model) Serialize(w io.Writer) (int, error) {
	return nn.Serialize(m, w)
}

func (m *Model) Deserialize(r io.Reader) (int, error) {
	return nn.Deserialize(m, r)
}

type Processor struct {
	opt              []interface{}
	model            *Model
	g                *ag.Graph
	transitionScores ag.Node
}

func (m *Model) NewProc(g *ag.Graph, opt ...interface{}) nn.Processor {
	p := &Processor{
		model:            m,
		opt:              opt,
		g:                g,
		transitionScores: g.NewWrap(m.TransitionScores),
	}
	p.init(opt)
	return p
}

func (p *Processor) init(opt []interface{}) {
	if len(opt) > 0 {
		log.Fatal("crf: invalid init option")
	}
}

func (p *Processor) Model() nn.Model       { return p.model }
func (p *Processor) Graph() *ag.Graph      { return p.g }
func (p *Processor) RequiresFullSeq() bool { return true }

func (p *Processor) Reset() {
	p.init(p.opt)
}

func (p *Processor) Forward(_ ...ag.Node) []ag.Node {
	panic("crf: Forward() not available. Use Predict() instead.")
}

func (p *Processor) Predict(emissionScores []ag.Node) []int {
	return Viterbi(p.model.TransitionScores.Value(), emissionScores)
}

func (p *Processor) NegativeLogLoss(emissionScores []ag.Node, target []int) ag.Node {
	goldScore := p.goldScore(emissionScores, target)
	totalScore := p.totalScore(emissionScores)
	return p.g.Sub(totalScore, goldScore)
}

func (p *Processor) goldScore(emissionScores []ag.Node, target []int) ag.Node {
	goldScore := p.g.At(emissionScores[0], target[0], 0)
	goldScore = p.g.Add(goldScore, p.g.At(p.transitionScores, 0, target[0]+1)) // start transition
	prevIndex := target[0] + 1
	for i := 1; i < len(target); i++ {
		goldScore = p.g.Add(goldScore, p.g.At(emissionScores[i], target[i], 0))
		goldScore = p.g.Add(goldScore, p.g.At(p.transitionScores, prevIndex, target[i]+1))
		prevIndex = target[i] + 1
	}
	goldScore = p.g.Add(goldScore, p.g.At(p.transitionScores, prevIndex, 0)) // end transition
	return goldScore
}

func (p *Processor) totalScore(predicted []ag.Node) ag.Node {
	totalVector := p.totalScoreStart(predicted[0])
	for i := 1; i < len(predicted); i++ {
		totalVector = p.totalScoreStep(totalVector, predicted[i])
	}
	totalVector = p.totalScoreEnd(totalVector)
	return p.g.Log(p.g.ReduceSum(totalVector))

}

func (p *Processor) totalScoreStart(stepVec ag.Node) ag.Node {
	size := p.model.TransitionScores.Value().Rows() - 1
	scores := make([]ag.Node, size)
	for i := 0; i < size; i++ {
		scores[i] = p.g.Add(p.g.AtVec(stepVec, i), p.g.At(p.transitionScores, 0, i+1))
	}
	return p.g.Concat(scores...)
}

func (p *Processor) totalScoreEnd(stepVec ag.Node) ag.Node {
	size := p.model.TransitionScores.Value().Rows() - 1
	scores := make([]ag.Node, size)
	for i := 0; i < size; i++ {
		vecTrans := p.g.Add(p.g.AtVec(stepVec, i), p.g.At(p.transitionScores, i+1, 0))
		scores[i] = p.g.Add(scores[i], p.g.Exp(vecTrans))
	}
	return p.g.Concat(scores...)
}

func (p *Processor) totalScoreStep(totalVec ag.Node, stepVec ag.Node) ag.Node {
	size := p.model.TransitionScores.Value().Rows() - 1
	scores := make([]ag.Node, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			vecSum := p.g.Add(p.g.AtVec(totalVec, i), p.g.AtVec(stepVec, j))
			vecTrans := p.g.Add(vecSum, p.g.At(p.transitionScores, i+1, j+1))
			scores[j] = p.g.Add(scores[j], p.g.Exp(vecTrans))
		}
	}
	for i := 0; i < size; i++ {
		scores[i] = p.g.Log(scores[i])
	}
	return p.g.Concat(scores...)
}
