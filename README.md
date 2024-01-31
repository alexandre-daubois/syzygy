Syzygy
======

Syzygy is a simple, lightweight, and fast process manager for Unix-like systems.

This program is written in Go and is inspired by Sueprvisor and PM2.

## Installation

### From Source

```bash
$ git clone
$ cd syzygy
$ go build .
```

## Usage

```bash
$ szg configuration_file.yaml
```

This will read the configuration file and start the processes defined in it.

### Interacting with Syzygy

Syzygy can be interacted with using the `szg` command without any arguments.
You'll then be prompted a command line interface to interact with Syzygy.

The following commands are available:

- `help`: Display the help message
- `list`: List the processes managed by Syzygy
- `start <configuration-file>`: Start a process
- `exit`: Exit the command line interface

## Configuration

Syzygy uses a YAML configuration file to define the processes to manage.

```yaml
logs: "/var/log/syzygy.log"
processes:
    process1:
        command: "echo hello world"
        cwd: "/tmp"
        env:
          - "FOO=bar"
        stop_signal: "SIGKILL"                      # SIGINT or SIGKILL
        restart: never                              # never, always
    process2:
        command: "ls -alh"
        restart: never
    sleepy:
        command: "sleep 3"
        restart: always
```
