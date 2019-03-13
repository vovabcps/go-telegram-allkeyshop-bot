package allkeyshop

import "github.com/gocolly/colly"

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

type Deals []Deal
