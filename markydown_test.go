package markydown_test

import (
	"fmt"

	"github.com/lmbarros/sbxs_go_markydown"
)

// Example shows how to use the parser to create a simple Markydown-to-HTML
// converter. This gives an idea on how the parser works, but don't expect
// this coverter to be foolproof -- it is not!
func Example() {
	doc := `
	# The title

	Some *text*.

	## Subtitle

	+ **Strong** item.

	+ [Linked](http://www.stackedboxes.com) item.

	### Subsubtitle

	Some\ text\
	and more.
	`

	var processor = newHTMLProcessor()
	markydown.Parse(doc, processor)

	// Output:
	// <html>
	// <body>
	// <h1>The title</h1>
	// <p>Some <em>text</em>.</p>
	// <h2>Subtitle</h2>
	// <ul>
	// <li><strong>Strong</strong> item.</li>
	// <li><a href="http://www.stackedboxes.com">Linked</a> item.</li>
	// </ul>
	// <h3>Subsubtitle</h3>
	// <p>Some text<br>and more.</p>
	// </body>
	// </html>
}

//
// Below is the definition of the HTML Processor that was passed to the parser
// above. Not pretty, but hopefully it has its educational value :-)
//

// htmlProcessor is a quick and dirty `markydown.Processor` that outputs HTML.
type htmlProcessor struct {
	parType   markydown.ParType   // The current paragraph type
	textStyle markydown.TextStyle // The current text style
}

func newHTMLProcessor() *htmlProcessor {
	p := &htmlProcessor{}
	p.textStyle = markydown.TextStyleRegular
	return p
}

func (p *htmlProcessor) StartDocument() {
	fmt.Println("<html>\n<body>")
}

func (p *htmlProcessor) EndDocument() {
	fmt.Println("</body>\n</html>")

}

func (p *htmlProcessor) StartParagraph(parType markydown.ParType) {

	if parType == markydown.ParTypeBulletedList && p.parType != markydown.ParTypeBulletedList {
		fmt.Println("<ul>")
	} else if parType != markydown.ParTypeBulletedList && p.parType == markydown.ParTypeBulletedList {
		fmt.Println("</ul>")
	}

	switch parType {
	case markydown.ParTypeText:
		fmt.Print("<p>")
	case markydown.ParTypeHeading1:
		fmt.Print("<h1>")
	case markydown.ParTypeHeading2:
		fmt.Print("<h2>")
	case markydown.ParTypeHeading3:
		fmt.Print("<h3>")
	case markydown.ParTypeBulletedList:
		fmt.Print("<li>")
	}

	p.parType = parType
}

func (p *htmlProcessor) EndParagraph(parType markydown.ParType) {
	switch parType {
	case markydown.ParTypeText:
		fmt.Println("</p>")
	case markydown.ParTypeHeading1:
		fmt.Println("</h1>")
	case markydown.ParTypeHeading2:
		fmt.Println("</h2>")
	case markydown.ParTypeHeading3:
		fmt.Println("</h3>")
	case markydown.ParTypeBulletedList:
		fmt.Println("</li>")
	}
}

func (p *htmlProcessor) Fragment(text string) {
	fmt.Print(text)
}

func (p *htmlProcessor) SpecialToken(token markydown.SpecialToken) {
	switch token {
	case markydown.SpecialTokenLineBreak:
		fmt.Print("<br>")
	case markydown.SpecialTokenSpace:
		fmt.Print(" ")
	}
}

func (p *htmlProcessor) ChangeTextStyle(style markydown.TextStyle) {
	// Close previous style
	if p.textStyle == markydown.TextStyleEmphasis {
		fmt.Print("</em>")
	} else if p.textStyle == markydown.TextStyleStrong {
		fmt.Print("</strong>")
	}

	// Start next style
	if style == markydown.TextStyleEmphasis {
		fmt.Print("<em>")
	} else if style == markydown.TextStyleStrong {
		fmt.Print("<strong>")
	}

	p.textStyle = style
}

func (p *htmlProcessor) StartLink(target string) {
	fmt.Print("<a href=\"" + target + "\">")
}
func (p *htmlProcessor) EndLink() {
	fmt.Print("</a>")
}
