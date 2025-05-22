# NWC Lightning Wallet API

A Go-based API for making Lightning Network payments between Nostr Wallet Connect (NWC) wallets.

## Features

- Secure API key authentication
- Convert EUR to millisatoshis using current exchange rates
- Make payments between NWC-compatible wallets
- Check wallet health and connectivity
- Swagger UI for easy API testing and documentation
- Docker support for easy deployment

## Prerequisites

- Go 1.24.3 or higher
- Docker and Docker Compose (optional, for containerized deployment)

## Configuration

Create a `.env` file in the root directory with the following variables:

```
# Wallet URIs
WALLET_NAME1="nostr+walletconnect://your-pubkey-here?relay=wss://relay.example.com/v1&secret=your-secret-here"
WALLET_NAME2="nostr+walletconnect://your-pubkey-here?relay=wss://relay.example.com/v1&secret=your-secret-here"

# API Key for authentication
NWC_API_KEY="your-api-key-here"
```

## Installation

### Local Development

```bash
# Clone the repository
git clone <repository-url>
cd nwc_app

# Install dependencies
go mod download

# Generate Swagger documentation
make swagger

# Run the application
make run
```

### Using Docker

```bash
# Build and start with Docker Compose
make docker-up

# Stop the application
make docker-down
```

## API Endpoints

### Health Check

```
GET /health?wallet_id=WALLET_NAME
```

Checks the health of a specific wallet.

### Convert EUR to Millisatoshis

```
GET /convert/eur-to-msats?amount=1&api_key=your-api-key
```

Converts a Euro amount to millisatoshis using the current exchange rate.

### Make a Payment

```
POST /nwc_payment?api_key=your-api-key
```

Request body:
```json
{
  "sender": "WALLET_NAME1", # Wallet URI from the .env file
  "recipient": "WALLET_NAME2", # Wallet URI from the .env file
  "euro_amount": 0.000001 # Amount in EUR
}
```

## Development

### Available Make Commands

- `make run` - Run the application locally
- `make build` - Build the application binary
- `make test` - Run all tests
- `make clean` - Clean build artifacts
- `make docker-build` - Build a Docker image
- `make docker-run` - Run the app in a Docker container
- `make docker-stop` - Stop and remove the Docker container
- `make docker-up` - Start the app using Docker Compose
- `make docker-down` - Stop the app using Docker Compose
- `make swagger` - Generate Swagger documentation
- `make lint` - Lint the code
- `make help` - Display help information

## API Documentation

The API documentation is available through Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## License

MIT
