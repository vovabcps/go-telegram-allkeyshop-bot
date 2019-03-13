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

func GetNextGames(games allkeyshop.Games, offset int, number int) (allkeyshop.Games, error) {
	l := len(games)
	if offset >= l {
		return nil, OutOfBondsError{offset, len(games)}
	}
	var result allkeyshop.Games
	if l-offset < number {
		result = games[offset:]
	} else {
		result = games[offset : offset+number]
	}
	return result, nil
}

func elapsed() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}
