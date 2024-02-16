package clients

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/nikhil1raghav/feedfinder/values"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Crawler struct {
	Client *http.Client
	UserAgent string
	Timeout time.Duration
}
func NewCrawler(useragent string, timeout time.Duration) *Crawler{
	c:=&Crawler{}
	c.Timeout = timeout
	c.UserAgent = useragent
	client:=&http.Client{
		Timeout: c.Timeout,
	}
	c.Client = client
	return c
}
func (c *Crawler) NewRequest(u string) (*http.Request, error){
	req, err:=http.NewRequest(http.MethodGet, u, nil)
	if err!=nil{
		return nil, err
	}
	req.Header.Add("User-Agent", c.UserAgent)
	return req,nil
}

func (c *Crawler) Get(u string) (string, error){
	req, err:=c.NewRequest(u)
	if err!=nil{
		return "", err
	}
	resp, err:=c.Client.Do(req)
	if err!=nil{
		return "",err
	}
	defer resp.Body.Close()
	body, err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return "",err
	}
	return string(body), err
}

//Parse the page and look for urls that end in a standard feed extension
func (c *Crawler) TypePass(doc *goquery.Document, u string) []string {
	links := make([]string, 0)
	base, err := url.Parse(u)
	if err != nil {
		log.Println("Error parsing url", err)
		return links
	}
	doc.Find("link").Each(func(i int, e *goquery.Selection){
		linkType,exist:=e.Attr("type")
		if !exist{
			return
		}
		for _, feedType:=range values.FeedTypes{
			if linkType==feedType{
				href,exists:=e.Attr("href")
				if !exists{
					continue
				}
				feedUrl, _:=url.Parse(href)
				links=append(links, base.ResolveReference(feedUrl).String())
			}
		}
	})
	log.Printf("Found %d feed links in typePass", len(links))
	return links
}

//Get all <a> tags in a page
func (c *Crawler) GetAllAnchors(doc *goquery.Document, u string) []string {

	anchors := make([]string, 0)
	base, err := url.Parse(u)
	if err != nil {
		log.Println("Error parsing url", err)
		return []string{}
	}
	log.Println("Getting all anchors")
	doc.Find("a").Each(func(i int, e *goquery.Selection){
		href, exists:=e.Attr("href")
		if !exists{
			return
		}
		if strings.Count(href,"://")>0{
			anchors=append(anchors, href)
		}else{
			localUrl,_:=url.Parse(href)
			anchors=append(anchors, base.ResolveReference(localUrl).String())
		}
	})
	return anchors
}
