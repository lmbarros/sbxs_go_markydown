package markydown

import (
	"testing"

	"github.com/lmbarros/sbxs_go_test/test/assert"
)

//
// testProcessor helpers
//

// parTypeToString converts a given ParType to a string value, as used by the
// testProcessor.
func parTypeToString(parType ParType) string {
	switch parType {
	case ParTypeText:
		return "P"
	case ParTypeHeading1:
		return "H1"
	case ParTypeHeading2:
		return "H2"
	case ParTypeHeading3:
		return "H3"
	case ParTypeBulletedList:
		return "UL"
	default:
		return "<WTF?!>"
	}
}

// specialTokenToString converts a given SpecialToken to a string value, as used
// by the testProcessor.
func specialTokenToString(token SpecialToken) string {
	switch token {
	case SpecialTokenSpace:
		return "SP"
	case SpecialTokenLineBreak:
		return "NL"
	default:
		return "<WTF?!>"
	}
}

// textStyleToString converts a given TextStyle to a string value, as used by
// the testProcessor.
func textStyleToString(style TextStyle) string {
	switch style {
	case TextStyleRegular:
		return "RE"
	case TextStyleEmphasis:
		return "EM"
	case TextStyleStrong:
		return "ST"
	default:
		return "<WTF?!>"
	}
}

//
// testProcessor itself
//

// testProcessor is a Markdown processor used for testing.
//
// It just append certain strings to its res member as its callbacks are called.
// Once parsing is complete, we can just check if the expected sequence of
// strings can be found in res.
type testProcessor struct {
	res []string
}

func (p *testProcessor) OnStartDocument() {
	p.res = append(p.res, "SD")
}

func (p *testProcessor) OnEndDocument() {
	p.res = append(p.res, "ED")
}

func (p *testProcessor) OnStartParagraph(parType ParType) {
	p.res = append(p.res, "SP-"+parTypeToString(parType))
}

func (p *testProcessor) OnEndParagraph(parType ParType) {
	p.res = append(p.res, "EP-"+parTypeToString(parType))
}

func (p *testProcessor) OnFragment(text string) {
	p.res = append(p.res, "F-"+text)
}

func (p *testProcessor) OnSpecialToken(token SpecialToken) {
	p.res = append(p.res, "ST-"+specialTokenToString(token))
}

func (p *testProcessor) OnChangeTextStyle(style TextStyle) {
	p.res = append(p.res, "TS-"+textStyleToString(style))
}

func (p *testProcessor) OnStartLink(target string) {
	p.res = append(p.res, "SL-"+target)
}

func (p *testProcessor) OnEndLink() {
	p.res = append(p.res, "EL")
}

//
// Real tests start here
//

// Tests parsing different variations of empty documents.
func TestParseBlanks(t *testing.T) {
	emptyResult := []string{"SD", "ED"}

	inputs := []string{
		"",
		" ",
		"\t",
		"   \t \t  ",
		"\n",
		"\r",
		"\n  ",
		"   \t\n\n \t\n    \t    \n\r  \r\n\t  ",
	}

	for _, v := range inputs {
		p := &testProcessor{}
		Parse(v, p)
		assert.Equal(t, p.res, emptyResult)
	}
}

// Tests parsing some single-word inputs.
func TestParseSingleWord(t *testing.T) {
	testData := map[string][]string{
		"one":          {"SD", "SP-P", "F-one", "EP-P", "ED"},
		"   två\t":     {"SD", "SP-P", "F-två", "EP-P", "ED"},
		"\n\ntrês  \n": {"SD", "SP-P", "F-três", "EP-P", "ED"},

		// Try an escaped character
		"\ndö\\rt \n": {"SD", "SP-P", "F-dört", "EP-P", "ED"},

		// A rogue backslash in the end of the input should be OK
		"fünf\\": {"SD", "SP-P", "F-fünf", "EP-P", "ED"},
	}

	for input, expected := range testData {
		p := &testProcessor{}
		Parse(input, p)
		assert.Equal(t, p.res, expected)
	}
}

// Tests parsing some simple (no formatting) paragraphs.
func TestParseSimpleParagraphs(t *testing.T) {
	testData := map[string][]string{
		// Normal text paragraph
		"bir iki":                      {"SD", "SP-P", "F-bir", "ST-SP", "F-iki", "EP-P", "ED"},
		"     água  \t seca\n  ":       {"SD", "SP-P", "F-água", "ST-SP", "F-seca", "EP-P", "ED"},
		"\n um  \n\r dois\n\r  três  ": {"SD", "SP-P", "F-um", "ST-SP", "F-dois", "ST-SP", "F-três", "EP-P", "ED"},

		// Some valid headings
		"#\tOompa    loompa":       {"SD", "SP-H1", "F-Oompa", "ST-SP", "F-loompa", "EP-H1", "ED"},
		"## doompadee\tdoo":        {"SD", "SP-H2", "F-doompadee", "ST-SP", "F-doo", "EP-H2", "ED"},
		"###    doompadah   dee  ": {"SD", "SP-H3", "F-doompadah", "ST-SP", "F-dee", "EP-H3", "ED"},

		// Need a space after the hash sign to be recognized as heading
		"#I've\n\tgot\t": {"SD", "SP-P", "F-#I've", "ST-SP", "F-got", "EP-P", "ED"},

		// Recognize only up to level-3 heading
		"#### a perfect\n": {"SD", "SP-P", "F-####", "ST-SP", "F-a", "ST-SP", "F-perfect", "EP-P", "ED"},

		// Bulleted lists
		"+ Puzzle\n\t for\n\n": {"SD", "SP-UL", "F-Puzzle", "ST-SP", "F-for", "EP-UL", "ED"},
		"+   \t you.":          {"SD", "SP-UL", "F-you.", "EP-UL", "ED"},

		// An actual list item requires a space after the bullet sign
		"+Só!": {"SD", "SP-P", "F-+Só!", "EP-P", "ED"},
	}

	for input, expected := range testData {
		p := &testProcessor{}
		Parse(input, p)
		assert.Equal(t, p.res, expected)
	}
}

// Tests parsing some simple formatting.
func TestParseSimpleFormatting(t *testing.T) {
	testData := map[string][]string{
		// Emphasis
		"*Now!*": {"SD", "SP-P", "TS-EM", "F-Now!", "TS-RE", "EP-P", "ED"},

		// Strong emphasis
		"Now, **go!**":   {"SD", "SP-P", "F-Now,", "ST-SP", "TS-ST", "F-go!", "TS-RE", "EP-P", "ED"},
		"**Now**,   go!": {"SD", "SP-P", "TS-ST", "F-Now", "TS-RE", "F-,", "ST-SP", "F-go!", "EP-P", "ED"},

		// Emphasis in the middle of words are ok
		"# Unbe*lie*vable!":   {"SD", "SP-H1", "F-Unbe", "TS-EM", "F-lie", "TS-RE", "F-vable!", "EP-H1", "ED"},
		"# Unbe**lie**vable!": {"SD", "SP-H1", "F-Unbe", "TS-ST", "F-lie", "TS-RE", "F-vable!", "EP-H1", "ED"},

		// Try some escaped emphasis characters
		"\\*Now!\\*":   {"SD", "SP-P", "F-*Now!*", "EP-P", "ED"},
		"*\\*Now!\\**": {"SD", "SP-P", "TS-EM", "F-*Now!*", "TS-RE", "EP-P", "ED"},
		"\\**Now!\\**": {"SD", "SP-P", "F-*", "TS-EM", "F-Now!*", "TS-RE", "EP-P", "ED"},
	}

	for input, expected := range testData {
		p := &testProcessor{}
		Parse(input, p)
		assert.Equal(t, p.res, expected)
	}
}

// Tests parsing some links.
func TestParseLinks(t *testing.T) {
	testData := map[string][]string{
		// Kinda boring link
		"Click [here](target).": {"SD", "SP-P", "F-Click", "ST-SP", "SL-target",
			"F-here", "EL", "F-.", "EP-P", "ED"},

		// Link with formatting (but formatting characters in the target must be ignored!)
		"+ Click [here, *please*!](the*tárgeτ*)": {"SD", "SP-UL", "F-Click", "ST-SP", "SL-the*tárgeτ*",
			"F-here,", "ST-SP", "TS-EM", "F-please", "TS-RE", "F-!", "EL", "EP-UL", "ED"},

		// Tricky escaped characters
		"### Click\\[ [her\\]e](t\\(arg\\)et).": {"SD", "SP-H3", "F-Click[", "ST-SP", "SL-t(arg)et",
			"F-her]e", "EL", "F-.", "EP-H3", "ED"},

		// Unclosed links must be handled as a regular text
		"[here": {"SD", "SP-P", "F-[here", "EP-P", "ED"},
		"x[x":   {"SD", "SP-P", "F-x[x", "EP-P", "ED"},

		// Targetless link is recognized as regular text.
		"[Kafka]": {"SD", "SP-P", "F-[Kafka]", "EP-P", "ED"},
	}

	for input, expected := range testData {
		p := &testProcessor{}
		Parse(input, p)
		assert.Equal(t, p.res, expected)
	}
}

// Tests parsing some multi-paragraph inputs.
func TestParseMultipleParagraphs(t *testing.T) {
	testData := map[string][]string{
		`foo bar
	baz

	anøther`: {
			"SD",
			"SP-P", "F-foo", "ST-SP", "F-bar", "ST-SP", "F-baz", "EP-P",
			"SP-P", "F-anøther", "EP-P",
			"ED"},

		`# Title

			+ Ïtem

			+ Another one`: {
			"SD",
			"SP-H1", "F-Title", "EP-H1",
			"SP-UL", "F-Ïtem", "EP-UL",
			"SP-UL", "F-Another", "ST-SP", "F-one", "EP-UL",
			"ED"},
	}

	for input, expected := range testData {
		p := &testProcessor{}
		Parse(input, p)
		assert.Equal(t, p.res, expected)
	}
}

// Tests parsing different kinds of newlines.
func TestParseNewLines(t *testing.T) {
	expectedResult := []string{"SD", "SP-P", "F-One", "ST-SP", "F-single", "ST-SP", "F-paragraph.", "EP-P", "ED"}

	inputs := []string{
		"One\nsingle\nparagraph.",
		"One\rsingle\rparagraph.",
		"One\n\rsingle\n\rparagraph.",
		"One\r\nsingle\r\nparagraph.",
		"One\r\nsingle\n\rparagraph.",
		"One\r\nsingle\rparagraph.",
		"One\nsingle\n\rparagraph.",
	}

	for _, v := range inputs {
		p := &testProcessor{}
		Parse(v, p)
		assert.Equal(t, p.res, expectedResult)
	}
}

// Tests inputs with hard spaces and hard line breaks.
func TestParseHardSpacing(t *testing.T) {
	testData := map[string][]string{
		"«\\ Où\\ ?\\ »": {"SD", "SP-P", "F-« Où ? »", "EP-P", "ED"},
		"blah\\ ":        {"SD", "SP-P", "F-blah ", "EP-P", "ED"},
		" blah\\  ":      {"SD", "SP-P", "F-blah ", "EP-P", "ED"},

		"line\\\nbreak":   {"SD", "SP-P", "F-line", "ST-NL", "F-break", "EP-P", "ED"},
		"line\\\rbreak":   {"SD", "SP-P", "F-line", "ST-NL", "F-break", "EP-P", "ED"},
		"line\\\r\nbreak": {"SD", "SP-P", "F-line", "ST-NL", "F-break", "EP-P", "ED"},
		"line\\\n\rbreak": {"SD", "SP-P", "F-line", "ST-NL", "F-break", "EP-P", "ED"},

		"here  \\\n there":   {"SD", "SP-P", "F-here", "ST-NL", "F-there", "EP-P", "ED"},
		"here  \\\n\\ there": {"SD", "SP-P", "F-here", "ST-NL", "F- there", "EP-P", "ED"},
		"here \\ \\\n there": {"SD", "SP-P", "F-here", "ST-SP", "F- ", "ST-NL", "F-there", "EP-P", "ED"},
	}

	for input, expected := range testData {
		p := &testProcessor{}
		Parse(input, p)
		assert.Equal(t, p.res, expected)
	}
}

// Tests parsing something looking like a real document.
func TestParseRealDocument(t *testing.T) {
	input := `# The  title

	Paragraph one.
	Still the *same paragraph*.

	## Subtitle

	+ First;

	+ [Second](http://www.example.com);

	+ **Third**, ok?

	### Sub*sub*ti\*tle

	Here \[we\] have some **more**   text\
	with        some hard  \
	breaks.

	`

	expected := []string{
		"SD",
		"SP-H1", "F-The", "ST-SP", "F-title", "EP-H1",
		"SP-P", "F-Paragraph", "ST-SP", "F-one.", "ST-SP", "F-Still", "ST-SP", "F-the", "ST-SP", "TS-EM", "F-same", "ST-SP", "F-paragraph", "TS-RE", "F-.", "EP-P",
		"SP-H2", "F-Subtitle", "EP-H2",
		"SP-UL", "F-First;", "EP-UL",
		"SP-UL", "SL-http://www.example.com", "F-Second", "EL", "F-;", "EP-UL",
		"SP-UL", "TS-ST", "F-Third", "TS-RE", "F-,", "ST-SP", "F-ok?", "EP-UL",
		"SP-H3", "F-Sub", "TS-EM", "F-sub", "TS-RE", "F-ti*tle", "EP-H3",
		"SP-P", "F-Here", "ST-SP", "F-[we]", "ST-SP", "F-have", "ST-SP", "F-some", "ST-SP", "TS-ST", "F-more", "TS-RE", "ST-SP", "F-text", "ST-NL",
		"F-with", "ST-SP", "F-some", "ST-SP", "F-hard", "ST-NL",
		"F-breaks.", "EP-P",
		"ED"}

	p := &testProcessor{}
	Parse(input, p)
	assert.Equal(t, p.res, expected)
}

//
// Benchmark
//

// noopProcessor is a Markdown processor that doesn't do anything.
type noopProcessor struct{}

func (p *noopProcessor) OnStartDocument()                  {}
func (p *noopProcessor) OnEndDocument()                    {}
func (p *noopProcessor) OnStartParagraph(parType ParType)  {}
func (p *noopProcessor) OnEndParagraph(parType ParType)    {}
func (p *noopProcessor) OnFragment(text string)            {}
func (p *noopProcessor) OnSpecialToken(token SpecialToken) {}
func (p *noopProcessor) OnChangeTextStyle(style TextStyle) {}
func (p *noopProcessor) OnStartLink(target string)         {}
func (p *noopProcessor) OnEndLink()                        {}

// Benchmarks the Markdown parser.
func BenchmarkParser(b *testing.B) {
	const input = `
	# A test text

	## Here's some text for benchmarking

    Here's some *text* few people will read, *but* will help me to evaluate
    my *Markyside*-related code. So, despite the low number of readers, this
	*matters*.

    ### And *here*, more of the *same*

    You now, I need a reasonable amount of text. Even some
    [fake links](with-fake-targets) are nice to have. So, I'll keep writing
    this nonsense until I believe this is *enough* text. **Some sentences are
    strongly emphasized**. *Others are mildly emphasized*. And still others are
    not emphasized at all.

    Maybe I should add a [real link](http://www.stackedboxes.com), too, even
    though nobody will *ever* click it. Well, *maybe* someone will. You, know,
    someone could be looking at this with an editor with clickable hyperlinks
	or something like this.  Well, who knows?

	Should I used some lorem ipsum instead? I found them so boring. *Hm*, this
	text here is also quite boring... maybe I should add a joke? Nice idea,
	here's one: A duck walks into a bar, orders a beer and says to the
	bartender: "put it on my bill!"

	In your opinion, this was:

	+ Funny

	+ **Really** funny

	+ Kinda funny

	## Some poetry

	Roses are red\
	Violets are blue\
	As Yoda would say\
	This long is getting too

    Nuff said!
	`

	p := &noopProcessor{}

	for i := 0; i < b.N; i++ {
		Parse(input, p)
	}
}
