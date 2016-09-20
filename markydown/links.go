package markydown

import "unicode/utf8"

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
	targetLen := 0 // length in bytes (includes escapes runes)

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
