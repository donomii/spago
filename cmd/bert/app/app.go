// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"github.com/urfave/cli/v2"
)

const (
	programName = "bert-server"
)

// BertApp contains everything needed to run the BERT demo client or server.
type BertApp struct {
	*cli.App
	address               string
	grpcAddress           string
	output                string
	model                 string
	repo                  string
	requestText           string
	requestText2          string
	passage               string
	question              string
	serverTimeoutSeconds  int
	serverMaxRequestBytes int
}

// NewBertApp returns BertApp objects. The app can be used as both a client and a server.
func NewBertApp() *BertApp {
	app := &BertApp{
		App: cli.NewApp(),
	}
	app.Name = programName
	app.HelpName = programName
	app.Usage = "A demo for question-answering based on BERT."
	app.Commands = []*cli.Command{
		newClientCommandFor(app),
		newServerCommandFor(app),
	}
	return app
}
