# Gator CLI

A command-line RSS aggregator and user management tool written in Go. It uses PostgreSQL for storage and supports user registration, login, feed management, and following/unfollowing feeds.

## Features
- User registration and login
- Add, follow, unfollow RSS feeds
- Browse and aggregate posts
- List users and feeds
- Reset database rows (for development/testing)
- Colorful terminal output

## What needs to be done
- Refactor so context are passed down the call stack, not just made everytime you call a function... (altough this format was required by this assignment, so not my fault :P)
- Maybe implement concurrent feed fetching?
- Maybe comeback and redo the whole thing, who knows...

## Requirements
- Go 1.20+
- PostgreSQL

## Setup

### 1. Clone the repository
```bash
git clone <your-repo-url>
cd gator-cli
```

### 2. Configure the database
- Create a PostgreSQL database.
- Enable the `uuid-ossp` extension:
  ```sql
  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
  ```
- Set your database URL in `~/.gatorconfig.json`:
  ```json
  {
    "db_url": "postgres://user:password@localhost:5432/yourdb?sslmode=disable",
    "current_user_name": ""
  }
  ```

### 3. Run Migrations
Install [Goose](https://github.com/pressly/goose) and run:
```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### 4. Generate Go code from SQL (optional, if you change queries)
Install [sqlc](https://sqlc.dev):
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
Then run:
```bash
sqlc generate
```

### 5. Build and Run
```bash
go build -o gator-cli
./gator-cli <command> [args]
```

## Usage

### Commands
- `register <username>`: Register a new user
- `login <username>`: Set the current user
- `users`: List all users
- `addfeed <name> <url>`: Add a new RSS feed (must be logged in)
- `feeds`: List all feeds
- `follow <feed_url>`: Follow a feed (must be logged in)
- `unfollow <feed_url>`: Unfollow a feed (must be logged in)
- `following`: List feeds you follow (must be logged in)
- `browse [limit]`: Browse posts (must be logged in; default limit is 2)
- `agg <duration>`: Aggregate feeds every duration (e.g., `1m`, `10s`)
- `reset`: Reset all users and feeds (dangerous!)

### Example
```bash
./gator-cli register alice
./gator-cli login alice
./gator-cli addfeed "Go Blog" https://blog.golang.org/feed.atom
./gator-cli follow https://blog.golang.org/feed.atom
./gator-cli browse 5
```

## Development
- Handlers are in `internal/app/handlers/`
- Middleware (e.g., login required) is in `internal/app/middleware/`
- Database access is in `internal/database/`
- Config is in `internal/config/`

## License
MIT
