package markydown

import "unicode/utf8"

// nextRune consumes and returns the next rune from the input.
//
// We could perhaps say that this function does the lexing in this
// implementation. Perhaps. Not all input passes through this function -- we
// take some shortcuts here and there -- and it has some smarts that I wouldn't
// expect in a real lexer. (Particularly when handling links; we do quite a bit
// of work here to simplify the work on the parser.)
//
// Apart from the links case mentioned above, this function doesn't know
// much the Markydown syntax.
//
// This function handles escaped characters and handles line breaks smartly (to
// deal with all that CRLF x CR x whatever mess).
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

	case isEmphasis(r):
		r, w = utf8.DecodeRuneInString(p.input)

		if isEmphasis(r) {
			p.input = p.input[w:]
			return runeTypeStrongEmphasis, '*', false
		}

		return runeTypeEmphasis, '*', false

	case isLinkStart(r):
		if p.lookAheadForLink() {
			return runeTypeLinkStart, r, false
		}
		return runeTypeText, r, false

	case isLinkEnd(r):
		if len(p.linkTarget) > 0 {
			return runeTypeLinkEnd, r, false
		}
		return runeTypeText, r, false

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
