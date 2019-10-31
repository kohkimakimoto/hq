# HQ

HQ is a simplistic, language agnostic job queue engine communicated by HTTP messages.

HQ is implemented as a standalone JSON over HTTP API server. In the job running workflow, it behaves like an asynchronous HTTP proxy server. When you push a job to the HQ server, it stores the job in the internal queue database and sends asynchronous HTTP POST request to a URL that specified in the job.

Worker applications that actually run the jobs are web applications. So you can implement applications for the jobs in Any programming language that can talk HTTP (such as Go, PHP, Python, etc).

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

> Note: Running HQ server without any configuration like the above can cause to lost queued jobs, because HQ uses temporary directory to store jobs. Therefore this should be used only on DEV environment. When you use HQ on your production environment, You should set a proper configuration file. see [Configuration](#configuration).

You can push a job by using the following `curl` command:

```
$ curl -XPOST http://localhost:19900/job -H "Content-Type: application/json" -d '{"url": "http://your-worker-app-server/example", "payload": {"message": "Hello world!"}}'
```

## Configuration

The config file must be written in [TOML](https://github.com/toml-lang/toml). You can specify the config file by `-c` or `-config-file` option when HQ runs like the following.

```
$ hq serve -c /path/to/config.toml
```

### Example

```toml
server_id = 0
addr = "0.0.0.0:19900"
data_dir = "/var/lib/hq"
log_level = "info"
log_file = "/var/log/hq/hq.log"
access_log_file = "/var/log/hq/access.log"
queues = 8192
dispatchers = 1
max_workers = 0
shutdown_timeout = 10
job_lifetime = 2419200
job_list_default_limit = 0
```

### Parameters

* `server_id` (number): This is used to generate Job ID. HQ uses [go-katsubushi](https://github.com/kayac/go-katsubushi) to allocate unique ID. If you want Job ID to be unique on mulitple servers, You need to set `server_id` unique on each servers. The default is `0`.

* `addr` (string): The listen address to the HQ server process. The default is `0.0.0.0:19900`.

* `data_dir` (string): The data directory to store all generated data by the HQ sever. You should set the parameter to keep jobs persistantly. If you doesn't set it, HQ uses a temporary directory that is deleted after the process terminates.

* `log_level` (string): The log level (`debug|info|warn|error`). The default is `info`.

* `log_file` (string): The log file path. If you do not set, HQ writes log to STDOUT.

* `access_log_file` (string):The access log file path. If you do not set, HQ writes log to STDOUT.

* `queues` (number): Size of queue. The default is `8192`.

* `dispatchers` (number): Number of dispatchers. The default is `runtime.NumCPU()`.

* `max_workers` (number): Number of max workers. The dispacher can execute multiple workers concurrently. This config is limit how many each dispatcher can run workers concurrently. For example, If you set `dispatcher = 2` and `max_workers = 3`, HQ can run max `6` workers at the same time. If you set `max_workers = 0`, each dispacher run only one worker synchronously. The default is `0`.

* `shutdown_timeout` (number): This is time how many seconds HQ waits executing jobs to finish in a shutdown process. If HQ server process receives `SIGINT` or `SIGTERM`, it try to shutdown itself. If HQ has executing workers, it waits workers to finish or number of seconds of this config. The default is `10`.

* `job_lifetime` (number): HQ removes old finished jobs automatically. This config sets time how many seconds HQ keeps jobs. If you set it `0`, HQ does not remove any jobs. The default is `2419200` (28 days).

* `job_list_default_limit` (number): The default `limit` value of [`GET /job`](#get-job). The default is `0` (no limit).

## Job

Job in HQ is a JSON as the following:

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
  "url": "http://your-worker-app-server/example",
  "waiting": false
}
```

To create and push a new job, You can use [`POST /job`](#post-job) API.
The pushed job is stored in the queue and executed by the HQ worker. The HQ worker constructs HTTP POST request from the job. You can customize this request headers and JSON payload by the job properties.

If the above example job is executed, HQ will send like the following HTTP request:

```
POST /example HTTP/1.1
Host: your-worker-app-server
Content-Type: application/json
User-Agent: HQ/1.0.0

{
  "message": "Hello world!"
}
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
  "url": "https://your-worker-app-server/example",
  "payload": {
    "message": "Hello world!"
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
  "comment": "This is an example job!",
  "createdAt": "2019-10-29T23:57:08.713Z",
  "err": "",
  "failure": false,
  "finishedAt": null,
  "headers": null,
  "id": "109440416981450752",
  "name": "example",
  "output": "",
  "payload": {
    "message": "Hello world!"
  },
  "running": false,
  "startedAt": null,
  "status": "unfinished",
  "statusCode": null,
  "success": false,
  "timeout": 0,
  "url": "https://your-worker-app-server/example",
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
