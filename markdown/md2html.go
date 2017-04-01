package markdown

import (
	"regexp"
	"strings"
)

var reBold, reItalic, reStrike, reLink, reCode, reImg *regexp.Regexp

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

func ToHtml(markdown string) string {
	html := ""
	pre := false
	quote := false
	li := false
	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		if isQuote(line) {
			if quote == true {
				// already in blockquote
			} else {
				// start a new blockquote
				html += "<blockquote>\n"
				quote = true
			}
		} else {
			if quote == true {
				// end current blockquote
				html += "</blockquote>\n"
				quote = false
			}
		}

		if isListItem(line) {
			if li == true {
				// already inside a list
			} else {
				// start a new list
				html += "<ul>\n"
				li = true
			}
		} else {
			if li == true {
				// end current list
				html += "</ul>\n"
				li = false
			}
		}

		l := strings.TrimSpace(line)

		if isH1(l) {
			html += "<h1>" + substr(l, 2) + "</h1>\n"
		} else if isH2(l) {
			html += "<h2>" + substr(l, 3) + "</h2>\n"
		} else if isH3(l) {
			html += "<h3>" + substr(l, 4) + "</h3>\n"
		} else if isPreTerminal(l) {
			html += "<pre class=\"terminal\">\n"
			pre = true
		} else if isPreCode(l) {
			html += "<pre class=\"code\">\n"
			pre = true
		} else if isPre(l) {
			if pre {
				html += "</pre>\n"
				pre = false
			} else {
				html += "<pre>\n"
				pre = true
			}
		} else if l == "" {
			// html += "<br/>\n"
			html += "\n"
		} else {
			if pre {
				// we use the original line in pre to preserve spaces
				html += inline(line) + "\n"
			} else if quote {
				html += inline(substr(l, 2)) + "<br/>\n"
			} else if li {
				html += "<li>" + inline(substr(l, 2)) + "\n"
			} else {
				html += "<p>" + inline(l) + "</p>\n"
			}
		}
	}
	return html
}

func isH1(line string) bool {
	return strings.HasPrefix(line, "# ")
}

func isH2(line string) bool {
	return strings.HasPrefix(line, "## ")
}

func isH3(line string) bool {
	return strings.HasPrefix(line, "### ")
}

func substr(line string, i int) string {
	return line[i:len(line)]
}

func isPreTerminal(line string) bool {
	return strings.HasPrefix(line, "```terminal")
}

func isPreCode(line string) bool {
	return strings.HasPrefix(line, "```code")
}

func isPre(line string) bool {
	return strings.HasPrefix(line, "```")
}

func isQuote(line string) bool {
	return strings.HasPrefix(line, "> ")
}

func isListItem(line string) bool {
	return strings.HasPrefix(line, "* ")
}

func inline(line string) string {
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
