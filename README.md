# HQ

HQ is a simplistic job queue engine.

----


Table of Contents

  - [Installation](#installation)
  - [Usage](#usage)
  - [Configuration](#configuration)
  - [HTTP API](#http-api)
  - [Commands](#commands)
  - [Author](#author)
  - [License](#license)


## Installation

As `hq` is a single binary command, to install it you can download the binrary from github release page and drop it in you $PATH

[Download latest version](https://github.com/kohkimakimoto/hq/releases/latest)

If you use CentOS7, you can also use RPM package that is stored in the same release page. It is useful because it configures systemd service automatically.

## Usage

## Configuration

## HTTP API

HQ provides functions as a RESTful HTTP API.

Overview of endpoints:

* `GET /`: Gets HQ info.
* `GET /stats`: Gets the HQ server statistics.
* `POST /job`: Pushes a new job.
* `GET /job`: Lists jobs.
* `GET /job/{id}`: Gets a job.
* `DELETE /job/{id}`: Deletes a job.
* `POST /job/{id}/restart`: Restarts a job.
* `POST /job/{id}/stop`: Stops a job.

## Commands

HQ provides command-line interface to manipulate HQ jobs. To view a list of the available commands, just run `hq` without any arguments:

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
