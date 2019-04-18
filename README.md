# HQ

HQ is a simplistic job queue engine.

----

## Installation

As `hq` is a single binary command, to install it you can download the binrary from github release page and drop it in you $PATH

[Download latest version](https://github.com/kohkimakimoto/hq/releases/latest)

If you use CentOS7, you can also use RPM package that is stored in the same release page. It is useful because it configures systemd service automatically.

## Usage

## Configuration

## HTTP API

HQ provides functions as a RESTful HTTP API.

Overview of endpoints:

* `GET /`:  
* `GET /stats`:
* `POST /job`:
* `GET /job`:
* `GET /job/{id}`:
* `DELETE /job/{id}`:
* `POST /job/{id}/restart`:
* `POST /job/{id}/stop`:

## Commands

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
