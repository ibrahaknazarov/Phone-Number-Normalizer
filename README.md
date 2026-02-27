# Phone Number Normalizer

A Go application that normalizes phone numbers stored in a PostgreSQL database by removing formatting characters and consolidating duplicates.

## Features

- Connects to PostgreSQL database
- Normalizes phone numbers (removes dashes, spaces, parentheses)
- Detects and handles duplicate normalized numbers
- Seeds test data automatically
- Environment-based configuration
- Docker and Docker Compose support
- Comprehensive test coverage

## Prerequisites

### Local Development
- Go 1.21 or higher
- PostgreSQL running on localhost:5432

### Docker
- Docker and Docker Compose

## Configuration

The application reads configuration from environment variables. Copy `.env.example` to `.env` and update with your settings:

```bash
cp .env.example .env
```

### Environment Variables
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: phone)
- `DB_SSLMODE` - SSL mode (default: disable)

## Installation & Usage

### Local Development

```bash
# Install dependencies
go mod download

# Run tests
go test -v

# Build
make build

# Run
make run
```

### Using Make Commands

```bash
make help      # Show available commands
make build     # Build the application
make test      # Run tests
make run       # Build and run
make fmt       # Format code
make vet       # Analyze code
make clean     # Clean build artifacts
```

### Docker Compose (Recommended)

The easiest way to run with a fresh PostgreSQL database:

```bash
docker-compose up --build
```

This will:
1. Start a PostgreSQL 15 container
2. Build the application image
3. Run the application with proper database connection

To stop:
```bash
docker-compose down
```

### Docker Standalone

```bash
# Build the image
docker build -t phone-normalizer .

# Run with your PostgreSQL instance
docker run \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=phone \
  phone-normalizer
```

## Phone Number Formats Supported

The normalizer handles various formats:
- `1234567890`
- `123 456 7891`
- `(123) 456 7892`
- `(123) 456-7893`
- `123-456-7894`
- `(123)456-7892`
- `+11234567890`
- `1-800-FLOWERS` (removes non-digit characters)

All formats are normalized to digits only: `1234567890`

## Testing

Run tests with:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Project Structure

```
.
├── main.go              # Application entry point
├── main_test.go         # Tests for normalize function
├── config/
│   └── config.go        # Configuration management
├── db/
│   └── phone.go         # Database operations
├── Makefile             # Build automation
├── Dockerfile           # Container image definition
├── docker-compose.yml   # Local development setup
├── .env.example         # Example environment variables
├── go.mod               # Module definition
└── README.md            # This file
```

## Development Workflow

1. Make changes to code
2. Run `make fmt` to format
3. Run `make vet` to analyze
4. Run `make test` to test
5. Run `make run` to test locally

## Database Schema

The application creates a simple phone_numbers table:

```sql
CREATE TABLE phone_numbers (
  id SERIAL PRIMARY KEY,
  value VARCHAR(255)
);
```

## License

MIT

