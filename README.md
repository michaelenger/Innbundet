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

Add the `--include-example-data` flag to include some RSS feeds:

```shell
innbundet migrate --include-example-data
```

Then you'll want to run the `sync` command to fetch any new feed items. This
should be done periodically:

```shell
innbundet sync
```

Run the web app using the `server` command:

```shell
innbundet server
```

This will serve the application on http://localhost:8080

All the commands will attempt to read the config from a config file, which is
assumed to be `config.yaml` but can also be set using the `--config` parameter:

```shell
innbundet sync --config my_config.yaml
```

### Config File

The configuration file is a YAML and can contain the following:

```yaml
database_file: innbundet.sqlite  # Path to the SQLite database file
description: Tiny RSS reader.    # Description of the page
items_per_page: 25               # Amount of items to show per page
title: Innbundet                 # Title of the page
```

If no file is present it will use all the defaults.

### Adding Feeds

You can add the feeds using the `add` subcommand. Just provide it with the URL
of the feed:

```shell
innbundet add https://michaelenger.com/feed.rss
```

You can also provide a web URL and it will look for a feed in the site's
`<link>` tags:

```shell
innbundet add https://michaelenger.com/
```

### Removing Feeds

You can remove a feed using the `remove` subcommand and the feed's ID:

```shell
innbundet remove 123
```
