# gearman-load-logger

`gearman-load-logger` queries the `gearman` server at an interval and outputs that information to logs for further querying.

## Motivation

[Gearman](http://gearman.org/) as a service does not provide much visibility, and due to its binary protocol can be difficult to easily query. We found ourselves running queries by hand to find out how many workers were running or how many jobs were queued. Instead, we'd like to have a process watch these parameters for us and output to a common format to use for monitoring and alerting elsewhere in our system.

## Features

The script outputs in [kayvee](https://github.com/Clever/kayvee) format:
- total jobs queued (`total`)
- currently running jobs (`running`)
- workers available to process jobs, including those currently running jobs (`workers`)

## Running

```bash
./gearman-load-logger --host yourgearmanhost.example.com --port 4730 --interval 1m
```

Note that you can usually rely on the default values of 4730 for the gearman port and 1 minute for the interval.

## Dependencies

`gearman-load-logger` uses [gearadmin](https://github.com/Clever/gearadmin) and [kayvee](https://github.com/Clever/kayvee), and requires `golang` to be installed (tested with version `1.3`).

However, it is easiest to just download from the [releases page](https://github.com/Clever/gearman-load-logger/releases).
