package markydown

import "unicode/utf8"

// nextRune consumes and returns the next rune from the input.
//
// We could perhaps say that this function does the lexing in this
// implementation. (Though not all input passes through this function, we take
// some shortcuts here and there) It handles escaped characters and handles line
// breaks smartly (to deal with all that CRLF x CR x whatever mess), but it
// doesn't know anything about the Markydown syntax. For instance, if it returns
// a result saying that the next rune is a heading marker, it doesn't mean we
// have a heading at that point of the input -- it might be just a hash sign
// used in the middle of a paragraph.
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
			return runeTypeLinkOpen, r, false
		}
		return runeTypeText, r, false

	case isLinkEnd(r):
		if len(p.linkTarget) > 0 {
			return runeTypeLinkClose, r, false
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

// lookAheadForLink first looks ahead to detect if the `[` we just found is
// really a link. Then, if we are indeed parsing a link, it looks ahead a bit
// further to obtain the link target.
//
// Returns true if we are parsing a link, false otherwise.
func (p *parser) lookAheadForLink() bool {
	// We are only looking ahead, use a copy and leave the real input untouched
	input := p.input

	for {
		r, w := utf8.DecodeRuneInString(input)

		switch {
		case len(input) == 0:
			return false

		case isLinkEnd(r):
			input = input[w:]
			return p.parseLinkTarget(input)

		case isEscape(r):
			input = input[w:]
			_, w = utf8.DecodeRuneInString(input)
			input = input[w:]

		default:
			input = input[w:]
		}
	}
}

// parseLinkTarget parses a link target from a given input string. Returns a
// Boolean indicating if a link target was actually found on input.
func (p *parser) parseLinkTarget(input string) bool {
	r, w := utf8.DecodeRuneInString(input)

	if !isLinkTargetStart(r) {
		return false
	}

	input = input[w:]

	target := input
	targetEnd := 0 // index into target (excludes escape runes)
	targetLen := 0 // lenght in bytes (includes escapes runes)

	for {
		r, w = utf8.DecodeRuneInString(input)

		switch {
		case len(input) == 0:
			return false

		case isLinkTargetEnd(r):
			p.linkTarget = target[:targetEnd]
			p.linkTargetLen = targetLen
			return true

		case isEscape(r):
			input = input[w:]
			target = target[:targetEnd] + target[targetEnd+1:]
			targetEnd += w
			targetLen += w

			r, w = utf8.DecodeRuneInString(input)
			input = input[w:]
			targetLen += w

		default:
			input = input[w:]
			targetEnd += w
			targetLen += w
		}
	}
}

// consumeLinkTarget chomps the link target that is expected to be right on the
// start of the input.
func (p *parser) consumeLinkTarget() {
	p.input = p.input[p.linkTargetLen+2:] // `+2` accounts for the parens themselves
	p.frag = p.input
	p.fragEnd = 0
	p.linkTarget = ""
	p.linkTargetLen = 0
}
