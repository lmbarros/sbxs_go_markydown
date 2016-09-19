package markydown

import (
	"unicode"
	"unicode/utf8"
)

// consumeRawSpaces chomps spaces from the input until some non-space is found.
//
// The "Raw" in the name is to indicate that we are looking at input characters
// without any special handling.
//
// We don't have to check for escape characters here because escaped spaces or
// new lines have special meaning and must be handled by the main processing
// loop.
func (p *parser) consumeRawSpaces() {
	for len(p.input) > 0 {
		r, w := utf8.DecodeRuneInString(p.input)

		if !unicode.IsSpace(r) {
			break
		}

		p.input = p.input[w:]
	}
}

// consumeRawHorizontalSpaces chomps horizontal spaces from the input.
func (p *parser) consumeRawHorizontalSpaces() {
	for len(p.input) > 0 {
		r, w := utf8.DecodeRuneInString(p.input)

		if !isHorizontalSpace(r) {
			break
		}

		p.input = p.input[w:]
	}
}

// consumeRawSpacesWithinParagraph chomps all spaces from the input, but do not
// slip to the next paragraph.
//
// Returns true if the paragraph goes on, or false if the paragraph is ending
// here.
func (p *parser) consumeRawSpacesWithinParagraph() bool {
	p.consumeRawHorizontalSpaces()
	r, w := utf8.DecodeRuneInString(p.input)
	if isNewLine(r) {
		p.input = p.input[w:]

		firstNewLine := r
		r, w = utf8.DecodeRuneInString(p.input)
		if isNewLine(r) && r != firstNewLine { // either CRLF or LFCR
			p.input = p.input[w:]
		}

		p.consumeRawHorizontalSpaces()
	}

	p.frag = p.input

	// Does the current paragraph go on?
	if len(p.input) == 0 {
		return false
	}

	r, _ = utf8.DecodeRuneInString(p.input)
	return !isNewLine(r)
}
