package main

import (
	"github.com/vovabcps/go-telegram-allkeyshop-bot/allkeyshop"
)

type State struct {
	chatState map[int64]allkeyshop.Games
}

func NewState() *State {
	return &State{make(map[int64]allkeyshop.Games)}
}

func (s *State) Get(chatId int64) allkeyshop.Games {
	c := s.chatState[chatId]
	return c
}

func (s *State) Contains(chatId int64) bool {
	_, c := s.chatState[chatId]
	return c
}

func (s *State) Add(chatId int64) {
	s.chatState[chatId] = allkeyshop.Games{}
}

func (s *State) Remove(chatId int64) {
	delete(s.chatState, chatId)
}

func (s *State) Games(chatId int64) allkeyshop.Games {
	return s.Get(chatId)
}


func (s *State) SetGames(i int64, games allkeyshop.Games) {
	s.chatState[i] = games
}
