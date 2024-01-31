Syzygy
======

Syzygy is a simple, lightweight, and fast process manager for Unix-like systems.

Syzygy is written in Go and is inspired by Sueprvisor and PM2.

## Installation

### From Source

```bash
$ git clone
$ cd szygy
$ go build
```

## Usage

```bash
$ szg configuration_file.yaml
```

## Configuration

Szygy uses a YAML configuration file to define the processes to manage.

```yaml
processes:
    process1:
        command: "echo hello world"
        cwd: "/tmp"
        env:
          - "FOO=bar"
        stop_signal: "SIGTERM"
        restart: never
        output_log_file: "/tmp/process1.out.log"
        events_log_file: "/tmp/process1.events.log"
    process2:
        command: "ls -alh"
        restart: never
    sleepy:
        command: "sleep 3"
        restart: always
```