# Kamimai - 紙舞

A migration manager written in Golang. Use it in run commands via the CLI.

[![GoDoc](https://godoc.org/github.com/eure/kamimai?status.svg)](https://godoc.org/github.com/eure/kamimai)
[![Go Report Card](https://goreportcard.com/badge/github.com/eure/kamimai)](https://goreportcard.com/report/github.com/eure/kamimai)

## Installation

`kamimai` is written in Go, so if you have Go installed you can install it with go get:

```shell
go get github.com/eure/kamimai/cmd/kamimai
```

Make sure that `kamimai` was installed correctly:

```shell
kamimai -h
```

## Usage:

### Create

```shell
# create new migration files
kamimai -path=./examples/mysql -env=test1 create migrate_name
```

### Up

```shell
# apply all available migrations
kamimai -path=./examples/mysql -env=test1 up

# apply the next n migrations
kamimai -path=./examples/mysql -env=test1 up n

# apply the given version migration
kamimai -path=./examples/mysql -env=test1 up -version=20060102150405
```

### Down

```shell
# rollback the previous migration
kamimai -path=./examples/mysql -env=test1 down

# rollback the previous n migrations
kamimai -path=./example/mysql -env=test1 down n

# rollback the given version migration
kamimai -path=./examples/mysql -env=test1 down -version=20060102150405
```

### Sync

```shell
# sync all migrations
kamimai -path=./examples/mysql -env=test1 sync
```

## Usage in Go code

```go
package main

import (
	"github.com/eure/kamimai"
	"github.com/eure/kamimai/core"
	_ "github.com/eure/kamimai/driver"
)

func main() {
	conf, err := core.NewConfig("examples/testdata")
	if err != nil {
		panic(err)
	}

	conf.WithEnv("development")

	// Sync
	kamimai.Sync(conf)

	// ...
```

## Drivers

### Availables

- MySQL

### Plan

- SQLite
- PostgreSQL
- _and more_

## Contribution

### Setup

```sh
docker-compose up -d

# kamimai の実行
MYSQL_HOST=127.0.0.1 MYSQL_USER=kamimai MYSQL_PASSWORD=kamimai go run ./cmd/kamimai --dry-run --env=development --path=./examples/mysql create test1
kamimai: created examples/mysql/migrations/002_test_up.sql
kamimai: created examples/mysql/migrations/002_test_down.sql
```

## License

[The MIT License (MIT)](https://github.com/eure/kamimai/blob/master/LICENSE)
