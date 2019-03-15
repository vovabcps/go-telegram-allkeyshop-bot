package allkeyshop

import (
	"github.com/gocolly/colly"
	s "strings"
)

const rootUrl = "https://www.allkeyshop.com/"


func NewAksAPI(options ...func(api *AksAPI)) *AksAPI {

	api := &AksAPI{colly.NewCollector(
		colly.AllowedDomains("www.allkeyshop.com"),
		colly.AllowURLRevisit(),
	)}
	for _, f := range options {
		f(api)
	}
	return api
}

func Cache(dir string) func(*AksAPI) {
	return func(api *AksAPI) {
		colly.CacheDir(dir)(api.collector)
	}
}

//Uses site search and returns slice of Game
func (api *AksAPI) Find(input string) (Games, error) {
	var games Games
	api.collector.OnHTML(".search-results-row-link", func(e *colly.HTMLElement) {
		temp := Game{}
		child := e.DOM.ChildrenFiltered(".search-results-row-game")
		temp.Title = child.ChildrenFiltered(".search-results-row-game-title").Text()
		temp.Details = child.ChildrenFiltered(".search-results-row-game-infos").Text()
		temp.Price = e.ChildText(".search-results-row-price")
		temp.Url = e.Attr("href")
		games = append(games, temp)
	})
	err := api.collector.Visit(rootUrl + "blog/catalogue/search-" + input)
	if err != nil {
		return Games{}, err
	}
	return games, nil
}

func (api *AksAPI) GetDeals(game Game) Deals {
	var deals []Deal
	api.collector.OnHTML(".offers-table-row", func(e *colly.HTMLElement) {
		temp := Deal{}
		temp.Shop, _ = e.DOM.Find(".offers-merchant").Attr("title")
		temp.Price = e.DOM.Find("[data-offer-price-container]").Text()
		temp.Platform = s.TrimSpace(s.Split(e.DOM.Find(".offers-merchant-reviews .d-block.d-xl-none").Text(), "\n")[2])
		url, _ := e.DOM.Find(".buy-btn").Attr("href")
		temp.Url = s.TrimPrefix(url, "//")
		deals = append(deals, temp)
	})
	if e := api.collector.Visit(game.Url); e != nil {
		panic(e)
	}
	return deals
}
