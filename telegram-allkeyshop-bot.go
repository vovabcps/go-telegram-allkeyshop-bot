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
	defer elapsed(update.Message.Chat.ID, text)()
	chatId := update.Message.Chat.ID

	switch {
	case !state.Contains(chatId):
		handleInitialRequest(bot, aks, &update, state)
	default:
		handleStatedRequest(bot, aks, &update, state)
	}
}

func handleInitialRequest(bot *tgbotapi.BotAPI, aks *allkeyshop.AksAPI, update *tgbotapi.Update, state *State) {
	var response string
	chatId := update.Message.Chat.ID

	switch {
	case update.Message.IsCommand():
		response = initialHelp
	default:
		sendMessage(bot, chatId, fetchingGamesMessage)
		games, e := findGames(aks, update.Message.Text)
		if e!= nil {
			response = internalErrorMessage
		} else {
			if len(games) <= 0 {
				response = noResultsMessage
			} else {
				state.Add(update.Message.Chat.ID, games)
				g := GetNextGames(games, 0, 10)
				response = FormatGames(g)
				response += "\n" + foundHelp
			}
		}
	}

	sendMessage(bot, chatId, response)
}

func handleStatedRequest(bot *tgbotapi.BotAPI, aks *allkeyshop.AksAPI, update *tgbotapi.Update, state *State) {
	var response string
	var input string
	chatId := update.Message.Chat.ID
	regex := regexp.MustCompile("^[0-9]*$")

	switch {
	case update.Message.IsCommand():
		input = update.Message.Command()
	default:
		input = update.Message.Text
	}

	switch {
	case regex.MatchString(input):
		index, _ := strconv.Atoi(input)
		game, e := getGameByNumber(state.Games(chatId), index)
		if e == nil {
			sendMessage(bot, chatId, fetchingDealsMessage)
			deals := getBestDeals(getDeals(aks, game), 5)
			response = FormatDeals(deals)
			state.Remove(chatId)
		} else {
			response = foundHelp
		}
	case input == "cancel":
		state.Remove(chatId)
		response = initialHelp
	default:
		response = foundHelp
	}

	sendMessage(bot, chatId, response)
}

func findGames(aks *allkeyshop.AksAPI, input string) (allkeyshop.Games, error) {
	games, e := aks.Find(input)
	return games, e
}

func getDeals(aks *allkeyshop.AksAPI, game allkeyshop.Game) allkeyshop.Deals {
	deals := aks.GetDeals(game)
	return deals
}

func sendMessage(bot *tgbotapi.BotAPI, chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
