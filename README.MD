# Show all available commands
make help

# Complete setup (your workflow)
make setup

# Or step by step:
make create-orders-table  # Create orders_table migration and run it
make run                  # Run the application

# Database operations
make create-migration name=add_email_to_users  # Create new migration
make migrate-up                                 # Run migrations
make migrate-down                               # Rollback last migration
make migrate-status                             # Check migration status
make reset-db                                   # Full database reset

# Development
make deps          # Install dependencies
make build         # Build the binary
make test          # Run tests
make tidy          # Clean up go.mod
make clean         # Remove build artifacts and database

# Hot reload (requires 'air' - auto-installed on first use)
make watch



<!-- create migration file -->
goose -dir ./migrations create orders_table sql

<!-- run migration -->

goose -dir ./migrations sqlite3 ./orders.db up

<!-- run app -->
go run ./main.go
