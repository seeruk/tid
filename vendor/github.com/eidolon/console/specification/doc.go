// Package specification is home to all parameter specification parsing. It includes a lexical
// scanner, and a parser for converting specification strings into Arguments or Options.
//
// There are two implementations of the parser, one for Arguments, and one for Options. The base
// parser is tested via the two implementations. All tests in this package test the public API of
// the package.
package specification
