package markydown

import "unicode/utf8"

// nextRune consumes and returns the next rune from the input.
//
// We could perhaps say that this function does the lexing in this
// implementation. It handles escaped characters and handles line breaks smartly
// (to deal with all that CRLF x CR x whatever mess), but it doesn't know
// anything about the Markydown syntax. For instance, if it returns a result
// saying that the next rune is a heading marker, it doesn't mean we have a
// heading at that point of the input -- it might be just a hash sign used in
// the middle of a paragraph.
//
// Three values are returned. The first is the rune type; as explained above,
// a `#` is always identified as a `runeTypeHeading`, even if it is not being
// used to represent a heading at that point. (And the same logic holds for
// other character classes.) The second value is the rune itself. And the third
// value indicates whether the lexed rune was escaped or not.
func (p *parser) nextRune() (runeType, rune, bool) {

	// Do we still have input data?
	if len(p.input) == 0 {
		return runeTypeEOI, utf8.RuneError, false
	}

	// Decode next rune
	r, w := utf8.DecodeRuneInString(p.input)
	p.input = p.input[w:]

	// Handle it
	switch {

	case isHorizontalSpace(r):
		return runeTypeSpace, ' ', false

	case isEscape(r):
		r, w = utf8.DecodeRuneInString(p.input)
		p.input = p.input[w:]
		if isNewLine(r) {
			firstNewLine := r
			r, w = utf8.DecodeRuneInString(p.input)
			if isNewLine(r) && r != firstNewLine { // either CRLF or LFCR
				p.input = p.input[w:]
			}
			return runeTypeNewLine, '\n', true
		}
		return runeTypeText, r, true

	case isNewLine(r):
		firstNewLine := r
		r, w = utf8.DecodeRuneInString(p.input)
		if isNewLine(r) && r != firstNewLine { // either CRLF or LFCR
			p.input = p.input[w:]
		}
		return runeTypeNewLine, '\n', false

	default:
		return runeTypeText, r, false
	}
}
