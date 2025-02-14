// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"github.com/nlpodyssey/spago/cmd/clientutils"
	mat "github.com/nlpodyssey/spago/pkg/mat32"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bert/grpcapi"
	"github.com/urfave/cli/v2"
	"math"
)

func newClientSimilarityCommandFor(app *BertApp) *cli.Command {
	return &cli.Command{
		Name:        "similarity",
		Usage:       "Perform text-similarity using BERT sentence encoding.",
		Description: "Run the " + programName + " client to determine the similarity between two texts.",
		Flags:       newClientSimilarityCommandFlagsFor(app),
		Action:      newClientSimilarityCommandActionFor(app),
	}
}

func newClientSimilarityCommandFlagsFor(app *BertApp) []cli.Flag {
	return clientutils.Flags(&app.address, &app.output, []cli.Flag{
		&cli.StringFlag{
			Name:        "text1",
			Destination: &app.requestText,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "text2",
			Destination: &app.requestText2,
			Required:    true,
		},
	})
}

func newClientSimilarityCommandActionFor(app *BertApp) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		clientutils.VerifyFlags(app.output)

		conn := clientutils.OpenConnection(app.address)
		client := grpcapi.NewBERTClient(conn)

		resp, err := client.Encode(context.Background(), &grpcapi.EncodeRequest{
			Text: app.requestText,
		})
		if err != nil {
			return err
		}

		resp2, err := client.Encode(context.Background(), &grpcapi.EncodeRequest{
			Text: app.requestText2,
		})
		if err != nil {
			return err
		}

		vec1 := normalize(f32SliceToFloatSlice(resp.Vector))
		vec2 := normalize(f32SliceToFloatSlice(resp2.Vector))
		similarity := vec1.DotUnitary(vec2)

		clientutils.Println(app.output, toFixed(similarity, 6))

		return nil
	}
}

func normalize(xs []mat.Float) mat.Matrix {
	return mat.NewVecDense(xs).Normalize2()
}

func f32SliceToFloatSlice(xs []float32) []mat.Float {
	ys := make([]mat.Float, len(xs))
	for i, f32 := range xs {
		ys[i] = mat.Float(f32)
	}
	return ys
}

func round(num mat.Float) int {
	return int(float64(num) + math.Copysign(0.5, float64(num)))
}

func toFixed(num mat.Float, precision int) mat.Float {
	output := mat.Pow(10, mat.Float(precision))
	return mat.Float(round(num*output)) / output
}
