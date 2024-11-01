[![codecov](https://codecov.io/github/mfenerich/kubectl-lr/graph/badge.svg?token=A02R6FB3CV)](https://codecov.io/github/mfenerich/kubectl-lr) [![Go CI/CD Pipeline](https://github.com/mfenerich/kubectl-lr/actions/workflows/go.yml/badge.svg)](https://github.com/mfenerich/kubectl-lr/actions/workflows/go.yml)

# Kubernetes LimitRange CLI Plugin

A `kubectl` plugin for creating `LimitRange` resources with specified CPU and memory limits in Kubernetes namespaces.

## Overview

This plugin provides a command-line interface for easily creating `LimitRange` resources in your Kubernetes cluster. It supports both client-side and server-side dry runs and outputs in YAML or JSON formats for previewing the resource before creation.

## Features

- Create `LimitRange` resources with configurable CPU and memory limits.
- Supports dry-run modes (`client` and `server`) to preview the resource without applying it.
- Outputs resource definitions in YAML or JSON format.
- Easy to use with intuitive command flags.

## Installation

Ensure you have Go installed and `kubectl` configured on your system. Clone this repository and run:

```bash
go build cmd/kubectl-limitrange/kubectl-lr.go
```

Move the binary to a directory in your `PATH`:

```bash
mv kubectl-lr /usr/local/bin/
```

## Run tests

```bash
go test ./pkg/cmd/... -v
```

### Usage

#### Basic Example

```bash
kubectl lr my-limitrange --namespace=my-namespace --max-cpu="1" --min-cpu=100m --max-memory=500Mi --min-memory=100Mi
```

#### Client-Side Dry Run

```bash
kubectl lr my-limitrange --namespace=my-namespace --max-cpu="2" --dry-run=client -o yaml
```

#### Server-Side Dry Run

```bash
kubectl lr my-limitrange --namespace=my-namespace --max-cpu="1" --dry-run=server -o json
```

### Command Flags

- `--max-cpu`: Maximum CPU limit for containers.
- `--min-cpu`: Minimum CPU limit for containers.
- `--default-cpu`: Default CPU limit for containers.
- `--default-request-cpu`: Default CPU request for containers.
- `--max-memory`: Maximum memory limit for containers.
- `--min-memory`: Minimum memory limit for containers.
- `-n, --namespace`: Namespace for the `LimitRange` resource (shorthand for `--namespace`).
- `--dry-run`: Dry-run mode (`client` or `server`).
- `-o, --output`: Output format (`yaml` or `json`).

### Example Commands

- Create a `LimitRange` with CPU and memory limits:
  ```bash
  kubectl lr my-limitrange --namespace=my-namespace --max-cpu="1" --min-cpu=100m --max-memory=500Mi --min-memory=100Mi
  ```

- Client-side dry-run:
  ```bash
  kubectl lr my-limitrange --namespace=my-namespace --max-cpu="2" --min-cpu=500m --dry-run=client -o yaml
  ```

- Server-side dry-run:
  ```bash
  kubectl lr my-limitrange --namespace=my-namespace --default-cpu=500m --default-request-cpu=200m --dry-run=server -o json
  ```

## Requirements

- Go 1.23 or later.
- `kubectl` configured on your system.

## Contributing

Feel free to contribute by submitting issues or pull requests. Any help to enhance the functionality or add new features is welcome!