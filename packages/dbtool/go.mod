module github.com/cinchprotocol/cinch-api/packages/dbtool

go 1.24.1

require (
	github.com/cinchprotocol/cinch-api/packages/core v0.0.0
	github.com/go-sql-driver/mysql v1.9.2
	github.com/golang-migrate/migrate/v4 v4.18.2
	github.com/spf13/cobra v1.8.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)

replace github.com/cinchprotocol/cinch-api/packages/core => ../core
