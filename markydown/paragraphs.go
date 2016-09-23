package markydown

import (
	"strings"
	"unicode/utf8"
)

// parseTextParagraph parses a regular text paragraph.
//
// This is used as a last resort: if the paragraph is not of any other paragraph
// types, we use this one. It is supposed to succeed no matter what, which is
// the reason why it doesn't return a Boolean indicating success or failure.
func (p *parser) parseTextParagraph() {

	p.processor.StartParagraph(ParTypeText)
	defer p.processor.EndParagraph(ParTypeText)

	p.parseParagraphContents()
}

// parseHeading parses a heading (of any supported level). Returns true if the
// parsing succeeded or false otherwise (in which case no input is consumed).
func (p *parser) parseHeading() bool {
	parType := ParTypeInvalid

	firstSpace := strings.IndexFunc(p.input, isHorizontalSpace)

	if strings.HasPrefix(p.input, "###") && firstSpace == 3 {
		parType = ParTypeHeading3
	} else if strings.HasPrefix(p.input, "##") && firstSpace == 2 {
		parType = ParTypeHeading2
	} else if strings.HasPrefix(p.input, "#") && firstSpace == 1 {
		parType = ParTypeHeading1
	}

	if parType == ParTypeInvalid {
		return false
	}

	p.input = p.input[firstSpace:]
	p.consumeRawHorizontalSpaces()

	p.processor.StartParagraph(parType)
	defer p.processor.EndParagraph(parType)

	p.parseParagraphContents()

	return true
}

// parseBulletedParagraph parses a paragraph that is a bulleted list item.
// Returns true if the parsing succeeded or false otherwise (in which case no
// input is consumed).
func (p *parser) parseBulletedParagraph() bool {
	firstSpace := strings.IndexFunc(p.input, isHorizontalSpace)

	r, w := utf8.DecodeRuneInString(p.input)

	if !isBullet(r) || firstSpace != w {
		return false
	}

	p.input = p.input[w:]
	p.consumeRawHorizontalSpaces()

	p.processor.StartParagraph(ParTypeBulletedList)
	defer p.processor.EndParagraph(ParTypeBulletedList)

	p.parseParagraphContents()

	return true
}

// parseParagraphContents parses the contents of a paragraph. The input must be
// on the first character of the contents (that is, things like `# ` and `+ `
// that mark the paragraph type must have been consumed already).
func (p *parser) parseParagraphContents() {

	p.frag = p.input
	p.fragEnd = 0

	for initialLen := len(p.input); ; initialLen = len(p.input) {
		theType, isEscaped := p.nextRune()

		switch theType {

		case runeTypeSpace:
			p.emitFragment()
			p.consumeRawSpacesWithinParagraph()
			if p.paragraphGoesOn() && !p.isHardLineBreakAhead() {
				p.processor.SpecialToken(SpecialTokenSpace)
			}

		case runeTypeEmphasis:
			p.emitFragment()

			if p.textStyle == TextStyleEmphasis {
				p.textStyle = TextStyleRegular
			} else {
				p.textStyle = TextStyleEmphasis
			}

			p.processor.ChangeTextStyle(p.textStyle)

		case runeTypeStrongEmphasis:
			p.emitFragment()

			if p.textStyle == TextStyleStrong {
				p.textStyle = TextStyleRegular
			} else {
				p.textStyle = TextStyleStrong
			}

			p.processor.ChangeTextStyle(p.textStyle)

		case runeTypeNewLine:
			p.emitFragment()
			p.consumeRawHorizontalSpaces()

			if isEscaped {
				p.processor.SpecialToken(SpecialTokenLineBreak)
			} else if p.paragraphGoesOn() {
				p.processor.SpecialToken(SpecialTokenSpace)
			}

			if len(p.input) == 0 {
				// End of input: return, as this is also the end of the paragraph
				return
			}
			r, _ := utf8.DecodeRuneInString(p.input)
			if isNewLine(r) {
				// Two consecutive new lines: we reached the end of the paragraph
				return
			}

		case runeTypeLinkStart:
			p.emitFragment()
			p.processor.StartLink(p.linkTarget)

		case runeTypeLinkEnd:
			p.emitFragment()
			p.processor.EndLink()
			p.consumeLinkTarget()

		case runeTypeEOI:
			p.emitFragment()
			return

		case runeTypeText:
			if isEscaped {
				p.frag = p.frag[:p.fragEnd] + p.frag[p.fragEnd+1:]
				p.fragEnd--
			}
			p.fragEnd += initialLen - len(p.input)
		}
	}
}
