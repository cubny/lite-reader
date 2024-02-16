# Feedfinder

Golang port of [this](https://github.com/dfm/feedfinder2) python library (not actively maintained).


It finds links to rss/atom/rdf feeds in a website. 
Better support for Twitter and reddit links and option to add more extensions

Wrote this to use in rewrite of this [bot](https://github.com/nikhil1raghav/rssbot) in golang.

```go
package main
import (
	"fmt"

	"github.com/nikhil1raghav/feedfinder"
)
func main() {

	f := feedfinder.NewFeedFinder()
    url:="old.reddit.com/r/unixporn"
	links, _ := f.FindFeeds(url)
	for _, link := range links {
		fmt.Println(link)
	}
}
```

Feedfinder supports following options
- `CheckAll` : if set to `true` initiates exhaustive search for feedurls, default `false`
- `UserAgent`: Custom user-agent string for the crawler, default `chrome user agent on windows`
- `Timeout` : timeout duration for the feedfinding operation (not supported yet but option is there), default `60 seconds`

```go
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/nikhil1raghav/feedfinder"
	"github.com/nikhil1raghav/feedfinder/values"
)

func main() {

	f := feedfinder.NewFeedFinder(
		feedfinder.UserAgent(values.ChromeUserAgent), //custom user agent
		feedfinder.CheckAll(true), //to check if url not found easily
		feedfinder.TimeOut(10*time.Second), //timeout to avoid long waiting times
	)
	url := flag.String("URL", "https://raghavnikhil.com/", "url to find feeds for")
	flag.Parse()
	links, _ := f.FindFeeds(*url)
	for _, link := range links {
		fmt.Println(link)
	}
}
```





