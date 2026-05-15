package render

import (
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

const noBreakSpace = '\u00A0'

// Render plain text from HTML DOM with awareness of block-level elements
// https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/innerText
func InnerText(doc *html.Node) string {
	var tb innerTextBuilder

	// Output a string expression as TeX.
	renderTeX := func(expr string, isBlock bool) {
		if isBlock {
			tb.WriteNewline(2, true)
			tb.WritePre("$$\n")
		} else {
			tb.WritePre("$")
		}
		tb.WritePre(strings.TrimSpace(expr))
		if isBlock {
			tb.WritePre("\n$$")
			tb.WriteNewline(2, true)
		} else {
			tb.WritePre("$")
		}
	}

	var render func(*html.Node, bool)
	render = func(n *html.Node, keepWhitespace bool) {
		if n.Type == html.TextNode {
			if keepWhitespace {
				tb.WritePre(n.Data)
			} else {
				// write each word to innerTextBuilder
				startOfWord := -1
				for i, r := range n.Data {
					if unicode.IsSpace(r) {
						if startOfWord >= 0 {
							tb.WriteWord(n.Data[startOfWord:i])
							startOfWord = -1
						}
						if r == noBreakSpace {
							tb.QueueSpace(noBreakSpace)
						} else {
							tb.QueueSpace(' ')
						}
					} else if startOfWord < 0 {
						startOfWord = i
					}
				}
				if startOfWord >= 0 {
					tb.WriteWord(n.Data[startOfWord:])
				}
			}
			return
		}

		if n.Type == html.ElementNode {
			if isHiddenElement(n) {
				return
			}
			switch n.Data {
			// These elements will never contain user-facing text nodes, so there is no need to
			// recurse into them.
			case "head", "meta", "style", "iframe",
				"audio", "video", "track", "source", "canvas", "svg", "map", "area":
				return
			case "script":
				if ok, isBlock := isMathjaxScript(n); ok {
					renderTeX(textContent(n), isBlock)
				}
				return
			case "math":
				isBlock := isDisplay(n, "block")
				if el := findAnnotation(n, "application/x-tex"); el != nil {
					renderTeX(textContent(el), isBlock)
				}
				return
			case "mjx-container":
				isBlock := isDisplay(n, "true")
				if tex := findLatex(n); tex != "" {
					renderTeX(tex, isBlock)
				}
				return
			case "br":
				tb.WriteNewline(1, false)
			case "hr", "p", "blockquote", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "dl", "table":
				tb.WriteNewline(2, true)
			case "pre":
				tb.WriteNewline(2, true)
				keepWhitespace = true
			case "th", "td":
				tb.QueueSpace('\t')
			case "div", "figure", "figcaption", "picture", "li", "dt", "dd",
				"header", "footer", "main", "section", "article", "aside", "nav", "address",
				"details", "summary", "dialog", "form", "fieldset",
				"caption", "thead", "tbody", "tfoot", "tr":
				tb.WriteNewline(1, true)
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			render(child, keepWhitespace)
		}
	}

	render(doc, false)
	return tb.String()
}

func textContent(n *html.Node) string {
	if child := n.FirstChild; child != nil && child.Type == html.TextNode && child.NextSibling == nil {
		return child.Data
	}
	var b strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.TextNode {
			continue
		}
		b.WriteString(child.Data)
	}
	return b.String()
}

func isDisplay(el *html.Node, val string) bool {
	for _, attr := range el.Attr {
		if attr.Key == "display" {
			return attr.Val == val
		}
	}
	return false
}

// Finds the deeply nested annotation element with desired encoding
func findAnnotation(parent *html.Node, mimeType string) *html.Node {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}
		if child.Data == "annotation" {
			for _, attr := range child.Attr {
				if attr.Key == "encoding" {
					if attr.Val == mimeType {
						return child
					}
					break
				}
			}
		}
		if found := findAnnotation(child, mimeType); found != nil {
			return found
		}
	}
	return nil
}

// Finds the original LaTeX expression stored in a "data-latex" attribute
func findLatex(parent *html.Node) string {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}
		for _, attr := range child.Attr {
			if attr.Key == "data-latex" {
				return attr.Val
			}
		}
		if found := findLatex(child); found != "" {
			return found
		}
	}
	return ""
}

// Detect a MathJax 2 script element.
func isMathjaxScript(el *html.Node) (bool, bool) {
	for _, attr := range el.Attr {
		if attr.Key == "type" {
			// <script type="math/tex; mode=display">
			mime := attr.Val
			isBlock := false
			if idx := strings.IndexByte(mime, ';'); idx >= 0 {
				isBlock = strings.Contains(mime[idx+1:], "mode=display")
				mime = mime[:idx]
			}
			return mime == "math/tex", isBlock
		}
	}
	return false, false
}

func isHiddenElement(el *html.Node) bool {
	for _, attr := range el.Attr {
		if attr.Key == "aria-hidden" {
			return attr.Val == "" || attr.Val == "true"
		}
	}
	return false
}

type innerTextBuilder struct {
	sb strings.Builder
	sp rune
	nl uint8
}

func (tb *innerTextBuilder) String() string {
	return tb.sb.String()
}

func (tb *innerTextBuilder) QueueSpace(c rune) {
	if tb.sp > 0 {
		return
	}
	tb.sp = c
}

func (tb *innerTextBuilder) WriteNewline(n uint8, collapse bool) {
	if collapse {
		if tb.nl >= n {
			return
		}
		n -= tb.nl
	}
	tb.nl += n
	if collapse && tb.sb.Len() == 0 {
		return
	}
	for ; n > 0; n-- {
		tb.sb.WriteByte('\n')
	}
}

func (tb *innerTextBuilder) WriteWord(w string) {
	if tb.sp > 0 && tb.nl == 0 {
		tb.sb.WriteRune(tb.sp)
	}
	tb.sb.WriteString(w)
	tb.nl = 0
	tb.sp = 0
}

func (tb *innerTextBuilder) WritePre(pre string) {
	if tb.sp > 0 && tb.nl == 0 {
		tb.sb.WriteRune(tb.sp)
	}
	tb.sb.WriteString(pre)
	tb.nl = 0
	tb.sp = 0
}
