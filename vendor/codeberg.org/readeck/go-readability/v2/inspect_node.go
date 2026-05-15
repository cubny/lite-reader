package readability

import (
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/net/html"
)

// inspectNode wraps a HTML node to use in structured logging
func inspectNode(node *html.Node) slog.LogValuer {
	return &inspectedNode{node}
}

type inspectedNode struct {
	node *html.Node
}

func (n *inspectedNode) LogValue() slog.Value {
	if n.node.Type == html.TextNode {
		return slog.StringValue(n.node.Data)
	}

	var tagPreview strings.Builder
	tagPreview.WriteString("<")
	tagPreview.WriteString(n.node.Data)

	hasOtherAttributes := false
	for _, attr := range n.node.Attr {
		switch strings.ToLower(attr.Key) {
		case "id", "class", "rel", "itemprop", "name", "type", "role", "for", "action", "method":
			fmt.Fprintf(&tagPreview, ` %s=%q`, attr.Key, attr.Val)
		case "src", "href":
			val := attr.Val
			if strings.HasPrefix(val, "data:") {
				if v, _, ok := strings.Cut(val, ","); ok {
					val = v + ",***"
				}
			} else if strings.HasPrefix(val, "javascript:") {
				val = "javascript:***"
			}
			fmt.Fprintf(&tagPreview, ` %s=%q`, attr.Key, val)
		default:
			if !strings.HasPrefix(attr.Key, "data-readability-") {
				hasOtherAttributes = true
			}
		}
	}
	if hasOtherAttributes {
		tagPreview.WriteString(" ...")
	}

	if n.node.FirstChild == nil {
		tagPreview.WriteString("/")
	}
	tagPreview.WriteString(">")

	if c := n.node.FirstChild; c != nil && c.Type == html.TextNode && hasContent(c.Data) {
		text := []rune(strings.TrimSpace(c.Data))
		if len(text) > 15 {
			tagPreview.WriteString(string(text[:12]))
			tagPreview.WriteString("...")
		} else {
			tagPreview.WriteString(string(text))
		}
	}

	return slog.StringValue(tagPreview.String())
}

func inspectXPath(node *html.Node) slog.LogValuer {
	return &xpathNode{node}
}

type xpathNode struct {
	*html.Node
}

func (n *xpathNode) LogValue() slog.Value {
	return slog.StringValue(getXPathSelector(n.Node))
}

// getXPathSelector constructs an XPath selector that uniquely identifies a DOM node.
func getXPathSelector(node *html.Node) string {
	p := node
	if node.Type == html.TextNode {
		p = node.Parent
	}
	var names []string

	for p.Parent != nil {
		elementPos := 1
		for s := p.PrevSibling; s != nil; s = s.PrevSibling {
			if s.Type == html.ElementNode && s.Data == p.Data {
				elementPos++
			}
		}
		names = append(names, "")
		copy(names[1:], names)
		names[0] = fmt.Sprintf("%s[%d]", p.Data, elementPos)
		p = p.Parent
	}

	if node.Type == html.TextNode {
		textNodePos := 1
		for s := node.PrevSibling; s != nil; s = s.PrevSibling {
			if s.Type == html.TextNode {
				textNodePos++
			}
		}
		names = append(names, fmt.Sprintf("text()[%d]", textNodePos))
	}

	return strings.Join(names, "/")
}
