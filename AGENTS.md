# AGENTS.md

> **Purpose**: This file provides context and instructions for AI coding agents working on this repository.

## 1. Project Overview
- **Type**: a monorepo containing the Interactive Brokers client and related API definitions for other services to connect to IBKR.
- **Core Stack**: Golang, Protobuf, ConnectRPC
- **Language**: Golang
- **Package Manager**: go modules

## 2. Getting Started

### Environment Setup
1.  This project uses `go` and `buf`. Make sure you have them installed.
2.  Copy the `.env.example` file in `ibkr-go` to `.env` and fill in the required values: `cp ibkr-go/.env.example ibkr-go/.env`

## 3. Operational Commands
The agent should use these exact commands to perform tasks. **Note**: All `go` commands should be run from the `ibkr-go` directory.

- **Install Dependencies**: `go mod download`
- **Generate Protobuf files**: `buf generate` (from the `proto` directory)
- **Start Dev Server**: `go run ./cmd/server`
- **Run Tests**: `go test -v -race ./...`
- **Lint/Format**: `golangci-lint run ./... --fix`

## 4. Directory Structure
- `ibkr-client`: The root of the monorepo.
- `ibkr-go`: Go-based Interactive Brokers client. It contains the main application logic, API endpoints, and handles communication with the IBKR API.
- `proto`: Houses the Protobuf definitions for the APIs used in the project. This includes the API definitions and generated code.
- `.github`: Contains GitHub-specific files, including workflow definitions for CI/CD, issue templates, and other repository management configurations.

## 5. CI/CD
This project has the following CI/CD pipelines in place:
- **Linting**: On every pull request to `main`, a GitHub Action runs `golangci-lint` to enforce code style and quality.
- **Testing**: On every pull request to `main`, a GitHub Action runs `go test` to ensure that all tests pass and that code coverage is maintained.

## 6. Coding Standards
Use coding standards as defined in `ibkr-go/.golangci.yaml`

## 7. Testing Guidelines
- **Unit Tests**: Write unit tests for all functions.

## 8. Rules & Boundaries
- **NEVER** commit `.env` files.
- **ALWAYS** run the linter before confirming a task is complete.
- If a file is too long (> 300 lines), suggest refactoring it into smaller modules.
