# aggreGator: Blog Aggregator in Go

A simple CLI tool for following RSS feeds and storing posts in a local postgres database.
Created as a [boot.dev project](https://www.boot.dev/courses/build-blog-aggregator-golang).

## Requirements

You'll need to install:

- Postgres
- Go

## Installation & Configuration

- Run "go install ." to install the gator CLI tool.
- Create a ".gatorconfig.json" file in your home directory (e.g. "~/.gatorconfig.json") and populate with the following content. You'll need to use a valid connection string for your environment.

```
{
        "db_url": "YOUR_POSTGRES_DB_CONNECTION_STRING_HERE",
        "current_user_name": "PICK_A_USERNAME"
}
```

## How to Use

Type "aggreGATOR help" to see available commands.
Run a command with "aggreGATOR \<command name\>".
