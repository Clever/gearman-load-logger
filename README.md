# gearman-load-logger

`gearman-load-logger` queries the `gearman` server and outputs that information to logs for further querying.

## Motivation

[Gearman](http://gearman.org/) as a service does not provide much visibility, and due to its binary protocol can be difficult to easily query. We found ourselves running queries by hand to find out how many workers were running or how many jobs were queued. Instead, we'd like to have a process watch these parameters for us and output to a common format to use for monitoring and alerting elsewhere in our system.

## Features

The script outputs in [kayvee](https://github.com/Clever/kayvee) format:
- total jobs queued (`total_jobs`)
- currently running jobs (`running_jobs`)
- workers available to process jobs, including those currently running jobs (`available_workers`)
- current job load on the workers (`worker_load`)

The worker load is calculated as:
`total_jobs / available_workers`

If there are no workers, a default value of 99 is output as the `worker_load`.

## Running

```bash
./gearman-load-logger --host yourgearmanhost.example.com --port 4730
```

Note that you can usually rely on the default value of 4730 for the gearman port.

## Dependencies

`gearman-load-logger` uses [gearadmin](https://github.com/Clever/gearadmin) and [kayvee](https://github.com/Clever/kayvee), and requires `golang` to be installed (tested with version `1.3`).

However, it is easiest to just download from the [releases page](https://github.com/Clever/gearman-load-logger/releases).
