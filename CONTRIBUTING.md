# Contributing to Kubernetes LimitRange CLI Plugin

First of all, thank you for taking the time to contribute! Contributions are always welcome, whether it's code, documentation, bug reports, or feature requests.

## Table of Contents

- [Contributing to Kubernetes LimitRange CLI Plugin](#contributing-to-kubernetes-limitrange-cli-plugin)
  - [Table of Contents](#table-of-contents)
  - [Code of Conduct](#code-of-conduct)
  - [How Can I Contribute?](#how-can-i-contribute)
    - [Reporting Bugs](#reporting-bugs)
    - [Suggesting Enhancements](#suggesting-enhancements)
    - [Pull Requests](#pull-requests)
  - [Development Guidelines](#development-guidelines)
    - [Setting Up Your Environment](#setting-up-your-environment)
    - [Running Tests](#running-tests)
    - [Linting and Code Quality](#linting-and-code-quality)
  - [Commit Message Guidelines](#commit-message-guidelines)
  - [Code Review Process](#code-review-process)
  - [Questions](#questions)

## Code of Conduct

Please note that this project is governed by a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report any behavior that violates this code to [marcel@feneri.ch](mailto:marcel@feneri.ch).

## How Can I Contribute?

### Reporting Bugs

If you find a bug, please create an issue and provide the following information:
- A clear and descriptive title for the issue.
- Steps to reproduce the problem.
- The version of Go and the Kubernetes version you are using.
- Any error messages or relevant logs.

### Suggesting Enhancements

We welcome feature requests! If you have a suggestion, please create an issue with:
- A detailed description of the feature.
- Why you think this feature is needed.
- Any potential challenges in implementing the feature.

### Pull Requests

If you are ready to submit your changes:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/my-feature`).
3. Make your changes and ensure your code is well-formatted and documented.
4. Commit your changes (`git commit -m "feat: add my new feature"`).
5. Push to the branch (`git push origin feature/my-feature`).
6. Open a Pull Request.

Please ensure your code passes all tests and linting checks before submission.

## Development Guidelines

### Setting Up Your Environment

Ensure that you have the following:
- Go version ${{ vars.GO_VERSION }} or later.
- `kubectl` configured on your system.

### Running Tests

Run the test suite to make sure your changes don't break anything:
```bash
go test ./pkg/cmd/... -v
```

### Linting and Code Quality

We use `golangci-lint` for linting. Run the following command to check your code:
```bash
golangci-lint run
```

We also use `gosec` for security checks:
```bash
gosec ./...
```

Ensure all linter and security checks pass before committing your code.

## Commit Message Guidelines

Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification:
- **feat**: A new feature.
- **fix**: A bug fix.
- **docs**: Documentation changes.
- **style**: Code style changes (formatting, etc.).
- **refactor**: Code refactoring without changing functionality.
- **test**: Adding or modifying tests.
- **chore**: Other changes that don't modify source or test files.

Example:
```
feat: add support for server-side dry run
```

## Code Review Process

1. Pull Requests will be reviewed by project maintainers.
2. Ensure your branch is up-to-date with `main`.
3. Address any review comments and make necessary changes.
4. Once approved, the PR will be merged by a maintainer.

## Questions

If you have any questions or need assistance, feel free to create an issue or reach out via [marcel@feneri.ch](marcel@feneri.ch).