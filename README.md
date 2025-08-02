# Heirloom Deployment API

A Go backend API for managing application deployments across different environments and regions.

## Features

- Release new versions of applications
- Rollback to previous versions
- View active deployments
- View deployment history

## API Endpoints

### Release a new version

```
POST /release
```

Request body:
```json
{
  "application": "auth-service",
  "environment": "production",
  "region": "us-east-1",
  "version": "v2.4.1",
  "deployed_by": "user@example.com"
}
```

### Rollback to a previous version

```
POST /rollback
```

Request body:
```json
{
  "application": "auth-service",
  "environment": "production",
  "region": "us-east-1",
  "version": "v2.3.0",  // Optional, if not provided, will rollback to the previous version
  "deployed_by": "user@example.com"
}
```

### Get active deployments

```
GET /deployments?application=auth-service
```

### Get deployment history

```
GET /history?application=auth-service&environment=production&region=us-east-1
```

## Setup

1. Create a PostgreSQL database using the schema in `src/data/schema.sql`
2. Configure the database connection in the command line flags or use the defaults

## Running the API

```bash
go run main.go --port=8080 --db-host=localhost --db-port=5432 --db-user=postgres --db-pass=postgres --db-name=heirloom
```

## Environment Variables

The API can be configured using command line flags:

- `--port`: Server port (default: 8080)
- `--db-host`: Database host (default: localhost)
- `--db-port`: Database port (default: 5432)
- `--db-user`: Database user (default: postgres)
- `--db-pass`: Database password (default: postgres)
- `--db-name`: Database name (default: heirloom)
- `--db-sslmode`: Database SSL mode (default: disable)