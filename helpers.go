package main

import (
	"fmt"
	"github.com/vovabcps/go-telegram-allkeyshop-bot/allkeyshop"
	"time"
)

func FormatGames(games allkeyshop.Games) string {
	result := ""
	for i, game := range games {
		result += fmt.Sprintf("/%d %s %s %s\n\n", i+1, game.Title, game.Details, game.Price)
	}
	return result
}

func FormatDeals(deals allkeyshop.Deals) string {
	result := ""
	for _, deal := range deals {
		result += fmt.Sprintf("%s %s€ %s\n\n", deal.Shop, deal.Price, deal.Url)
	}
	return result
}

func GetNextGames(games allkeyshop.Games, offset int, number int) allkeyshop.Games {
	var result allkeyshop.Games
	l := len(games)
	if offset >= l {
		result = nil
	}
	if l-offset < number {
		result = games[offset:]
	} else {
		result = games[offset : offset+number]
	}
	return result
}

func getBestDeals(deals allkeyshop.Deals, n int) allkeyshop.Deals {
	var result allkeyshop.Deals
	l := len(deals)

	if l-n < 0 {
		result = deals[0:]
	} else {
		result = deals[0:n]
	}
	return result
}

type IndexError struct {
	Msg string
}

func (e *IndexError) Error() string {
	return e.Msg
}

func getGameByNumber(games allkeyshop.Games, n int) (allkeyshop.Game, error) {
	if len(games) < n {
		 return allkeyshop.Game{}, &IndexError{}
	}
	return games[n-1], nil
}

func elapsed(chat int64, msg string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s from %d took %v\n", msg, chat, time.Since(start))
	}
}
