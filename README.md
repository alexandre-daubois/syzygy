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
        restart: never                              # never, always, unless-stopped
    process2:
        command: "ls -alh"
        restart: never
    sleepy:
        command: "sleep 3"
        restart: unless-stopped
```

### Logs

Syzygy can log its activity to a file. This is useful to debug issues.
You can define the path to the log file using the `logs` property.

### Processes

Processes are defined in the `processes` section of the configuration file.

Each process is defined by a name and a set of properties. The name is used to
identify the process in the logs and in the command line interface.

You can change the working directory of the process using the `cwd` property, add
environment variables using the `env` property, and change the signal used to
stop the process using the `stop_signal` property.

### Restart Policies

Syzygy supports the following restart policies:

- `never`: The process will never be restarted
- `always`: The process will always be restarted
- `unless-stopped`: The process will be restarted unless it has been stopped by the user
