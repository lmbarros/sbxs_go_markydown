package markydown

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

	parser.parseDocument()
}

// parser stores all the parsing state.
type parser struct {
	input     string    // The input that was not consumed yet.
	processor Processor // Processor processing the parsed data.
	frag      string    // The current text fragment being parsed, along with the rest of the input
	fragEnd   int       // Index into frag indicating the end of fragment being parsed
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

	// TODO: Try parsing other paragraph types

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

// parseParagraphContents parses the contents of a paragraph. The input must be
// on the first character of the contents (that is, things like `# ` and `+ `
// that mark the paragraph type must have been consumed already).
func (p *parser) parseParagraphContents() {

	p.frag = p.input
	p.fragEnd = 0

	for initialLen := len(p.input); ; initialLen = len(p.input) {
		//initialLen := len(p.input)
		// runeType, theRune, isEscaped := p.nextRune()
		theType, _, isEscaped := p.nextRune()

		switch theType {

		case runeTypeSpace:
			p.emitFragment()
			if p.consumeRawSpacesWithinParagraph() {
				p.processor.onSpecialToken(SpecialTokenSpace)
			}

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
// parse a new fragment. This is a no-op if the current fragment is empty.
func (p *parser) emitFragment() {
	if p.fragEnd > 0 {
		p.processor.onFragment(p.frag[:p.fragEnd])
		p.fragEnd = 0
		p.frag = p.input
	}
}
