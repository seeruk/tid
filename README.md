<p align="center">
    <img src="https://raw.githubusercontent.com/SeerUK/tid/master/logo.png" height="256" />
</p>

# tid

A simple, CLI-based time tracking utility for personal time tracking. `tid` uses [Bolt][1] as a
storage backend, and does not use any kind of active daemon.

## Installation

```
$ go get -u -v github.com/SeerUK/tid/...
```

## Usage

Full usage information can be seen by using the built-in contextual help. This can accessed by
running:

```
$ tid --help
$ # Or for command help:
$ tid start --help
```

Here's some simple general usage:

```
$ tid start "Working on AI"
Started tracking 'Working on AI' (fdb6f0d)

$ tid status
...

$ tid stop
Stopped tracking 'Working on AI' (fdb6f0d)

# Gone for lunch...

$ tid resume
Resumed tracking 'Working on AI' (fdb6f0d)

$ tid report
...
```

(`...` is used in place of larger output)

### Starting an Entry Timer

```
$ tid start "A note"
```

The note is required, but can by any string value. It's used so when you view the status or the
report you know what you've been tracking. Try make it something identifiable. Maybe this will just
be an issue ID from your issue tracker?

### Stopping an Entry Timer

```
$ tid stop
```

Stop always stops the currently active timer, if you don't have an active timer, it won't do
anything.

### Resuming an Entry Timer

```
$ tid resume
$ tid resume fdb6f0d
```

The resume command allows you to resume the most recently stopped entry, or a specific entry by
passing in that entry's hash. If you don't have a most recently stopped entry then you would have to
pass in an entry hash to use resume (e.g. if you remove the entry being tracked).

### Status of an Entry

```
$ tid status
$ tid status fdb6f0d
$ tid status --format="{{.Entry.Duration}} on '{{.Entry.Note}}'"
```

You can view the status of the currently tracked entry (the most recently started or resumed entry)
or you can view the status of a specific entry. The output is similar to the report output.

### Report your Timesheet

```
$ tid report
$ tid report --start=2017-02-01 --end=2017-02-28
$ tid report --start=(tiddate --months=-6)
$ tid report --no-summary
$ tid report --format="{{.Entry.Hash}} {{.Entry.Note}}" --no-summary
```

The report command is quite powerful and gives you a lot of different ways to view timesheet data.
By default the output will display a summary, and a table of the entries. You can control the output
by passing other options like `--format` which is useful for scripting.

### Adding an Entry

```
$ tid add 2001-01-01 24h "Welcome to the 21st century!"
$ tid add (tiddate --days=-1) 1h10m "Call with Google"
```

Sometimes you just forget to track something, and maybe it was a day ago! This command lets you add
something on a specific timesheet. You can edit and remove just like normal anyway. You must always
specify a timesheet date, a duration, and a note.

### Editing an Entry

```
$ tid edit fdb6f0d --duration=30s
$ tid edit fdb6f0d --note="Working on AI killer"
$ tid edit fdb6f0d --offset=1m
```

People make mistakes, and you probably will when tracking your time too. You can update the note, or
the duration of an entry easily with the edit command. If you specify an `--offset` it will add that
to the duration (you can specify negative offsets to subtract time too, like `--offset=-12s`). You
cannot specify both a duration and an offset at the same time.

### Removing an Entry

```
$ tid remove fdb6f0d
```

Removing will permantently delete an entry. If the entry is running, you will be returned to a
stopped state. If it was being tracked most recently, you will no longer have any entry to resume.

## Completions

Completions are provided for Fish and are located with obvious names in the `completions/` 
directory. Installation will probably look something like this:

```
$ mkdir -p ~/.config/fish/completions/
$ cp completions/tid.fish ~/.config/fish/completions/
```

Completion covers commands, options of commands, and entries where applicable.

## License

MIT

[1]: https://github.com/boltdb/bolt
