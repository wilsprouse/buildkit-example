# buildkit-example

A simple Go program that builds a "Hello World" Docker image using BuildKit and stores it in containerd.

## Prerequisites

On Ubuntu, you need:
- Go 1.21 or later
- BuildKit daemon running
- containerd running

### Installing Dependencies

```bash
# Install Go (if not already installed)
sudo apt update
sudo apt install -y golang-go

# Install BuildKit
sudo apt install -y buildkit

# Install containerd (usually comes with Docker or can be installed separately)
sudo apt install -y containerd

# Start buildkit daemon
sudo systemctl start buildkit
sudo systemctl enable buildkit

# Start containerd
sudo systemctl start containerd
sudo systemctl enable containerd
```

## Building the Project

```bash
go build -o buildkit-example main.go
```

## Running

```bash
sudo ./buildkit-example
```

This will:
1. Connect to the BuildKit daemon
2. Build a simple alpine-based "Hello World" image
3. Store it in containerd's image store (in the "buildkit" namespace)

## Verifying the Image

After running, you can verify the image was created:

```bash
# List images in the buildkit namespace
sudo ctr -n buildkit images list | grep hello-world

# Run the image
sudo ctr -n buildkit run --rm hello-world:latest test-container
```

## How It Works

The program uses BuildKit's Low-Level Build (LLB) API to:
1. Start with an Alpine Linux base image
2. Run a simple echo command during build
3. Export the result as an image to containerd

BuildKit uses containerd as its content store backend, so images built with BuildKit are automatically available in containerd without needing a separate push step.
