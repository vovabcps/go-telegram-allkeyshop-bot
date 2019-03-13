package main

import "github.com/vovabcps/go-telegram-allkeyshop-bot/allkeyshop"

type State struct {
	chatState map[int64]ChatState
}

func NewState() *State {
	return &State{make(map[int64]ChatState)}
}

func (s *State) Get(chatId int64) *ChatState {
	c := s.chatState[chatId]
	return &c
}

func (s *State) Contains(chatId int64) bool {
	_, c := s.chatState[chatId]
	return c
}

func (s *State) Add(chatId int64, games allkeyshop.Games) {
	s.chatState[chatId] = ChatState{games, 0}
}

func (s *State) Remove(chatId int64) {
	delete(s.chatState, chatId)
}

func (s *State) Games(chatId int64) allkeyshop.Games {
	return s.Get(chatId).Games
}

func (s *State) Offset(chatId int64) int64 {
	return s.Get(chatId).Offset
}

type ChatState struct {
	Games  allkeyshop.Games
	Offset int64
}
