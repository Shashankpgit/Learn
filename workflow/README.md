## Universal Information Tokenization System

## Installation

Installing the dependencies:

```bash
go mod tidy
```

Set the environment variables:

```bash
cp .env.example .env
```

> **Note:** Open `.env` and modify the environment variables if needed.

## Keycloak Integration

If you want to enable Keycloak authentication, add these environment vars to your `.env` file:

```bash
# Keycloak Configuration
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=finternet
KEYCLOAK_CLIENT_ID=units
KEYCLOAK_CLIENT_SECRET=7AQqLock6kdnoX6ApY0FMjJKeqEOR2pd

# Keycloak Admin Credentials (for user creation)
KEYCLOAK_ADMIN_USER=admin
KEYCLOAK_ADMIN_PASSWORD=admin
```

**Note:** When Keycloak is configured, user registration will automatically create users in Keycloak, and login will use Keycloak's OAuth2 token API.


## Commands

### Running locally:

```bash
make start
```

### Running with live reload:

```bash
air
```

> **Note:** Make sure you have Air installed. See [How to install Air](https://github.com/cosmtrek/air#installation)

### Testing:

```bash
# run all tests
make tests
```

### API Documentation

To view the list of available APIs and their specifications, run the server and go to http://localhost:3000/v1/docs in your browser.