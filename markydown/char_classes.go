package markydown

import "unicode"

// isNewLine checks if a givn rune is a new line.
func isNewLine(r rune) bool {
	return r == '\n' || r == '\r'
}

// isHorizontalSpace checks if a given rune is an horizontal space.
func isHorizontalSpace(r rune) bool {
	return !isNewLine(r) && unicode.IsSpace(r)
}

// isBullet checks if a given rune can be used as bullet in a bulleted list.
func isBullet(r rune) bool {
	return r == '+'
}

// isEmphasis checks if a given rune can be used to emphasize text.
func isEmphasis(r rune) bool {
	return r == '*'
}

// isEscape checks if a given rune can be used to escape other runes.
func isEscape(r rune) bool {
	return r == '\\'
}

// isLinkStart checks if a given rune can be used to start a link.
func isLinkStart(r rune) bool {
	return r == '['
}

// isLinkEnd checks if a given rune can be used to end a link.
func isLinkEnd(r rune) bool {
	return r == ']'
}

// isLinkTargetStart checks if a given rune can be used to start a link target.
func isLinkTargetStart(r rune) bool {
	return r == '('
}

// isLinkTargetEnd checks if a given rune can be used to end a link target.
func isLinkTargetEnd(r rune) bool {
	return r == ')'
}
