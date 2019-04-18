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

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
