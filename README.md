# Phone Number Normalizer

A Go application that normalizes phone numbers stored in a PostgreSQL database by removing formatting characters and consolidating duplicates.

## Features

- Connects to PostgreSQL database
- Normalizes phone numbers (removes dashes, spaces, parentheses)
- Detects and handles duplicate normalized numbers
- Seeds test data automatically

## Prerequisites

- Go 1.21 or higher
- PostgreSQL running on localhost:5432
- Database credentials:
  - User: `postgres`
  - Password: `postgres`

## Usage

```bash
go run main.go
```

The application will:
1. Drop and recreate the `phone` database
2. Create the `phone_numbers` table
3. Seed test data with various phone number formats
4. Normalize all numbers and consolidate duplicates

## Phone Number Formats Supported

The normalizer handles various formats:
- `1234567890`
- `123 456 7891`
- `(123) 456 7892`
- `(123) 456-7893`
- `123-456-7894`
- `(123)456-7892`

All formats are normalized to digits only: `1234567890`

## Testing

```bash
go test -v
```

## License

MIT
