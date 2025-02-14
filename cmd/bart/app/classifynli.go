// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"github.com/nlpodyssey/spago/cmd/clientutils"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/server/grpcapi"
	"github.com/urfave/cli/v2"
	"strings"
)

func newClientClassifyNLICommandFor(app *BartApp) *cli.Command {
	return &cli.Command{
		Name:        "classify-nli",
		Usage:       "Perform zero-shot classification using BART fine-tuned for Natural Language Inference (NLI).",
		Description: "Run the " + programName + " client to perform zero-shot classification.",
		Flags:       newClientClassifyNLICommandFlagsFor(app),
		Action:      newClientClassifyNLICommandActionFor(app),
	}
}

func newClientClassifyNLICommandFlagsFor(app *BartApp) []cli.Flag {
	return clientutils.Flags(&app.grpcAddress, &app.output, []cli.Flag{
		&cli.StringFlag{
			Name:        "text",
			Destination: &app.requestText,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "labels",
			Usage:       "candidate labels separated by `,`",
			Destination: &app.commaSepLabels,
			Required:    true,
		},
		&cli.BoolFlag{
			Name:        "multi-class",
			Destination: &app.multiClass,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "hypothesis-template",
			Destination: &app.requestText2,
			Required:    false,
		},
	})
}

func newClientClassifyNLICommandActionFor(app *BartApp) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		clientutils.VerifyFlags(app.output)

		conn := clientutils.OpenConnection(app.grpcAddress)
		client := grpcapi.NewBARTClient(conn)

		resp, err := client.ClassifyNLI(context.Background(), &grpcapi.ClassifyNLIRequest{
			Text:               app.requestText,
			HypothesisTemplate: app.requestText2,
			PossibleLabels: func() []string {
				labels := make([]string, 0)
				for _, x := range strings.Split(app.commaSepLabels, ",") {
					labels = append(labels, strings.Trim(x, " "))
				}
				return labels
			}(),
			MultiClass: app.multiClass,
		})
		if err != nil {
			return err
		}

		clientutils.Println(app.output, resp)
		return nil
	}
}
