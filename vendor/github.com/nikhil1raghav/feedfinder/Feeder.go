package feedfinder

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"

	"github.com/nikhil1raghav/feedfinder/clients"
	"github.com/nikhil1raghav/feedfinder/utils"
	"github.com/nikhil1raghav/feedfinder/values"
)

type FeedFinder struct {
	UserAgent  string
	TimeOut    time.Duration
	CheckAll   bool
	Crawler    clients.Crawler
}

func NewFeedFinder(options ...func(*FeedFinder)) *FeedFinder {
	f := &FeedFinder{}
	f.Init()
	for _, fn := range options {
		fn(f)
	}

	return f
}
func UserAgent(ua string) func(*FeedFinder) {
	return func(f *FeedFinder) {
		f.UserAgent = ua
		f.Crawler.UserAgent = ua
	}
}
func CheckAll(checkall bool) func(*FeedFinder) {
	return func(f *FeedFinder) {
		f.CheckAll = checkall
	}
}
func TimeOut(timeout time.Duration) func(*FeedFinder) {
	return func(f *FeedFinder) {
		f.TimeOut = timeout
		f.Crawler.Timeout = timeout
	}
}
func (f *FeedFinder) Init() {
	f.UserAgent = values.ChromeUserAgent
	f.CheckAll = false
	f.TimeOut = 60*time.Second
	f.Crawler = *clients.NewCrawler(f.UserAgent, f.TimeOut)
}
func (f *FeedFinder) getFeed(url string) (string, error) {
	resp, err := f.Crawler.Get(url)
	if err != nil {
		log.Printf("Error getting %s, %s", url, err.Error())
		return "", err
	}
	return resp,nil
}

func (f *FeedFinder) isFeedData(data string) bool {
	dataString := strings.ToLower(data)
	if strings.Count(dataString, "<html") > 0 {
		return false
	}
	for _, header := range values.FeedHeaders {
		if strings.Count(dataString, header) > 0 {
			return true
		}
	}
	return false
}
func (f *FeedFinder) isFeedUrl(url string) bool {
	url = strings.ToLower(url)
	for _, suffix := range values.FeedUrlSuffix {
		if strings.HasSuffix(url, suffix) {
			return true
		}
	}
	return false
}
func (f *FeedFinder) isFeedLike(url string) bool {
	url = strings.ToLower(url)
	for _, word := range values.FeedLike {
		if strings.Count(url, word) > 0 {
			return true
		}
	}
	return false
}

//called when all else failed
//validate feed after guessing
func (f *FeedFinder) guessUrls(u string) []string {
	guessed := make([]string, 0)
	for _, suffix := range values.GuessWords {
		url, err := utils.JoinUrl(u, suffix)
		if err != nil {
			log.Println(err)
			continue
		} else if validFeed, _ := f.isFeed(url); validFeed {

			guessed = append(guessed, url)
		}
	}
	return guessed

}
func (f *FeedFinder) isFeed(u string) (bool, error) {
	data, err := f.getFeed(u)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return f.isFeedData(data), nil
}
func (f *FeedFinder) FindFeeds(url string) ([]string, error) {
	url = utils.ForceUrl(url)
	feedUrls := make([]string, 0)
	if validFeed, _ := f.isFeed(url); validFeed {
		feedUrls = append(feedUrls, url)
	}
	data, err:= f.Crawler.Get(url)
	if err!=nil{
		log.Println("Couldn't parse the page ",err)
		return []string{}, err
	}
	htmlDoc, err:=goquery.NewDocumentFromReader(strings.NewReader(data))
	if err!=nil{
		log.Println("Error converting to goquery doc", err)
		return []string{},err
	}

	feedUrls = append(feedUrls, f.Crawler.TypePass(htmlDoc,url)...)

	if len(feedUrls) > 0 && !f.CheckAll {
		return feedUrls, nil
	}

	anchors := f.Crawler.GetAllAnchors(htmlDoc, url)
	filteredUrl := make([]string, 0)
	for _, anchor := range anchors {
		if f.isFeedUrl(anchor) {
			filteredUrl = append(filteredUrl, anchor)
		} else if f.isFeedLike(anchor) {
			filteredUrl = append(filteredUrl, anchor)
		}
	}

	for _, u := range filteredUrl {
		if validFeed, _ := f.isFeed(u); validFeed {
			feedUrls = append(feedUrls, u)
		}
	}

	if len(feedUrls) > 0 && !f.CheckAll {
		return feedUrls, nil
	}
	feedUrls = append(feedUrls, f.guessUrls(url)...)

	return feedUrls, nil
}
