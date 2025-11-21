package api

import (
	"fmt"
	"html"
	"strings"

	"github.com/cubny/lite-reader/internal/app/feed"
)

// renderFeedRow generates an HTML fragment for a single feed row
func renderFeedRow(f *feed.Feed) string {
	return fmt.Sprintf(`<li id="%d" class="feed" data-feed-id="%d">
	<img src="https://t1.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=%s" alt="">
	<div class="feedtitle">%s</div>
	<div class="count"><span>%s</span></div>
</li>`,
		f.ID,
		f.ID,
		html.EscapeString(f.Link),
		html.EscapeString(f.Title),
		formatUnreadCount(f.UnreadCount),
	)
}

// renderFeedList generates an HTML fragment for the complete feed list
func renderFeedList(feeds []*feed.Feed) string {
	var builder strings.Builder

	// Always include the Unread and Starred static feeds first
	builder.WriteString(`<li id="unread" class="feed">
	<div class="count"><span></span></div>
	<i class="icon-circle"></i>
	<div class="feedtitle">Unread</div>
</li>
<li id="starred" class="feed">
	<div class="count"><span></span></div>
	<i class="icon-star"></i>
	<div class="feedtitle">Starred</div>
</li>
`)

	if len(feeds) == 0 {
		// No user feeds, but we still have unread/starred
		return builder.String()
	}

	// Add user feeds
	for _, f := range feeds {
		builder.WriteString(renderFeedRow(f))
		builder.WriteString("\n")
	}
	return builder.String()
}

// renderFeedError returns an error message fragment
func renderFeedError(msg string) string {
	errorStyle := "padding: 10px; color: #d9534f; background: #f8d7da; " +
		"border: 1px solid #f5c6cb; border-radius: 4px; margin: 10px;"
	return fmt.Sprintf(`<div class="feed-error" style="%s">
	<i class="icon-exclamation-sign"></i> %s
</div>`, errorStyle, html.EscapeString(msg))
}

// formatUnreadCount formats the unread count for display
func formatUnreadCount(count int) string {
	if count > 0 {
		return fmt.Sprintf("%d", count)
	}
	return ""
}
