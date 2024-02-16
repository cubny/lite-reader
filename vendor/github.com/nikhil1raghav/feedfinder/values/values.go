package values

var FeedHeaders = []string{"<rdf", "<rss", "<feed"}
var FeedUrlSuffix = []string{".rss", ".atom", ".xml", ".rdf"}
var FeedLike = []string{"rss", "feed", "atom", "rdf", "xml"}

var GuessWords = []string{
	"atom.xml",
	"index.xml",
	"index.atom",
	"index.rdf",
	"rss.xml",
	"index.rss",
}

var FeedTypes = []string{
	"application/rss+xml",
	"text/xml",
	"application/atom+xml",
	"application/x.atom+xml",
	"application/x-atom+xml",
}

const ChromeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.61 Safari/537.36"
const FirefoxUserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:100.0) Gecko/20100101 Firefox/100.0"
