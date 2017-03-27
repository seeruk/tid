# Console [![Travis Build Status][travis-badge]][travis-build]

Facilitates the process of creating complex command-line applications in Go.

## Usage

You can see usage examples in the `_examples/` folder in this repository.

## Motivation

Command line applications provide an easy to use interface for interacting with your application. Go
has great support for making these types of applications built into the standard library via the 
flags package, but it doesn't really facilitate the development of complex command line 
applications.

This library is designed to help make more complex, type-safe console applications, in a consistent,
simple, and easy to use way. It's designed to be lightweight and simple, but powerful and 
configurable. Another goal of this library is to make it so that commands are easily testable.

## Todo List

Consider this library pre-v1.0. In other words, there may be backwards compatibility breaking 
changes. If this bothers you, use vendoring.

* Better usage documentation.
* Better overall application testing. Although all of the code is covered, there are still a few 
more tests that could be written.
* More helpful `Input` type (but one that still maintains testability).
* Command testing helpers. Currently commands can be tested, however the tests can be a little 
verbose. I need to plan around how I can make command tests a little easier to write.

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
