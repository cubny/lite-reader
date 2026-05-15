package readability

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"codeberg.org/readeck/go-readability/v2/render"
	"github.com/araddon/dateparse"
	"golang.org/x/net/html"
)

var ErrTimestampMissing = errors.New("timestamp not found in document")

// Article is the final readable content of a parsed article.
type Article struct {
	// Node is the top-level container of cleaned-up article content. It may be nil if there were
	// errors or if article content was blank.
	Node *html.Node

	title         string
	byline        string
	excerpt       string
	siteName      string
	image         string
	favicon       string
	language      string
	publishedTime string
	modifiedTime  string
}

func (a Article) RenderText(w io.Writer) error {
	if a.Node == nil {
		return fmt.Errorf("the Node field is nil")
	}
	text := render.InnerText(a.Node)
	_, err := fmt.Fprint(w, text)
	return err
}

func (a Article) RenderHTML(w io.Writer) error {
	if a.Node == nil {
		return fmt.Errorf("the Node field is nil")
	}
	return html.Render(w, a.Node)
}

func (a Article) Title() string {
	return a.title
}

func (a Article) Byline() string {
	return a.byline
}

func (a Article) Excerpt() string {
	excerpt := a.excerpt
	if excerpt == "" {
		if paragraph := getElementByTagName(a.Node, "p"); paragraph != nil {
			excerpt = strings.TrimSpace(render.InnerText(paragraph))
		}
	}
	return strings.Join(strings.Fields(excerpt), " ")
}

func (a Article) SiteName() string {
	return a.siteName
}

func (a Article) ImageURL() string {
	return a.image
}

func (a Article) Favicon() string {
	return a.favicon
}

func (a Article) Language() string {
	return a.language
}

// PublishedTime is the time when the article was published. If no timestamp was found in the
// article metadata, the error will be ErrTimestampMissing.
func (a Article) PublishedTime() (time.Time, error) {
	if a.publishedTime == "" {
		return time.Time{}, ErrTimestampMissing
	}
	t, err := dateparse.ParseAny(a.publishedTime)
	if err != nil {
		return t, fmt.Errorf("error parsing publishedTime: %w", err)
	}
	return t, nil
}

// ModifiedTime is the time when the article was modified. If no timestamp was found in the
// article metadata, the error will be ErrTimestampMissing.
func (a Article) ModifiedTime() (time.Time, error) {
	if a.modifiedTime == "" {
		return time.Time{}, ErrTimestampMissing
	}
	t, err := dateparse.ParseAny(a.modifiedTime)
	if err != nil {
		return t, fmt.Errorf("error parsing modifiedTime: %w", err)
	}
	return t, nil
}
