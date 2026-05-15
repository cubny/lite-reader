package readability

import (
	"fmt"
	"io"
	nurl "net/url"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

// Parse parses a reader and find the main readable content.
func (ps *Parser) Parse(input io.Reader, pageURL *nurl.URL) (Article, error) {
	// Parse input
	doc, err := dom.Parse(input)
	if err != nil {
		return Article{}, fmt.Errorf("failed to parse input: %v", err)
	}

	return ps.ParseAndMutate(doc, pageURL)
}

// ParseDocument parses the specified document and find the main readable content.
func (ps *Parser) ParseDocument(doc *html.Node, pageURL *nurl.URL) (Article, error) {
	// Clone document to make sure the original kept untouched
	return ps.ParseAndMutate(dom.Clone(doc, true), pageURL)
}

// ParseAndMutate is like ParseDocument, but mutates doc during parsing.
func (ps *Parser) ParseAndMutate(doc *html.Node, pageURL *nurl.URL) (Article, error) {
	ps.doc = doc

	// Reset parser data
	ps.articleTitle = ""
	ps.articleByline = ""
	ps.articleDir = ""
	ps.articleSiteName = ""
	ps.documentURI = pageURL
	ps.attempts = []parseAttempt{}
	// These flags could get modified during subsequent passes in grabArticle
	ps.flags = flags{
		stripUnlikelys:     true,
		useWeightClasses:   true,
		cleanConditionally: true,
	}

	// Avoid parsing too large documents, as per configuration option
	if ps.MaxElemsToParse > 0 {
		numTags := len(dom.GetElementsByTagName(ps.doc, "*"))
		if numTags > ps.MaxElemsToParse {
			return Article{}, fmt.Errorf("documents too large: %d elements", numTags)
		}
	}

	// Unwrap image from noscript
	ps.unwrapNoscriptImages(ps.doc)

	// Extract JSON-LD metadata before removing scripts
	var jsonLd map[string]string
	if !ps.DisableJSONLD {
		jsonLd, _ = ps.getJSONLD()
	}

	// Remove script tags from the document.
	ps.removeScripts(ps.doc)

	// Prepares the HTML document
	ps.prepDocument()

	// Fetch metadata
	metadata := ps.getArticleMetadata(jsonLd)
	ps.articleTitle = metadata["title"]
	ps.articleByline = metadata["byline"]

	// Try to grab article content
	articleContent := ps.grabArticle()
	var readableNode *html.Node

	if articleContent != nil {
		ps.postProcessContent(articleContent)
		readableNode = dom.FirstElementChild(articleContent)
	}

	return Article{
		title:         ps.articleTitle,
		byline:        ps.articleByline,
		Node:          readableNode,
		excerpt:       metadata["excerpt"],
		siteName:      metadata["siteName"],
		image:         metadata["image"],
		favicon:       metadata["favicon"],
		language:      ps.articleLang,
		publishedTime: metadata["publishedTime"],
		modifiedTime:  metadata["modifiedTime"],
	}, nil
}
