package markydown

// Parse parses a Markdown document passed as a string and lets the passed
// Processor do its work as the document is parsed.
//
// It works in the same spirit as the Template Method design pattern.
func Parse(document string, processor Processor) {
}

// parser stores all the parsing state.
type parser struct {
}
