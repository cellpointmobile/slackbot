package event

import (
	"fmt"
	"net/http"
	"os"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"time"
)

const IFTTT_BASE_URL = "https://maker.ifttt.com/trigger"

func Power_on() {
	fmt.Println("power on event fired")
	http.Get(IFTTT_BASE_URL + "/power_on/with/key/" + os.Getenv("IFTTT_TOKEN") )
}

func Power_off() {
	fmt.Println("power off event fired")
	http.Get(IFTTT_BASE_URL + "/power_off/with/key/" + os.Getenv("IFTTT_TOKEN") )
}

func Get_random_quote() (author string, quote string) {
	urls := []string{
		"https://www.brainyquote.com/quotes/keywords/coffee.html",
		"https://www.brainyquote.com/quotes/keywords/coffee_2.html?vm=l",
		"https://www.brainyquote.com/quotes/keywords/coffee_3.html?vm=l",
		"https://www.brainyquote.com/quotes/keywords/coffee_4.html?vm=l",
	}

	quotes := map[string]string{}
	authors := []string{}

	for _, url := range urls {
		doc, err := goquery.NewDocument(url)

		if err != nil {
			log.Fatalln(err)
		}

		doc.Find(".bq_list_i").Each(func(i int, sel *goquery.Selection) {
			quote := sel.Find(".b-qt").Text()
			author := sel.Find(".bq-aut").Text()
			authors = append(authors, author)
			quotes[author] = quote
		})
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(authors)
	return authors[n], quotes[authors[n]]
}
