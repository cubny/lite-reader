package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmcdole/gofeed"
)

func main() {
	router := httprouter.New()

	http.Get("/", func() {
		http.FileServer(http.Dir("./public"))
	})
	http.Get("/addFeed/:id", addFeed)
	http.ListenAndServe(":3000", nil)
}

func addFeed(w http.ResponseWriter, req *http.Request) {
	fp := gofeed.NewParser()
	id := req.URL.Query().Get("id")
	feed, _ := fp.ParseURL(id)
	fmt.Fprintf(w, feed.Title)
}
