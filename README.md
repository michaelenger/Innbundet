# innbundet

Personal RSS reader, inspired by [reedi](https://github.com/facundoolano/feedi).

## Requirements

* Go

## Usage

Build the application:

```shell
go build .
```

Migrate the database models (this only needs to be done once):

```shell
innbundet migrate
```

Run the web app using the `server` command:

```shell
innbundet server
```

This will serve the application on http://localhost:8080
