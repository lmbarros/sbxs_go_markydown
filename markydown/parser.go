package markydown

import (
	"strings"
	"unicode/utf8"
)

// Parse parses a Markdown document passed as a string and lets the passed
// Processor do its work as the document is parsed.
//
// It works in the same spirit as the Template Method design pattern.
func Parse(document string, processor Processor) {
	parser := &parser{}

	parser.input = document
	parser.processor = processor
	parser.frag = parser.input
	parser.fragEnd = 0
	parser.textStyle = TextStyleRegular
	parser.linkTarget = ""
	parser.linkTargetLen = 0

	parser.parseDocument()
}

// parser stores all the parsing state.
type parser struct {
	input         string    // The input that was not consumed yet.
	processor     Processor // Processor processing the parsed data.
	frag          string    // The current text fragment being parsed, along with the rest of the input
	fragEnd       int       // Index into frag indicating the end of fragment being parsed
	textStyle     TextStyle // The current text style
	linkTarget    string    // The current link target; if empty, we are not parsing a link
	linkTargetLen int       // The length of the target link in the input string (may be greater than len(linkTarget), because of escapes)
}

// parseDocument parses the whole Markydown document.
func (p *parser) parseDocument() {
	p.processor.onStartDocument()
	defer p.processor.onEndDocument()

	for p.parseAnyParagraph() {
		continue
	}
}

// parseAnyParagraph detects the type of the next paragraph on the input and
// parses it.
//
// Returns true if it could parse an actual paragraph, or false if the end of
// input was reached.
func (p *parser) parseAnyParagraph() bool {
	// Chomp spaces, check if something is left
	p.consumeRawSpaces()
	if len(p.input) == 0 {
		return false
	}

	// Try parsing each of the "special" paragraph types.
	if p.parseHeading() {
		return true
	}

	if p.parseBulletedParagraph() {
		return true
	}

	// If everything fails, parse as a regular text paragraph
	p.parseTextParagraph()
	return true
}

// parseTextParagraph parses a regular text paragraph.
//
// This is used as a last resort: if the paragraph is not of any other paragraph
// types, we use this one. It is supposed to succeed no matter what, which is
// the reason why it doesn't return a Boolean indicating success or failure.
func (p *parser) parseTextParagraph() {

	p.processor.onStartParagraph(ParTypeText)
	defer p.processor.onEndParagraph(ParTypeText)

	p.parseParagraphContents()
}

// parseHeading parses a heading (of any supported level). Returns true if the
// parsing suceeded or false otherwise (in which case no input is consumed).
func (p *parser) parseHeading() bool {
	parType := parTypeInvalid

	firstSpace := strings.IndexFunc(p.input, isHorizontalSpace)

	if strings.HasPrefix(p.input, "###") && firstSpace == 3 {
		parType = ParTypeHeading3
	} else if strings.HasPrefix(p.input, "##") && firstSpace == 2 {
		parType = ParTypeHeading2
	} else if strings.HasPrefix(p.input, "#") && firstSpace == 1 {
		parType = ParTypeHeading1
	}

	if parType == parTypeInvalid {
		return false
	}

	p.input = p.input[firstSpace:]
	p.consumeRawHorizontalSpaces()

	p.processor.onStartParagraph(parType)
	defer p.processor.onEndParagraph(parType)

	p.parseParagraphContents()

	return true
}

// parseBulletedParagraph parses a paragraph that is a bulleted list item.
// Returns true if the parsing suceeded or false otherwise (in which case no
// input is consumed).
func (p *parser) parseBulletedParagraph() bool {
	firstSpace := strings.IndexFunc(p.input, isHorizontalSpace)

	r, w := utf8.DecodeRuneInString(p.input)

	if !isBullet(r) || firstSpace != w {
		return false
	}

	p.input = p.input[w:]
	p.consumeRawHorizontalSpaces()

	p.processor.onStartParagraph(ParTypeBulletedList)
	defer p.processor.onEndParagraph(ParTypeBulletedList)

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
		theType, _, isEscaped := p.nextRune()

		switch theType {

		case runeTypeSpace:
			p.emitFragment()
			p.consumeRawSpacesWithinParagraph()
			if p.paragraphGoesOn() {
				p.processor.onSpecialToken(SpecialTokenSpace)
			}

		case runeTypeEmphasis:
			p.emitFragment()

			if p.textStyle == TextStyleEmphasis {
				p.textStyle = TextStyleRegular
			} else {
				p.textStyle = TextStyleEmphasis
			}

			p.processor.onChangeTextStyle(p.textStyle)

		case runeTypeStrongEmphasis:
			p.emitFragment()

			if p.textStyle == TextStyleStrong {
				p.textStyle = TextStyleRegular
			} else {
				p.textStyle = TextStyleStrong
			}

			p.processor.onChangeTextStyle(p.textStyle)

		case runeTypeNewLine:
			p.emitFragment()
			p.consumeRawHorizontalSpaces()
			if p.paragraphGoesOn() {
				p.processor.onSpecialToken(SpecialTokenSpace)
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

		case runeTypeLinkOpen:
			p.emitFragment()
			p.processor.onStartLink(p.linkTarget)

		case runeTypeLinkClose:
			p.emitFragment()
			p.processor.onEndLink()
			p.consumeLinkTarget()

		case runeTypeEOI:
			p.emitFragment()
			return

		default:
			if isEscaped {
				p.frag = p.frag[:p.fragEnd] + p.frag[p.fragEnd+1:]
				p.fragEnd--
			}
			p.fragEnd += initialLen - len(p.input)
		}
	}
}

// emitFragment tells the processor that the current text fragment was parsed
// and resets the fragment-related parser state, in order to make it ready to
// parse a new fragment.
//
// If the current fragment is empty, this will not emit anything, but will reset
// the internal state so that we start a new fragment from the current point in
// the input.
func (p *parser) emitFragment() {
	if p.fragEnd > 0 {
		p.processor.onFragment(p.frag[:p.fragEnd])
	}

	p.fragEnd = 0
	p.frag = p.input
}
