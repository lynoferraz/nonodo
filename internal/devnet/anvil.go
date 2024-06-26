// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package devnet

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/calindra/nonodo/internal/supervisor"
)

// Default port for the Ethereum node.
const AnvilDefaultAddress = "127.0.0.1"
const AnvilDefaultPort = 8545

// Generate the devnet state and embed it in the Go binary.
//
//go:generate go run ./gen-devnet-state
//go:embed anvil_state.json
var devnetState []byte

const stateFileName = "anvil_state.json"

const anvilCommand = "anvil"

// Start the anvil process in the host machine.
type AnvilWorker struct {
	Address string
	Port    int
	Verbose bool
}

func (w AnvilWorker) String() string {
	return anvilCommand
}

func (w AnvilWorker) Start(ctx context.Context, ready chan<- struct{}) error {
	dir, err := makeStateTemp()
	if err != nil {
		return err
	}
	defer removeTemp(dir)
	slog.Debug("anvil: created temp dir with state file", "dir", dir)

	var server supervisor.ServerWorker
	server.Name = anvilCommand
	server.Command = anvilCommand
	server.Port = w.Port
	server.Args = append(server.Args, "--host", fmt.Sprint(w.Address))
	server.Args = append(server.Args, "--port", fmt.Sprint(w.Port))
	server.Args = append(server.Args, "--load-state", path.Join(dir, stateFileName))
	if !w.Verbose {
		server.Args = append(server.Args, "--silent")
	}
	return server.Start(ctx, ready)
}

// Create a temporary directory with the state file in it.
// The directory should be removed by the callee.
func makeStateTemp() (string, error) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("anvil: failed to create temp dir: %w", err)
	}
	stateFile := path.Join(tempDir, stateFileName)
	const permissions = 0644
	err = os.WriteFile(stateFile, devnetState, permissions)
	if err != nil {
		return "", fmt.Errorf("anvil: failed to write state file: %w", err)
	}
	return tempDir, nil
}

// Delete the temporary directory.
func removeTemp(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		slog.Warn("anvil: failed to remove temp file", "error", err)
	}
}
