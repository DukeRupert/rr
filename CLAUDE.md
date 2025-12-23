# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based order reminder service for Rockabilly Roasting, a coffee wholesale business. The application integrates with the Orderspace API for customer/order management and Postmark for email delivery. It runs as a scheduled service that sends automated reminder emails to active customers.

## Development Commands

### Build and Run
```bash
# Build the application
go build -o main ./cmd/main.go

# Run locally
go run cmd/main.go

# Run with Docker Compose (production - uses Postmark)
docker compose up --build

# Run with Docker Compose (development - uses Mailhog)
docker compose -f docker-compose.dev.yml up --build
# View captured emails at http://localhost:8025

# Build Docker image
docker build -t rr .
```

### Testing and Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Run the application (requires .env file)
./main
```

### Environment Setup
Create a `.env` file in the root directory with:
```
ORDERSPACE_CLIENT_ID=your_client_id
ORDERSPACE_CLIENT_SECRET=your_client_secret
POSTMARK_SERVER_TOKEN=your_postmark_token
DATABASE_URL=./rockabilly.db  # Optional, defaults to rockabilly.db

# Optional: For local development with SMTP (e.g., Mailhog)
# SMTP_HOST=localhost
# SMTP_PORT=1025
```

Note: Either `POSTMARK_SERVER_TOKEN` or `SMTP_HOST` must be set. When `SMTP_HOST` is set, emails are sent via SMTP instead of Postmark API.

## Architecture

### Application Structure
The application follows a standard Go project layout with `cmd/` for entrypoints and `internal/` for application code.

**Entry Point**: `cmd/main.go`
- Initializes Echo web server on port 8080
- Sets up database connection (SQLite)
- Configures Postmark email client
- Configures Orderspace API client with OAuth2 token management
- Starts the reminder scheduler service
- Sets up API routes

### Core Components

**Config (`internal/config/`)**
- Loads environment variables from `.env` file using godotenv
- Validates required credentials (Orderspace, Postmark)
- DATABASE_URL defaults to "rockabilly.db" if not specified

**Database (`internal/database/`)**
- SQLite database with automatic schema initialization
- Tables: `customers`, `addresses`, `orders`, `order_lines`, `customer_notifications`, `tokens`
- The `tokens` table stores OAuth access tokens from Orderspace
- The `customer_notifications` table controls email preferences (customers can opt out via `email_notify_days` flag)

**Orderspace Client (`internal/orderspace/`)**
- OAuth2 client credentials flow with automatic token refresh
- Tokens are cached in the database and refreshed when expired (25-minute validity window)
- `Client.GetValidToken()` automatically handles token lifecycle
- `Client.MakeAuthenticatedRequest()` abstracts authenticated API calls
- Supports customer and order data synchronization with Orderspace API

**Email Client (`internal/email/`)**
- `Sender` interface allows swapping between email backends
- `Client` (Postmark): Production email delivery via Postmark API
- `SMTPClient`: Development email delivery via SMTP (for use with Mailhog)
- Email client is selected at startup based on `SMTP_HOST` environment variable

**Reminder Service (`internal/services/`)**
- **Scheduler**: Uses `gocron/v2` to run weekly on Fridays at 10:00 AM MST (America/Denver timezone)
- **SendOrderReminders()**:
  - Fetches customers updated in the last 6 weeks from Orderspace
  - Checks `customer_notifications.email_notify_days` to respect opt-out preferences (defaults to true if no record exists)
  - Sends branded reminder emails asking customers to place orders by Saturday afternoon
  - Logs all send attempts with success/failure status
- **PreviewOrderReminders()**: Sends a preview email to `logan@fireflysoftware.dev` listing all customers who will receive reminders

**API Routes (`internal/api/`)**
- `GET /health` - Health check endpoint, returns `{"status": "ok"}`
- `GET /api/customers` - Fetch customers from Orderspace
- `GET /api/orders` - Fetch orders from Orderspace
- `GET /api/email/preview-reminders` - Trigger preview email showing which customers will receive reminders
- `POST /api/email/send-adhoc` - Send custom ad-hoc emails to all recent customers (for corrections, updates, etc.)
  - Request body: `{"subject": "...", "htmlBody": "...", "textBody": "..."}`
  - Returns: `{"sent": N, "failed": N, "skipped": N, "details": [...]}`

### Data Models (`internal/models/`)
- Customer, Order, OrderLine, Address structures map to both Orderspace API responses and database schema
- EmailAddresses struct has separate fields for orders, dispatches, and invoices
- Status enums with validation methods for Customer and Order statuses

## Key Implementation Details

### Timezone Handling
All scheduled tasks run in Mountain Standard Time (America/Denver). The scheduler is explicitly configured with this location in `services/reminder.go:20`.

### Token Management
The Orderspace client maintains OAuth tokens in the database:
1. Check for existing valid token (< 25 minutes old)
2. If expired or missing, request new token via client credentials flow
3. Store new token with creation timestamp
4. Use token for all authenticated API requests

### Email Notification Control
Customers are included in reminder emails by default. To exclude a customer, add a record to `customer_notifications` with `email_notify_days = false`. The query in `SendOrderReminders()` uses `COALESCE` to default to `true` if no preference record exists.

### Database Schema
- Foreign keys are enabled via `PRAGMA foreign_keys = ON`
- Addresses are linked to customers with ON DELETE CASCADE
- Order lines are linked to orders with ON DELETE CASCADE
- JSON fields (email_addresses, buyers) are stored as TEXT and serialized in Go
- Indexes on frequently queried fields: customer status, order customer_id, order status, delivery_date

## Port Configuration
- Application listens on port 8080
- Docker Compose maps host port 1234 to container port 8080
- Database is persisted in Docker volume `db-data` mounted at `/data`
