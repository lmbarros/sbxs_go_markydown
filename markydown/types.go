package markydown

// Processor is something that processes a Markydown document as it is parsed.
//
// The Markydown parser works kinda like in Template Method pattern: you call
// the parser and it calls Processor's methods as it parses the data.
type Processor interface {
	StartDocument()
	EndDocument()
	StartParagraph(parType ParType)
	EndParagraph(parType ParType)
	Fragment(text string)
	SpecialToken(token SpecialToken)
	ChangeTextStyle(style TextStyle)
	StartLink(target string)
	EndLink()
}

// ParType is a paragraph type.
type ParType int

const (
	// ParTypeInvalid is an invalid paragraph type.
	ParTypeInvalid ParType = iota

	// ParTypeText is a regular paragraph.
	ParTypeText

	// ParTypeHeading1 is a level-1 heading.
	ParTypeHeading1

	// ParTypeHeading2 is a level-2 heading.
	ParTypeHeading2

	// ParTypeHeading3 is a level-3 heading.
	ParTypeHeading3

	//ParTypeBulletedList is a bulleted list paragraph.
	ParTypeBulletedList
)

// TextStyle is a "semantic" style a text can be rendered in.
//
// By "semantic", I mean that this does not describe how the text is to be
// physically rendered. For example, TextStyleEmphasis says that some fragment
// of text is to be emphasized, but it doesn't tell if the text ig going to be
// in italics, bold, in a different color, or something else.
type TextStyle int

const (
	// TextStyleRegular represents regular text, without any kind of emphasis.
	TextStyleRegular TextStyle = iota

	// TextStyleEmphasis represents emphasized text.
	TextStyleEmphasis

	// TextStyleStrong represents strongly emphasized text.
	TextStyleStrong
)

// SpecialToken is something that is not a text or paragraph and needs special
// handling.
type SpecialToken int

const (
	// SpecialTokenSpace represents a blank space -- one of those invisible
	// things we use to separate words, you know.
	SpecialTokenSpace SpecialToken = iota

	// SpecialTokenLineBreak represents a hard line break within the same
	// paragraph.
	SpecialTokenLineBreak
)

// runeType represents a rune type.
type runeType int

const (
	runeTypeText runeType = iota
	runeTypeEOI           // End of input
	runeTypeEmphasis
	runeTypeStrongEmphasis
	runeTypeSpace
	runeTypeNewLine
	runeTypeLinkStart
	runeTypeLinkEnd
)
