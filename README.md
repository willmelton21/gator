Sure! Here's the full `README.md` content formatted in Markdown:

```markdown
# Gator CLI

Gator is a simple command-line tool built with Go that lets users manage feeds, follow other users, and aggregate content from sources. This project uses a PostgreSQL database for data storage and expects a configuration file to manage user sessions.

## Prerequisites

To run this project, you’ll need the following installed on your machine:

- [Go](https://go.dev/dl/) (version 1.18 or higher recommended)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

You can install the Gator CLI using `go install`:

```bash
go install github.com/your-username/gator@latest
```

Make sure your `$GOPATH/bin` is in your `PATH` so you can run the `gator` command from anywhere.

## Configuration

Before using Gator, you'll need to create a configuration file in your home directory named `.gatorconfig.json`. This file stores your database connection info and the currently logged-in user.

### Example `.gatorconfig.json`

```json
{
  "db_url": "postgres://user:password@localhost:5432/gator_db?sslmode=disable",
  "current_user_name": ""
}
```

> ⚠️ Replace `user`, `password`, and `gator_db` with your actual Postgres credentials and database name.

## Running Gator

Once installed and configured, you can run the CLI by typing `gator` followed by any of the supported commands.

### Available Commands

| Command       | Description                            |
|---------------|----------------------------------------|
| `register`    | Register a new user                    |
| `login`       | Log in as an existing user             |
| `browse`      | Browse aggregated content              |
| `follow`      | Follow another user                    |
| `unfollow`    | Unfollow a user                        |
| `following`   | View users you are following           |
| `users`       | List all registered users              |
| `agg`         | Run content aggregation manually       |
| `addfeed`     | Add a new feed URL                     |
| `feeds`       | List your subscribed feeds             |
| `reset`       | Delete your account                    |

Some commands, such as `browse`, `follow`, `addfeed`, etc., require you to be logged in.

## Notes

- Config file is automatically read from your home directory when running any command.
- You can change the current user with the `login` command, which updates your `.gatorconfig.json`.
```
