# HQ

HQ is a simplistic, language agnostic job queue engine communicated by HTTP messages.

HQ is implemented as a standalone JSON over HTTP API server. In the job running workflow, it behaves like an asynchronous HTTP proxy server. When you push a job to the HQ server, it stores the job and sends asynchronous HTTP POST request to a URL that specified in the job.

Worker applications that actually run the jobs are web applications. So you can implement the workers in Any programming language that can talk HTTP.

```
                   ┌────────────────────────────────────────────────────────────────┐
┌───┐              │HQ                                                      ┌──────┐│              ┌──────────┐
│app│──POST /job──>┼┐          ┌───────────────────┐                      ┌>│worker│┼──POST /xxx──>│worker app│
└───┘              ││          │queue              │                      │ └──────┘│              └──────────┘
┌───┐              ││          │┌───┐ ┌───┐   ┌───┐│          ┌──────────┐│ ┌──────┐│              ┌──────────┐
│app│──POST /job──>┼┼─enqueue->││job│ │job│...│job││-dequeue->│dispatcher│┼>│worker│┼──POST /xxx──>│worker app│
└───┘              ││          │└───┘ └───┘   └───┘│          └──────────┘│ └──────┘│              └──────────┘
┌───┐              ││          └───────────────────┘                      │ ┌──────┐│              ┌──────────┐
│app│──POST /job──>┼┘                                                     └>│worker│┼──POST /xxx──>│worker app│
└───┘              │                                                        └──────┘│              └──────────┘
                   └────────────────────────────────────────────────────────────────┘
```

## Table of Contents

  - [Installation](#installation)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
    - [Example](#example)
    - [Parameters](#parameters)
  - [Job](#job)
  - [HTTP API](#http-api)
    - [`GET /`](#get-)
    - [`GET /stats`](#get-stats)
    - [`POST /job`](#post-job)
    - [`GET /job`](#get-job)
    - [`GET /job/{id}`](#get-jobid)
    - [`DELETE /job/{id}`](#delete-jobid)
    - [`POST /job/{id}/restart`](#post-jobidrestart)
    - [`POST /job/{id}/stop`](#post-jobidstop)
  - [Commands](#commands)
  - [Author](#author)
  - [License](#license)

## Installation

HQ lives in a single binary without external dependencies. You can download a precompiled binary at the Github releases page.

[Download latest version](https://github.com/kohkimakimoto/hq/releases/latest)

If you use CentOS7, you can also use RPM package that is stored in the same releases page. It is useful because it configures systemd service automatically.

## Getting Started

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

The config file must be written in [TOML](https://github.com/toml-lang/toml). You can specify the config file by `-c` or `-config-file` option when HQ runs like the following.

```
$ hq serve -c /path/to/config.toml
```

### Example

```toml
# server_id
server_id = 0

# addr
addr = "0.0.0.0:19900"

# data_dir
data_dir = "/var/lib/hq"

# log_level
log_level = "info"

# log_file
log_file = "/var/log/hq/hq.log"

# access_log_file
access_log_file = "/var/log/hq/access.log"

# queues
queues = 8192

# dispatchers
dispatchers = 1

# max_workers
max_workers = 0

# shutdown_timeout
shutdown_timeout = 10

# job_lifetime
job_lifetime = 2419200

# job_list_default_limit
job_list_default_limit = 0
```

### Parameters

* `server_id` (number): This is used to generate Job ID. HQ uses [go-katsubushi](https://github.com/kayac/go-katsubushi) to allocate unique ID. If you want Job ID to be unique on mulitple servers, You need to set `server_id` unique on each servers. The default is `0`.

* `addr` (string): The listen address to the HQ server process. The default is `0.0.0.0:19900`.

* `data_dir` (string): The data directory to store all generated data by the HQ sever. You should set the parameter to keep jobs persistantly. If you doesn't set it, HQ uses a temporary directory that is deleted after the process terminates.

* `log_level` (string): The log level (`debug|info|warn|error`). The default is `info`.

* `log_file` (string):

* `access_log_file` (string):

* `queues` (number):

* `dispatchers` (number):

* `max_workers` (number):

* `shutdown_timeout` (number):

* `job_lifetime` (number):

* `job_list_default_limit` (number):


## Job

Job in HQ is a JSON like the following:

```json
{
  "canceled": false,
  "comment": "This is an example job!",
  "createdAt": "2019-10-29T07:32:26.054Z",
  "err": "",
  "failure": false,
  "finishedAt": "2019-10-29T07:32:28.548Z",
  "headers": null,
  "id": "109192606348480512",
  "name": "example-job",
  "output": "OK",
  "payload": {
    "message": "Hello world!"
  },
  "running": false,
  "startedAt": "2019-10-29T07:32:28.252Z",
  "status": "success",
  "statusCode": 200,
  "success": true,
  "timeout": 0,
  "url": "http://your-worker-server/worker/example",
  "waiting": false
}
```

To create a new job, You can use [`POST /job`](#post-job) API.

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

By default, the output of all HTTP API requests is minimized JSON. If the client passes `pretty` on the query string, formatted JSON will be returned.

### `GET /`

#### Request

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

#### Request

```http
GET /stats
```

#### Response


```json
{
  "version": "0.3.0",
  "commitHash": "6fe8ba18835f531e16166180feb7335c519df662",
  "serverId": 0,
  "queues": 8192,
  "dispatchers": 8,
  "maxWorkers": 0,
  "shutdownTimeout": 10,
  "jobLifetime": 2419200,
  "jobLifetimeStr": "4 weeks",
  "jobListDefaultLimit": 0,
  "queueMax": 8192,
  "queueUsage": 0,
  "numWaitingJobs": 0,
  "numRunningJobs": 0,
  "numWorkers": 0,
  "numJobs": 0
}
```

### `POST /job`

Pushes a new job.

#### Request

```http
POST /job
```

```json
{
  "name": "example",
  "comment": "This is an example job!",
  "url": "https://your-worker-serevr/",
  "payload": {
    "foo": "bar"
  },
  "header": {
    "X-Custom-Token": "xxxxxxx"
  },
  "timeout": 0
}
```
- `url` (string,required):
- `name` (string):
- `comment` (string):
- `payload` (json):
- `header` (json):
- `timeout` (number):

#### Response

```json
{
  "canceled": false,
  "comment": "",
  "createdAt": "2019-10-29T23:57:08.713Z",
  "err": "",
  "failure": false,
  "finishedAt": null,
  "headers": null,
  "id": "109440416981450752",
  "name": "default",
  "output": "",
  "payload": {
    "foo": "bar"
  },
  "running": false,
  "startedAt": null,
  "status": "unfinished",
  "statusCode": null,
  "success": false,
  "timeout": 0,
  "url": "https://your-worker-serevr/",
  "waiting": false
}
```

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
