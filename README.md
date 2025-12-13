# ibkr-client

This repository contains the Interactive Brokers (IBKR) client and related API definitions. It's a monorepo built with Golang, Protobuf, and ConnectRPC.

## Project Structure

-   `ibkr-go`: The main Go application that connects to the IBKR API.
-   `proto`: Protobuf definitions for the APIs.
-   `.github`: GitHub Actions workflows for CI/CD.

## Getting Started

### Prerequisites

-   Go
-   [Buf](https://buf.build/docs/installation)

### Installation

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/majidmvulle/ibkr-client.git
    cd ibkr-client
    ```

2.  **Set up environment variables:**

    Navigate to the `ibkr-go` directory and create a `.env` file from the example:

    ```sh
    cd ibkr-go
    cp .env.example .env
    ```

    Update the `.env` file with your IBKR account details.

3.  **Install Go dependencies:**

    From the `ibkr-go` directory:

    ```sh
    go mod download
    ```

4.  **Generate Protobuf code:**

    From the `proto` directory:

    ```sh
    buf generate
    ```

## Usage

### Running the Server

To start the development server, run the following command from the `ibkr-go` directory:

```sh
go run ./cmd/server
```

### Running Tests

To run the test suite, use the following command from the `ibkr-go` directory:

```sh
go test -v -race ./...
```

## CI/CD

This project uses GitHub Actions for continuous integration and deployment. On every pull request to `main`:

-   **Linting**: `golangci-lint` is run to check for code style and quality.
-   **Testing**: `go test` is run to ensure all tests pass.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
