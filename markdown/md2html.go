package markdown

import (
	"regexp"
	"strings"
)

var reBold, reItalic, reStrike, reLink, reCode, reImg *regexp.Regexp

type Parser struct {
	html  string
	pre   bool
	quote bool
	li    bool
}

func init() {
	// text **in bold**
	reBold = regexp.MustCompile("(\\*\\*)(.*?)(\\*\\*)")

	// text *in italic*
	reItalic = regexp.MustCompile("(\\*)(.*?)(\\*)")

	// ~~striked text~~
	reStrike = regexp.MustCompile("(~~)(.*?)(~~)")

	// [some text](http://somewhere.org)
	// \\[							starts with [
	//    ([^\\]]*?)		followed by any character except ] using lazy match
	// \\]							followed by an ending ]
	// \\(							followed by (
	//    (.*?)					followed by any character using lazy match
	// \\)      				ends with )
	//
	// Using ([^\\]]*?) instead of (.*?) to prevent the parser from
	// picking up text in brackets that is not an URL, for example:
	// [1] [book x](http://link/to/bookx)
	reLink = regexp.MustCompile("\\[([^\\]]*?)\\]\\((.*?)\\)")

	// `hello world`
	reCode = regexp.MustCompile("(`)(.*?)(`)")

	// ![image caption](http://somewhere.org/image.png)
	reImg = regexp.MustCompile("!\\[([^\\]]*?)\\]\\((.*?)\\)")
}

func (p Parser) ToHtml(markdown string) string {
	p.html = ""
	p.pre = false
	p.quote = false
	p.li = false
	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		if p.isQuote(line) {
			if p.quote == true {
				// already in blockquote
			} else {
				// start a new blockquote
				p.html += "<blockquote>\n"
				p.quote = true
			}
		} else {
			if p.quote == true {
				// end current blockquote
				p.html += "</blockquote>\n"
				p.quote = false
			}
		}

		if p.isListItem(line) {
			if p.li == true {
				// already inside a list
			} else {
				// start a new list
				p.html += "<ul>\n"
				p.li = true
			}
		} else {
			if p.li == true {
				// end current list
				p.html += "</ul>\n"
				p.li = false
			}
		}

		l := strings.TrimSpace(line)

		if p.isH1(l) {
			p.html += "<h1>" + substr(l, 2) + "</h1>\n"
		} else if p.isH2(l) {
			p.html += "<h2>" + substr(l, 3) + "</h2>\n"
		} else if p.isH3(l) {
			p.html += "<h3>" + substr(l, 4) + "</h3>\n"
		} else if p.isPreTerminal(l) {
			p.html += "<pre class=\"terminal\">\n"
			p.pre = true
		} else if p.isPreCode(l) {
			p.html += "<pre class=\"code\">\n"
			p.pre = true
		} else if p.isPre(l) {
			if p.pre {
				p.html += "</pre>\n"
				p.pre = false
			} else {
				p.html += "<pre>\n"
				p.pre = true
			}
		} else if l == "" {
			// html += "<br/>\n"
			p.html += "\n"
		} else {
			if p.pre {
				// we use the original line in pre to preserve spaces
				p.html += p.inline(line) + "\n"
			} else if p.quote {
				p.html += p.inline(substr(l, 2)) + "<br/>\n"
			} else if p.li {
				p.html += "<li>" + p.inline(substr(l, 2)) + "\n"
			} else {
				p.html += "<p>" + p.inline(l) + "</p>\n"
			}
		}
	}
	return p.html
}

func (p Parser) isH1(line string) bool {
	if p.pre {
		return false
	}
	return strings.HasPrefix(line, "# ")
}

func (p Parser) isH2(line string) bool {
	if p.pre {
		return false
	}
	return strings.HasPrefix(line, "## ")
}

func (p Parser) isH3(line string) bool {
	if p.pre {
		return false
	}
	return strings.HasPrefix(line, "### ")
}

func substr(line string, i int) string {
	return line[i:len(line)]
}

func (p Parser) isPreTerminal(line string) bool {
	return strings.HasPrefix(line, "```terminal")
}

func (p Parser) isPreCode(line string) bool {
	return strings.HasPrefix(line, "```code")
}

func (p Parser) isPre(line string) bool {
	return strings.HasPrefix(line, "```")
}

func (p Parser) isQuote(line string) bool {
	return strings.HasPrefix(line, "> ")
}

func (p Parser) isListItem(line string) bool {
	return strings.HasPrefix(line, "* ")
}

func (p Parser) inline(line string) string {
	// TODO: encode & to &amp;
	line = strings.Replace(line, "<", "&lt;", -1)
	line = strings.Replace(line, ">", "&gt;", -1)

	// allow for <img ... />
	line = strings.Replace(line, "&lt;img ", "<img ", -1)
	line = strings.Replace(line, "/&gt;", "/>", -1)

	// allow for <sup> </sup>
	line = strings.Replace(line, "&lt;sup&gt;", "<sup>", -1)
	line = strings.Replace(line, "&lt;/sup&gt;", "</sup>", -1)

	line = reImg.ReplaceAllString(line, "<img src=\"$2\" alt=\"$1\" title=\"$1\" />")
	line = reBold.ReplaceAllString(line, "<b>$2</b>")
	line = reItalic.ReplaceAllString(line, "<i>$2</i>")
	line = reStrike.ReplaceAllString(line, "<strike>$2</strike>")
	line = reLink.ReplaceAllString(line, "<a href=\"$2\">$1</a>")
	line = reCode.ReplaceAllString(line, "<code>$2</code>")
	return line
}
