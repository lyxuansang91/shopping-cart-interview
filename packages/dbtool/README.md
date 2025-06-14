# Database Tool (dbtool)

A CLI tool for managing database operations across all services in the Cinch API.

## Features

- Run migrations up/down for all services or specified services
- Run migrations to specific versions
- Run seeders for all services or specified services

## Installation

```bash
# From the root of the project
go install ./packages/dbtool/cmd/dbtool
```

## Usage

### Running Migrations

Run migrations up to latest version for all services:

```bash
dbtool migrate up --uri "mysql://user:password@tcp(localhost:3306)/database"
```

Run migrations up to a specific version:

```bash
dbtool migrate up 000002 --uri "mysql://user:password@tcp(localhost:3306)/database"
```

Run migrations down to a specific version:

```bash
dbtool migrate down 000001 --uri "mysql://user:password@tcp(localhost:3306)/database"
```

Run migrations for specific services:

```bash
dbtool migrate up --uri "mysql://user:password@tcp(localhost:3306)/database" --services payments,users
dbtool migrate down 000001 --uri "mysql://user:password@tcp(localhost:3306)/database" --services payments
```

### Running Seeders

Run seeders for all services:

```bash
dbtool seed --uri "mysql://user:password@tcp(localhost:3306)/database"
```

Run seeders for specific services:

```bash
dbtool seed --uri "mysql://user:password@tcp(localhost:3306)/database" --services payments,users
```

### Using with Make

The project includes several Make targets for convenience:

```bash
# Run all migrations up to latest version
make dbtool-migrate

# Run migrations down to a specific version
make dbtool-migrate-down revision=000001

# Run all seeders
make dbtool-seed
```

## Directory Structure

The tool expects the following directory structure for each service:

```
services/
  service-name/
    database/
      migrations/
        000001_initial.up.sql
        000001_initial.down.sql
        000002_add_users.up.sql
        000002_add_users.down.sql
      seed.sql
```

## Migration Versioning

Migration files should follow the format:

```
<version>_<description>.(up|down).sql
```

Examples:
- `000001_initial.up.sql`
- `000001_initial.down.sql`
- `000002_add_users.up.sql`
- `000002_add_users.down.sql`

When specifying versions in commands, you can use either the full filename or just the version number:

```bash
dbtool migrate up 000002
# or
dbtool migrate up 000002_add_users.up.sql
```

## Development

To build the tool:

```bash
go build -o dbtool ./packages/dbtool/cmd/dbtool
# or
make dbtool
```
