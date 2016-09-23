# StackedBoxes' Markydown Parser in Go

[![GoDoc](https://godoc.org/github.com/lmbarros/sbxs_go_markydown?status.svg)](https://godoc.org/github.com/lmbarros/sbxs_go_markydown) [![Go Report Card](https://goreportcard.com/badge/github.com/lmbarros/sbxs_go_markydown)](https://goreportcard.com/report/github.com/lmbarros/sbxs_go_markydown) ![License](https://img.shields.io/github/license/lmbarros/sbxs_go_markydown.svg)

This package provides a Markydown parser in Go.

Markydown is a small subset of
[Markdown](http://daringfireball.net/projects/markdown/) that provides just
enough features for my needs, but you may consider too simple. (Technically it
is not a subset, but I thought it was similar enough to be described as one.)

This is *not* a Markydown to HTML converter. The parser takes a `Processor` as
parameter, and calls methods like `OnStartParagraph` and `OnChangeTextStyle` as
it parses its input. It's up to you to provide a `Processor` implementation that
does whatever you need. (That said, I provide an example that implements a
simple Markydown-to-HTML converter.)

## Markydown

Markydown isn't terribly well-specified. The example below should give you an
idea of what it is like. If not, though luck. That's all I have.

```
# Markydown example

## Headings are supported

### But only up to level 3

Notice that only atx-style headings are supported, and only partially, as you
can't add trailing hashes (#) to headings. Well, you can, but they will
be included as part of the heading.


Use one or more blank lines to separate
paragraphs. Leading
     and trailing spaces
are ignored.

Text can be *emphasized* or **strongly emphasized**, but you must use
asterisks -- underscores are treated as any other character. Unlike in Markdown,
an asterisk surrounded by spaces still is considered an emphasis mark. So, you
need to escape characters in things like this: 3 \* 7 = 21.

Escaped spaces are treated as non-breaking spaces, so you may want to write
things like 100\ kg and «\ Où\ ?\ ». Any character can be escaped, \O\K\?

You can force a line break by ending\
the line with a backslash. Think that you\
are escaping the newline character. The Markdown\
"trailing spaces" syntax does not work here.

You can create [links](www.example.com) but you cannot add a link title.

Bulleted lists are also supported, however:

+ They cannot be nested.

+ You must use "plus" signs as the bullets.

   + Items needn't be aligned and can
     span multiple
lines.

+ And you must leave an empty line between items.

And that's all.
```

## License

All code here is under the MIT License.
