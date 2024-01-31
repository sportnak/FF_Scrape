package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func main() {
	// Create a new collector
	c := colly.NewCollector()
	athletes := make(map[string]*Athlete)

	// maps := make(map[string]map[string]Athlete)

	// Set up callbacks to be executed when an element is found
	c.OnHTML("div.Boxscore__Team", func(e *colly.HTMLElement) {
		// local_athletes := make(map[string]Athlete)

		title := e.DOM.Find(".TeamTitle__Name").Text()
		fmt.Println(title)
		is_passing := strings.Contains(title, "Passing")
		is_rushing := strings.Contains(title, "Rushing")
		is_receiving := strings.Contains(title, "Receiving")

		playerNames := e.DOM.Find("a.Boxscore__Athlete_Name")
		var players []string
		playerNames.Each(func(index int, element *goquery.Selection) {
			players = append(players, element.Text())
		})

		if len(players) == 0 {
			return
		}

		headerNames := e.DOM.Find("div.Table__Scroller th")
		var headers []string
		headerNames.Each(func(index int, element *goquery.Selection) {
			header := element.Text()
			headers = append(headers, header)
		})

		if len(headers) == 0 {
			return
		}

		playerStats := e.DOM.Find("div.Table__Scroller tbody tr")
		playerStats.Each(func(index int, element *goquery.Selection) {
			values := element.Find("td")

			if index > len(players)-1 {
				return
			}
			player := players[index]
			athlete, exists := athletes[player]

			if !exists {
				athletes[player] = &Athlete{
					Name: player,
				}
				athlete = athletes[player]
			}

			values.Each(func(statIndex int, valueElement *goquery.Selection) {
				key := headers[statIndex]
				data := Data{
					Key:   key,
					Value: valueElement.Text(),
					RUSH:  is_rushing,
					PASS:  is_passing,
					REC:   is_receiving,
				}

				athlete.SetData(data)
			})
		})

		// fmt.Println("Players", players, headers)
	})

	// c.OnHTML("a.Boxscore__Athlete_Name", func(e *colly.HTMLElement) {
	// 	playerName := e.Text
	// 	player := athletes[playerName]
	// 	if (player == nil) {
	// 		player = Athelete(
	// 			Name: playerName
	// 		)
	// 	}
	// 	if ()
	// })

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		// Do any additional processing here
		fmt.Println("Scraping complete for:", r.Request.URL.String())

		InitSheet()
		for key := range athletes {
			UpdatePlayer(athletes[key])
		}
	})

	// Specify the URL to scrape
	url := "https://www.espn.com/nfl/boxscore/_/gameId/401547380"

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

}
