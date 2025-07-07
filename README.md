# billsplittr

A modern bill splitting application built with Go, designed to help groups of friends easily split expenses and manage shared costs.

## Features

- **User Management**: User registration, authentication, and profile management
- **Friendship System**: Add and manage friends for easy expense sharing
- **Group Expenses**: Create and manage shared expenses with multiple participants
- **Debt Tracking**: Automatic calculation and tracking of debts between users
- **Transfer Methods**: Support for various payment methods
- **RESTful API**: Clean API design for frontend integration

## Tech Stack

- **Backend**: Go 1.24.4
- **Web Framework**: Gin
- **Database**: PostgreSQL/MySQL (configurable)
- **ORM**: GORM
- **Authentication**: JWT tokens
- **Containerization**: Docker
- **Testing**: Go testing framework with testify

## Project Structure

```
billsplittr/
├── cmd/
│   └── api/                    # Application entry point
├── internal/
│   ├── entity/                 # Domain entities
│   ├── dto/                    # Data transfer objects
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic layer
│   ├── delivery/               # HTTP handlers and routing
│   ├── mapper/                 # Entity-DTO mapping
│   ├── provider/               # Dependency injection
│   ├── helper/                 # Utility helpers
│   ├── util/                   # Common utilities
│   ├── appconstant/           # Application constants
│   └── tests/                  # Test files
├── db/                         # Database migrations and seeds
├── bin/                        # Compiled binaries
└── tmp/                        # Temporary files
```

## Getting Started

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL or MySQL database
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/itsLeonB/billsplittr.git
cd billsplittr
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

4. Configure your database settings in `.env`:
```env
SQLDB_HOST=localhost
SQLDB_PORT=5432
SQLDB_USER=postgres
SQLDB_PASSWORD=postgres
SQLDB_NAME=billsplittr
SQLDB_DRIVER=postgres

APP_ENV=debug
APP_PORT=8080
APP_TIMEOUT=5s
APP_CLIENTURLS=http://localhost:5173
APP_TIMEZONE=Asia/Jakarta

AUTH_SECRETKEY=your-secret-key
AUTH_TOKENDURATION=9h
AUTH_COOKIEDURATION=9h
AUTH_ISSUER=billsplittr
```

### Running the Application

#### Development Mode (with hot reload)

Make sure you have [air](https://github.com/cosmtrek/air) installed:
```bash
go install github.com/cosmtrek/air@latest
```

Then run:
```bash
make api-hotreload
```

#### Production Mode

Build and run the application:
```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

#### Using Docker

Build the Docker image:
```bash
docker build -t billsplittr .
```

Run the container:
```bash
docker run -p 8080:8080 --env-file .env billsplittr
```

## Development

### Available Make Commands

- `make api-hotreload` - Run the API with hot reload (requires air)
- `make lint` - Run linting (requires golangci-lint)
- `make test` - Run tests

### Testing

Run the test suite:
```bash
go test ./internal/tests/...
```

Or using make:
```bash
make test
```

### Code Quality

Run linting:
```bash
make lint
```

## API Documentation

The application provides a RESTful API for managing users, expenses, and friendships. The API runs on port 8080 by default.

Key endpoints include:
- User authentication and management
- Friend management
- Group expense creation and management
- Debt calculation and tracking

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## Contact

For questions or support, please open an issue on GitHub.
