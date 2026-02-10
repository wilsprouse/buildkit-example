package main

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/progress/progressui"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// Connect to buildkitd
	c, err := client.New(ctx, "unix:///run/buildkit/buildkitd.sock")
	if err != nil {
		return fmt.Errorf("failed to connect to buildkit: %w", err)
	}
	defer c.Close()

	// Define the build using LLB (Low-Level Build)
	// This creates a simple alpine-based image with hello world
	state := llb.Image("alpine:latest")

	// Convert to definition
	def, err := state.Marshal(ctx)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Build options - export to containerd image store
	solveOpt := client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: client.ExporterImage,
				Attrs: map[string]string{
					"name": "hello-world:latest",
				},
			},
		},
	}

	// Display progress
	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err := c.Solve(ctx, def, solveOpt, ch)
		return err
	})

	eg.Go(func() error {
		_, err := progressui.DisplaySolveStatus(ctx, nil, os.Stderr, ch)
		return err
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println("\nImage built successfully: hello-world:latest")
	fmt.Println("Image is stored in containerd (buildkit uses containerd as backend)")
	fmt.Println("\nYou can list it with:")
	fmt.Println("  sudo ctr -n buildkit images list | grep hello-world")
	fmt.Println("\nYou can run it with:")
	fmt.Println("  sudo ctr -n buildkit run --rm hello-world:latest test-container")

	return nil
}
