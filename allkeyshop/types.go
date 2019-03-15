package allkeyshop

import (
	"github.com/gocolly/colly"
)

type Game struct {
	Title   string
	Details string
	Url     string
	Price   string
}

type Deal struct {
	Shop     string
	Platform string
	Price    string
	Url      string
}

type AksAPI struct {
	collector *colly.Collector
}

type Games []Game

type IndexError struct {
	Msg string
}

func (e *IndexError) Error() string {
	return e.Msg
}

func (games Games) Get(n int) (Game, error) {
	if len(games) < n+1 {
		return Game{}, &IndexError{}
	}
	return games[n], nil
}

type Deals []Deal
