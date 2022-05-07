// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"github.com/nlpodyssey/spago/pkg/nlp/tokenizers/bpetokenizer"
	"github.com/nlpodyssey/spago/pkg/nlp/tokenizers/sentencepiece"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/head/conditionalgeneration"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/head/sequenceclassification"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/loader"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/server"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/huggingface"
	"github.com/nlpodyssey/spago/pkg/utils/httputils"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

func newServerCommandFor(app *BartApp) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Run the " + programName + " as gRPC/HTTP server.",
		Description: "Run the " + programName + " indicating the model path (NOT the model file).",
		Flags:       newServerCommandFlagsFor(app),
		Action:      newServerCommandActionFor(app),
	}
}

func newServerCommandFlagsFor(app *BartApp) []cli.Flag {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-address",
			Usage:       "Changes the bind address of the gRPC server.",
			Value:       "0.0.0.0:1976",
			Destination: &app.grpcAddress,
		},
		&cli.StringFlag{
			Name:        "address",
			Usage:       "Changes the bind address of the HTTP JSON server.",
			Value:       "0.0.0.0:1987",
			Destination: &app.address,
		},
		&cli.StringFlag{
			Name:        "repo",
			Usage:       "Specifies the path to the models.",
			Value:       path.Join(usr.HomeDir, ".spago"),
			Destination: &app.repo,
		},
		&cli.StringFlag{
			Name:        "model, m",
			Required:    true,
			Usage:       "Specifies the model name.",
			Destination: &app.model,
		},
		&cli.IntFlag{
			Name:        "timeout",
			Usage:       "Server read, write, and idle timeout duration in seconds.",
			Value:       httputils.DefaultTimeoutSeconds,
			Destination: &app.serverTimeoutSeconds,
		},
		&cli.IntFlag{
			Name:        "max-request-size",
			Usage:       "Maximum number of bytes the server will read parsing the request content.",
			Value:       httputils.DefaultMaxRequestBytes,
			Destination: &app.serverMaxRequestBytes,
		},
	}
}

func newServerCommandActionFor(app *BartApp) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if err := pullModel(app); err != nil {
			return err
		}

		modelPath := filepath.Join(app.repo, app.model)

		model, err := loader.Load(modelPath)
		if err != nil {
			log.Fatal(err)
		}
		defer model.Close()

		var bpeTokenizer *bpetokenizer.BPETokenizer
		var spTokenizer *sentencepiece.Tokenizer

		switch model.(type) {
		case *sequenceclassification.Model:
			bpeTokenizer, err = bpetokenizer.NewFromModelFolder(modelPath)
			if err != nil {
				return err
			}
		case *conditionalgeneration.Model:
			spTokenizer, err = sentencepiece.NewFromModelFolder(modelPath, false)
			if err != nil {
				return err
			}
		default:
			panic("bart: invalid model type")
		}

		fmt.Printf("Start gRPC server listening on %s.\n", app.grpcAddress)

		fmt.Printf("Start HTTP server listening on %s.\n", app.address)

		s := server.NewServer(model, bpeTokenizer, spTokenizer)
		s.TimeoutSeconds = app.serverTimeoutSeconds
		s.MaxRequestBytes = app.serverMaxRequestBytes
		s.StartDefaultHTTPServer(app.address)
		s.StartDefaultServer(app.grpcAddress)

		return nil
	}
}

const defaultModelFile = "spago_model.bin"

func pullModel(app *BartApp) error {
	modelPath := filepath.Join(app.repo, app.model)
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		fmt.Printf("Unable to find `%s` locally.\n", modelPath)
		fmt.Printf("Pulling `%s` from Hugging Face models hub...\n", app.model)
		// make sure the models path exists
		if _, err := os.Stat(app.repo); os.IsNotExist(err) {
			if err := os.MkdirAll(app.repo, 0755); err != nil {
				return err
			}
		}
		err = huggingface.NewDownloader(app.repo, app.model, false).Download()
		if err != nil {
			return err
		}
		fmt.Printf("Converting model...\n")
		err = huggingface.NewConverter(app.repo, app.model).Convert()
		if err != nil {
			return err
		}

		return nil
	}

	if _, err := os.Stat(path.Join(modelPath, defaultModelFile)); os.IsNotExist(err) {
		fmt.Printf("Unable to find `%s` in the model directory.\n", defaultModelFile)
		fmt.Printf("Assuming there is a Hugging Face model to convert...\n")
		err = huggingface.NewConverter(app.repo, app.model).Convert()
		if err != nil {
			return err
		}
	}

	return nil
}
