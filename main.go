package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func Scrape(keyword string) {
	targetURL := "https://www.mercari.com/jp/search/?keyword=" + url.QueryEscape(keyword)
	res, err := http.Get(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".items-box-content .items-box").Each(func(i int, s *goquery.Selection) {
		itemName := s.Find(".items-box-name").Text()
		itemPrice := s.Find(".items-box-price").Text()
		itemImage, _ := s.Find("img").Attr("data-src")
		fmt.Printf("name => %s\n", itemName)
		fmt.Printf("price => %s\n", itemPrice)
		fmt.Printf("image => %s\n", itemImage)
	})
}

func main() {
	var (
		keyword = flag.String("keyword", "", "Search item keyword")
	)
	flag.Parse()
	Scrape(*keyword)
}
