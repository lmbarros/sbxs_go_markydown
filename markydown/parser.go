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
	p.processor.OnStartDocument()
	defer p.processor.OnEndDocument()

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

// emitFragment tells the processor that the current text fragment was parsed
// and resets the fragment-related parser state, in order to make it ready to
// parse a new fragment.
//
// If the current fragment is empty, this will not emit anything, but will reset
// the internal state so that we start a new fragment from the current point in
// the input.
func (p *parser) emitFragment() {
	if p.fragEnd > 0 {
		p.processor.OnFragment(p.frag[:p.fragEnd])
	}

	p.fragEnd = 0
	p.frag = p.input
}
