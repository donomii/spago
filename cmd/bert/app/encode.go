// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"log"

	"github.com/nlpodyssey/spago/cmd/clientutils"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bert/grpcapi"
	"github.com/urfave/cli/v2"
)

func newClientEncodeCommandFor(app *BertApp) *cli.Command {
	return &cli.Command{
		Name:        "encode",
		Usage:       "Perform sentence2vec encoding using BERT.",
		Description: "Run the " + programName + " client for sentence encoding.",
		Flags:       newClientEncodeCommandFlagsFor(app),
		Action:      newClientEncodeCommandActionFor(app),
	}
}

func newClientEncodeCommandFlagsFor(app *BertApp) []cli.Flag {
	return clientutils.Flags(&app.address, &app.output, []cli.Flag{
		&cli.StringFlag{
			Name:        "text",
			Destination: &app.requestText,
			Required:    true,
		},
	})
}

func newClientEncodeCommandActionFor(app *BertApp) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		clientutils.VerifyFlags(app.output)

		conn := clientutils.OpenConnection(app.address)
		client := grpcapi.NewBERTClient(conn)

		resp, err := client.Encode(context.Background(), &grpcapi.EncodeRequest{
			Text: app.requestText,
		})

		if err != nil {
			log.Fatalln(err)
		}

		clientutils.Println(app.output, resp)

		return nil
	}
}
