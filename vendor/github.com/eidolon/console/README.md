# Console [![Travis Build Status][travis-badge]][travis-build]

Facilitates the process of creating complex command-line applications in Go.

## Usage

You can see usage examples in the `_examples/` folder in this repository.

## Motivation

Command line applications provide an easy to use interface for interacting with your application. Go
has great support for making these types of applications built into the standard library via the 
flags package, but it doesn't really facilitate the development of complex command line 
applications.

This library is designed to help make more complex, safe console applications, in a consistent,
simple, and easy to use way. It's designed to be lightweight, but powerful and configurable. Another 
goal of this library is to make it so that commands are easily testable.

## Todo List

There are still some things I'd like to get done with this library, as is reflected by the pre-v1.0
state. Here's a priority ordered todo list:

* Documentation.
* Mutually exclusive options.
* More complete set of tests.
* More helpful `Input` type.
* Test helpers.
* Multiple environment variables per-option.

## License

MIT

## Contributions

Feel free to open a [pull request][1], or file an [issue][2] on Github. I always welcome 
contributions as long as they're for the benefit of all (potential) users of this project.

If you're unsure about anything, feel free to ask about it in an issue before you get your heart set
on fixing it yourself.

[1]: https://github.com/eidolon/console/pulls
[2]: https://github.com/eidolon/console/issues

[travis-badge]: https://img.shields.io/travis/eidolon/console.svg
[travis-build]: https://travis-ci.org/eidolon/console
