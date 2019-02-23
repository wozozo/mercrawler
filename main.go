package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/nlopes/slack"
)

type MercariItem struct {
	name  string
	price string
	image string
	url   string
}

type NotifyArgs struct {
	token       string
	channel     string
	mercariItem MercariItem
}

// Scrape mercari item list
func Scrape(keyword string) []MercariItem {
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

	var mercariItems []MercariItem

	doc.Find(".items-box-content .items-box").Each(func(i int, s *goquery.Selection) {
		itemName := s.Find(".items-box-name").Text()
		itemPrice := s.Find(".items-box-price").Text()
		itemImage, _ := s.Find("img").Attr("data-src")
		itemURL, _ := s.Find("a").Attr("href")

		mercariItems = append(mercariItems, MercariItem{
			name:  itemName,
			price: itemPrice,
			image: itemImage,
			url:   itemURL,
		})
	})

	return mercariItems
}

// Notify to Slack channel
func Notify(notifyArgs NotifyArgs) {
	api := slack.New(notifyArgs.token)
	mercariItem := notifyArgs.mercariItem

	attachment := slack.MsgOptionAttachments(slack.Attachment{
		Title:     mercariItem.price + ": " + mercariItem.name,
		ImageURL:  mercariItem.image,
		TitleLink: mercariItem.url,
	})

	_, respTimestamp, err := api.PostMessage(notifyArgs.channel, attachment)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Succeed: %s => %s", respTimestamp, mercariItem.name)
	}
}

func main() {
	var (
		keyword = flag.String("keyword", "", "Search item keyword")
	)
	flag.Parse()
	items := Scrape(*keyword)

	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")

	for i := 0; i < len(items); i++ {
		notifyArgs := NotifyArgs{
			token:       slackToken,
			channel:     slackChannel,
			mercariItem: items[i],
		}
		Notify(notifyArgs)
	}
}
