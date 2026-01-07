package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func buildETPURL(page int) string {
	u, _ := url.Parse("https://new.etpgpb.ru/procedures/")
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("search", "аренда оборудования")
	q.Set("sort", "by_relevance")
	u.RawQuery = q.Encode()
	return u.String()
}

func ExampleScrape() {
	client := &http.Client{}

	for page := 1; page <= 10; page++ {
		link := buildETPURL(page)
		fmt.Println("PAGE:", page, link)

		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != 200 {
			res.Body.Close()
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("items:", doc.Find(".proceduresList__item").Length())

		doc.Find(".proceduresList__item").Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Find(".vTitle.vTitle--1").Text())

			priceSel := s.Find(".vTitle.vTitle--2.cardBody__price")
			price, ok := priceSel.Attr("content")
			if !ok {
				price = strings.TrimSpace(priceSel.Text())
			}

			href, _ := s.Find("a").Attr("href")
			fullURL := href
			if strings.HasPrefix(href, "/") {
				fullURL = "https://new.etpgpb.ru" + href
			}

			if title != "" {
				fmt.Println(title, price, fullURL)
			}
		})
	}
}

func main() {
	ExampleScrape()
}
