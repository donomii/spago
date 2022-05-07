// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"github.com/nlpodyssey/spago/pkg/nlp/sequencelabeler"
	"github.com/nlpodyssey/spago/pkg/utils/httputils"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

func newServerCommandFor(app *NERApp) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Run the " + programName + " as gRPC/HTTP server.",
		Description: "You must indicate the directory that contains the spaGO neural models.",
		Flags:       newServerCommandFlagsFor(app),
		Action:      newServerCommandActionFor(app),
	}
}

func newServerCommandFlagsFor(app *NERApp) []cli.Flag {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "address",
			Usage:       "Specifies the bind-address of the server.",
			Value:       "0.0.0.0:1987",
			Destination: &app.address,
		},
		&cli.StringFlag{
			Name:        "grpc-address",
			Usage:       "Changes the bind address of the gRPC server.",
			Value:       "0.0.0.0:1976",
			Destination: &app.grpcAddress,
		},
		&cli.StringFlag{
			Name:        "repo",
			Usage:       "Specifies the path to the models.",
			Value:       path.Join(usr.HomeDir, ".spago"),
			Destination: &app.repo,
		},
		&cli.StringFlag{
			Name:        "model",
			Usage:       "Specifies the name of the model to use.",
			Destination: &app.modelName,
			Required:    true,
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

func newServerCommandActionFor(app *NERApp) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		modelsFolder := app.repo
		if _, err := os.Stat(modelsFolder); os.IsNotExist(err) {
			log.Fatal(err)
		}

		modelName := app.modelName
		modelPath := filepath.Join(modelsFolder, modelName)
		if _, err := os.Stat(modelPath); os.IsNotExist(err) {
			switch url, ok := predefinedModels[modelName]; {
			case ok:
				fmt.Printf("Fetch model from `%s`\n", url)
				if err := httputils.DownloadFile(fmt.Sprintf("%s-compressed", modelPath), url); err != nil {
					return err
				}
				r, err := os.Open(fmt.Sprintf("%s-compressed", modelPath))
				if err != nil {
					return err
				}
				fmt.Print("Extracting compressed model... ")
				extractTarGz(r, modelsFolder)
				fmt.Println("ok")
			default:
				return err
			}
		}

		model, err := sequencelabeler.LoadModel(modelPath)
		if err != nil {
			return err
		}
		defer model.Close()

		fmt.Printf("Start HTTP server listening on %s.\n", app.address)

		fmt.Printf("Start gRPC server listening on %s.\n", app.grpcAddress)

		server := sequencelabeler.NewServer(model)
		server.TimeoutSeconds = app.serverTimeoutSeconds
		server.MaxRequestBytes = app.serverMaxRequestBytes
		server.Start(app.address, app.grpcAddress)

		return nil
	}
}
