# Backend Stripe Integration

This project demonstrates a backend integration with Stripe using the Gin framework in Go. It includes various functionalities like managing customers, invoices, payment intents, setup intents, and cards. The project is structured to follow best practices in project organization, dependency management, and environment configuration.

## Table of Contents

- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Project Structure](#project-structure)
- [Routes](#routes)
- [Middleware](#middleware)
- [Services](#services)
- [Repositories](#repositories)
- [Models](#models)
- [Running with stripe-mock](#running-with-stripe-mock)
- [Logging](#logging)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

- Go 1.16+
- Docker (optional, for running `stripe-mock` with Docker)
- PostgreSQL (or any other database supported by GORM)

### Installation

1. Clone the repository:
    ```sh
    git clone https://gitlab.com/amcop-saas-platform/vcs/vcs.git
    cd vcs/src/Billing
    ```

2. Install dependencies:
    ```sh
    go mod download
    ```

3. Create a `.env` file in the root of the `backend-stripe` directory with the following content:
    ```env
    STRIPE_API_KEY=your_stripe_api_key
    STRIPE_SECRET_KEY=your_stripe_secret_key
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    PORT=8080
    ```

4. Start the PostgreSQL database and ensure it's running.

5. Run the application:
    ```sh
    go run cmd/main.go
    ```

## Configuration

The configuration is managed using the `viper` package. All configurations are loaded from the `.env` file. The `config/config.go` file handles loading and validation of the configuration.

## Project Structure

- `cmd/main.go`: Entry point of the application.
- `config/`: Contains configuration-related code.
- `internal/api/`: Defines the API routes.
- `internal/database/`: Contains database initialization and model definitions.
- `internal/handlers/`: Handlers for different API endpoints.
- `internal/services/`: Business logic for handling different services.
- `internal/repositories/`: Data access layer for interacting with the database.
- `middleware/`: Middleware functions for the Gin router.
- `providers/`: Initializes and provides external services like Stripe.
- `utils/`: Utility functions like logging and error handling.

## Routes

The routes are defined in the `internal/api/routes.go` file. Below are the primary routes and their functionalities:

- **Card Routes**
  - `POST /create-payment-method`: Create a payment method.
  - `POST /add-card`: Add a card to a customer.
  - `DELETE /delete-card`: Delete a card.
  - `GET /list-cards`: List all cards for a customer.

- **Customer Routes**
  - `POST /register-customer`: Register a new customer.

- **Invoice Routes**
  - `POST /create-invoice`: Create an invoice.
  - `POST /send-invoice`: Send an invoice.
  - `GET /retrieve-invoice/:id`: Retrieve an invoice by ID.
  - `GET /list-customer-invoices/:customer_id`: List all invoices for a customer.
  - `POST /pay-invoice`: Pay an invoice.

- **Payment Routes**
  - `POST /create-payment-intent`: Create a payment intent.
  - `POST /create-setup-intent`: Create a setup intent.
  - `POST /confirm-payment-intent`: Confirm a payment intent.

## Middleware

The middleware functions are defined in the `middleware/middleware.go` file. They include:

- `UserIDMiddleware`: Extracts and validates the user ID from the request header.
- `CORSMiddleware`: Handles Cross-Origin Resource Sharing (CORS).
- `TestMiddleware`: A test middleware to validate the user ID.

## Services

The service layer contains the business logic and interacts with Stripe and the database through repositories. Some of the services include:

- `PaymentService`: Interacts with Stripe for payment operations.
- `InvoiceService`: Handles invoice creation, sending, and retrieval.
- `CardService`: Manages card-related operations.

## Repositories

The repository layer handles data persistence and retrieval. It uses GORM for database operations. Some of the repositories include:

- `CardRepositoryGorm`
- `InvoiceRepositoryGorm`
- `CustomerRepositoryGorm`
- `PaymentIntentRepositoryGorm`
- `SetupIntentRepositoryGorm`

## Models

The models define the structure of the database tables. They are located in `internal/database/models/`. Some of the models include:

- `Card`
- `Customer`
- `Invoice`
- `PaymentIntent`
- `SetupIntent`

## Running with stripe-mock

To test the integration with Stripe without making actual API calls, you can use `stripe-mock`. There are two ways to run `stripe-mock`:

### Using Docker

1. Run stripe-mock using Docker:
    ```sh
    docker run -p 12111-12112:12111-12112 stripe/stripe-mock
    ```

2. Update your `.env` file to point to the local stripe-mock server:
    ```env
    STRIPE_API_KEY=sk_test_4eC39HqLyjWDarjtT1zdp7dc
    STRIPE_SECRET_KEY=sk_test_4eC39HqLyjWDarjtT1zdp7dc
    ```

3. Run your application as usual:
    ```sh
    go run cmd/main.go
    ```

### Using Go Install

1. Install stripe-mock using Go:
    ```sh
    go install github.com/stripe/stripe-mock@latest
    ```

2. Run stripe-mock:
    ```sh
    stripe-mock
    ```

3. Update your `.env` file to point to the local stripe-mock server:
    ```env
    STRIPE_API_KEY=sk_test_4eC39HqLyjWDarjtT1zdp7dc
    STRIPE_SECRET_KEY=sk_test_4eC39HqLyjWDarjtT1zdp7dc
    ```

4. Run your application as usual:
    ```sh
    go run cmd/main.go
    ```

## Logging

Logging is handled using the `zap` package. The logger is initialized in `utils/logger.go`. All logs are written in JSON format with ISO8601 timestamps.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.