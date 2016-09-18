// Package markydown provides a parser for Markydown, a small subset of Markdown
// that is enough for my needs. YMMV.
//
// This is not a Markydown-to-HTML (or to any other format, for that matter)
// converter. The provided parser will simply call some methods of a
// user-supplied `Processor` object as it detects, for example, that a new
// paragraph started, the formatting changed or some text is to be "emited".
package markydown
