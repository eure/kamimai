# Kamimai - 紙舞

A migration manager written in Golang. Use it in run commands via the CLI.

[![GoDoc](https://godoc.org/github.com/eure/kamimai?status.svg)](https://godoc.org/github.com/eure/kamimai)
[![wercker status](https://app.wercker.com/status/14ef1b50d66855fbe9c5f17e6e60edd8/s/master "wercker status")](https://app.wercker.com/project/bykey/14ef1b50d66855fbe9c5f17e6e60edd8)

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
kamimai -path=./example/mysql -env=test1 create migrate_name
```

### Up

```shell
# apply all available migrations
kamimai -path=./example/mysql -env=test1 up

# apply the next n migrations
kamimai -path=./example/mysql -env=test1 up n

# apply the given version migration
kamimai -path=./example/mysql -env=test1 up -version=20060102150405
```

### Down

```shell
# rollback the previous migration
kamimai -path=./example/mysql -env=test1 down

# rollback the previous n migrations
kamimai -path=./example/mysql -env=test1 down n

# rollback the given version migration
kamimai -path=./example/mysql -env=test1 down -version=20060102150405
```

### Sync

```shell
# sync all migrations
kamimai -path=./example/mysql -env=test1 sync
```

## Usage in Go code 

_T.B.D._

## Drivers

### Availables

- MySQL

### Plan

- SQLite
- PostgreSQL
- _and more_

## License

[The MIT License (MIT)](http://kaneshin.mit-license.org/)


## Author

Shintaro Kaneko <kaneshin0120@gmail.com>

