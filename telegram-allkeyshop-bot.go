package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vovabcps/go-telegram-allkeyshop-bot/allkeyshop"
	"log"
	"regexp"
	"strconv"
)

func main() {
	state := NewState()
	aks := allkeyshop.NewAksAPI()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening for updates...")
	u := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		go handleUpdate(update, bot, aks, state)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, aks *allkeyshop.AksAPI, state *State) {
	if update.Message == nil {
		return
	}
	chat := update.Message.Chat
	text := update.Message.Text
	log.Printf("Message from %s %d: %s", chat.UserName, chat.ID, text)
	defer elapsed()()
	var response string

	chatId := update.Message.Chat.ID

	switch {
	case !state.Contains(chatId):
		response = handleInitial(aks, &update, state)
	default:
		response = handleStated(aks, &update, state)
	}

	msg := tgbotapi.NewMessage(chatId, response)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err.Error())
	}
}

func handleInitial(aks *allkeyshop.AksAPI, update *tgbotapi.Update, state *State) string {
	var response string
	if update.Message.IsCommand() {
		return initialHelp
	}
	games := findGames(aks, update.Message.Text)
	if len(games) <= 0 {
		return "No results."
	}
	state.Add(update.Message.Chat.ID, games)
	g, e := GetNextGames(games, 0, 10)
	if e != nil {
		log.Println(e.Error())
		response = "Internal error."
	} else {
		response = FormatGames(g)
	}
	response += "\n" + foundHelp
	return response
}

func handleStated(aks *allkeyshop.AksAPI, update *tgbotapi.Update, state *State) string {
	var response string
	chatId := update.Message.Chat.ID
	if !update.Message.IsCommand() {
		return foundHelp
	}
	command := update.Message.Command()
	if command == "cancel" {
		state.Remove(chatId)
		return initialHelp
	}
	match, _ := regexp.MatchString("^[0-9]*$", command)
	if !match {
		return foundHelp
	}
	input, _ := strconv.Atoi(command)
	game := state.Games(chatId)[input-1]
	deals := getDeals(aks, game)
	response = FormatDeals(deals)
	state.Remove(chatId)
	return response
}

func findGames(aks *allkeyshop.AksAPI, input string) allkeyshop.Games {
	games := aks.Find(input)
	return games
}

func getDeals(aks *allkeyshop.AksAPI, game allkeyshop.Game) allkeyshop.Deals {
	deals := aks.GetDeals(game)
	return deals
}
