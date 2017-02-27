<p align="center">
    <img src="https://raw.githubusercontent.com/SeerUK/tid/master/logo.png" height="256" />
</p>

# tid

A simple, CLI-based time tracking utility for personal time tracking. `tid` uses [Bolt][1] as a
storage backend, and does not use any kind of active daemon.

## Installation

```
$ go install github.com/SeerUK/tid/...
```

## Usage

Full usage information can be seen by using the built-in contextual help. This can accessed by 
running:
 
```
$ tid --help
$ # Or for command help:
$ tid start --help
```

**Starting a timer:**

```
$ tid start "A note"
```

**Stopping a timer:**

```
$ tid stop
```

## License

MIT

[1]: https://github.com/boltdb/bolt
