# HQ

HQ is a simplistic job queue engine communicated by HTTP messages.

## Table of Contents

  - [Installation](#installation)
  - [Usage](#usage)
  - [Configuration](#configuration)
  - [HTTP API](#http-api)
  - [Commands](#commands)
  - [Author](#author)
  - [License](#license)

## Installation

[Download latest version](https://github.com/kohkimakimoto/hq/releases/latest)

## Usage

It is very easy to get started on HQ. After installing HQ, run `hq serve` in your terminal.

```
$ hq serve
2019-04-18T18:56:25+09:00 WARN Your 'data_dir' configuration is not set. HQ server uses a temporary directory that is deleted after the process terminates.
2019-04-18T18:56:25+09:00 WARN Created temporary data directory: /var/folders/7q/7yfsnkpj09n6s1pvlktkpp6h0000gn/T/hq_data_718224193
2019-04-18T18:56:25+09:00 INFO Opened data directory: /var/folders/7q/7yfsnkpj09n6s1pvlktkpp6h0000gn/T/hq_data_718224193
2019-04-18T18:56:25+09:00 INFO Opened boltdb: /var/folders/7q/7yfsnkpj09n6s1pvlktkpp6h0000gn/T/hq_data_718224193/server.bolt
2019-04-18T18:56:25+09:00 INFO The server Listening on 0.0.0.0:19900 (pid: 74090)
```

## Configuration

TODO: Write document


Example:

```toml
# server_id
server_id = 0

# log_level
log_level = "info"

# addr
addr = "0.0.0.0:19900"

# data_dir
data_dir = "/var/lib/hq"

# access_log_file
access_log_file = "/var/log/hq/access.log"

# queues
queues = 8192

# dispatchers
#dispatchers = 1

# max_workers
max_workers = 0

# shutdown_timeout
shutdown_timeout = 10

# job_lifetime
# job_lifetime = 2419200

# job_list_default_limit
job_list_default_limit = 0
```

## HTTP API

HQ core functions are provided via RESTful HTTP API.

Overview of endpoints:

 - [`GET /`](#get-): Gets HQ info.
 - [`GET /stats`](#get-stats): Gets the HQ server statistics.
 - [`POST /job`](#post-job): Pushes a new job.
 - [`GET /job`](#get-job): Lists jobs.
 - [`GET /job/{id}`](#get-jobid): Gets a job.
 - [`DELETE /job/{id}`](#delete-jobid): Deletes a job.
 - [`POST /job/{id}/restart`](#post-jobidrestart): Restarts a job.
 - [`POST /job/{id}/stop`](#post-jobidstop): Stops a job.

### `GET /`

Gets HQ info.

```http
GET /
```

#### Response

```json
{
  "version": "0.3.0",
  "commitHash": "6fe8ba18835f531e16166180feb7335c519df662"
}
```

### `GET /stats`

Gets the HQ server statistics.

### `POST /job`

Pushes a new job.

### `GET /job`

Lists jobs.

### `GET /job/{id}`

Gets a job.

### `DELETE /job/{id}`

Deletes a job.

### `POST /job/{id}/restart`

Restarts a job.

### `POST /job/{id}/stop`

Stops a job.

## Commands

HQ also provides command-line interface to communicate HQ server. To view a list of the available commands, just run `hq` without any arguments:

```
Usage: hq [<options...>] <command>

Simplistic job queue engine
version 0.3.0 (237ea6640ff100150fc9202a1a78322b321cddff)

Options:
  --help, -h     show help
  --version, -v  print the version
  
Commands:
  delete   Deletes a job
  info     Displays a job detail
  list     Lists jobs
  push     Pushes a new job.
  restart  Restarts a job
  serve    Starts the HQ server process
  stats    Displays the HQ server statistics.
  stop     Stops a job
  help, h  Shows a list of commands or help for one command
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
