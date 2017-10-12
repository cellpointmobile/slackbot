package event

import (
	"fmt"
	"net/http"
	"os"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
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

func Get_random_quote() string {
	doc, err := goquery.NewDocument("https://www.brainyquote.com/quotes/keywords/coffee.html")
	if err != nil {
		log.Fatalln(err)
	}

	quotes := []string{}

	doc.Find(".b-qt").Each(func(i int, sel *goquery.Selection) {
		quotes = append(quotes, sel.Text())
	})

	n := rand.Int() % len(quotes)
	return quotes[n]
}
