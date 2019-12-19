# Pricewatcher API
This is the API for the price watching application, workers can be found [here](https://github.com/laetificat/pricewatcher-worker), 
front-end can be found [here](https://github.com/laetificat/pricewatcher-web).

This project stores a history of prices for given products using API endpoints. You can manually create and update the 
prices yourself or use a [worker](https://github.com/laetificat/pricewatcher-worker) to do it for you.

## Installing
You can download binaries from the [releases page](https://github.com/laetificat/pricewatcher-api/releases), you can also
clone this project and run it directly or build a binary with `make build`.

## Running
### add
You can add a new price/watcher by running `pricewatcher add https://yoururlhere`, the following flags are supported:
```text
    --domain string   define the domain, for example: bol.com, ebay.nl, coolblue.nl, etc
-h, --help            help for add
```

### list domains
You can list all the supported domains by running `pricewatcher list domains`, the following flags are supported:
```text
-h, --help            help for list
```

### list watchers
You can list all the watcher by running `pricewatcher list watchers`, the following flags are supported:
```text
-h, --help            help for list
```

### remove
You can remove a watcher by ID by running `pricewatcher remove 1`, or remote them all by running `pricewatcher remove --all`, 
the following flags are supported:
```text
-a, --all    remove all the watchers
-h, --help   help for remove
```

### watch
You can start all the watchers by running `pricewatcher watch`, the following flags are available:
```text
-h, --help               help for watch
-t, --timeout duration   the amount of minutes to wait before checking (default 10m0s)
```

### webserver
You can start the webserver by running `pricewatcher webserver`, `webserver` supports the following flags:
```text
-a, --address string   the ip address to bind the webserver on (default "http://localhost:8080")
-h, --help             help for webserver
```

## Example configuration
```toml
# The database file to use/create.
database_file = "watchers.db"

[log]
    # The minimum log level.
    minimum_level = "info"

    # Sentry configuration
    [log.sentry]
        # Enable Sentry logging.
        enabled = true
        # The DSN you get from Sentry.
        dsn = "https://yoursentry@sentry.io/link"

    # Loggly configuration.
    [log.loggly]
        # Enable Loggly logging.
        enabled = true
        # The token you get from Loggly.
        token = "yourtokenhere"

[notification]
    # Email notification settings.
    [notification.email]
        # Enable email notifications.
        enabled = true
        # List of addresses to send a notification Email.
        addresses = [ "k.heruer@gmail.com" ]
        # The template to use for the email notifications.
        template = "templates/notification.htm"

[webserver]
    # The address for the webserver to listen on.
    address = "http://localhost:8080"

[watcher]
    # Timeout in minutes for the watchers to run their checks.
    timeout = 10
    # The amount of hours the price timestamp should be in the past before adding it to the queue.
    check_interval = 24
```

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md)

## License
See [LICENSE.md](LICENSE.md)