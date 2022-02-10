package mindgame

import (
	"fmt"
	"strconv"
)

type RoomInfo struct {
	owner        string
	participants []string
	status       string
	gameData     GameInfo
}

var maxLevel = 3

func (r *RoomInfo) StartRound(level int) {
	data := GameInfo{}
	data.lives = r.gameData.lives
	data.level = level
	var dealtCards []int
	for val := range selectN(len(r.participants) * level) {
		dealtCards = append(dealtCards, val)
	}
	for i, val := range r.participants {
		data.players = append(data.players, Player{val, dealtCards[i*level : (i+1)*level]})
	}
	r.gameData = data
}

func (r *RoomInfo) StartGame() {
	r.gameData.lives = len(r.participants)
	r.StartRound(1)
	fmt.Printf("%v\n", r.gameData.players)
}

func (r *RoomInfo) PlayCard(userId, card string) string {
	data := GameInfo{}
	data.lives = r.gameData.lives
	data.level = r.gameData.level
	data.players = r.gameData.players
	data.playedCards = r.gameData.playedCards
	correctCard := 999
	cardInt, err := strconv.Atoi(card)
	fmt.Printf("Game state pre card play: %v", r.gameData)
	if err != nil {
		panic("error converting card to int")
	}
	data.playedCards = append(data.playedCards, cardInt)
	for i, player := range data.players {
		var toRemove []int
		for j, card := range player.cards {
			if card < correctCard {
				correctCard = card
			}
			if card == cardInt {
				if player.id != userId {
					return "illegalMove"
				} else {
					toRemove = append(toRemove, j)
				}
			}
		}
		if len(toRemove) > 0 {
			data.players[i].cards = RemoveFromSlice(toRemove, player.cards)
		}
	}

	if cardInt == correctCard {
		hasPassedLevel := hasPassedLevel(r.gameData.players)
		if hasPassedLevel {
			if data.level+1 > maxLevel {
				return "gameWin"
			}
			r.StartRound(data.level + 1)
			// r.gameData = data
			return "gameNext"
		} else {
			r.gameData = data
			return "correct"
		}
	} else {
		for i, player := range r.gameData.players {
			var toRemove []int
			for j, card := range player.cards {
				if card <= cardInt {
					// remove the card
					println(j)
					toRemove = append(toRemove, j)
				}
			}
			if len(toRemove) > 0 {
				data.players[i].cards = RemoveFromSlice(toRemove, player.cards)
			}
		}
		data.lives--
		if data.lives <= 0 {
			r.gameData = data
			return "gameLose"
		} else {
			hasPassedLevel := hasPassedLevel(r.gameData.players)
			if hasPassedLevel {
				if data.level+1 > maxLevel {
					return "gameWin"
				}
				r.gameData = data
				fmt.Printf("Game Data Fail: %v\n", r.gameData)
				r.StartRound(data.level + 1)
				return "gameNext"
			}
			r.gameData = data
			return "incorrect"
		}
	}
}

func (r *RoomInfo) RemoveUser(userId string) string {
	for i, val := range r.participants {
		if val == userId {
			r.participants = RemoveFromStringSlice([]int{i}, r.participants)
			break
		}
	}
	if userId == r.owner {
		if len(r.participants) <= 0 {
			return "delete"
		} else {
			r.owner = r.participants[0]
		}
	}
	if r.status == "inProgress" {
		data := GameInfo{r.gameData.lives, r.gameData.level, r.gameData.players, r.gameData.playedCards}
		for i, val := range data.players {
			if val.id == userId {
				data.players = RemoveFromPlayerSlice([]int{i}, data.players)
				break
			}
		}
		r.gameData = data
		r.StartRound(data.level)
		return "newRound"
	}
	return "updateLobby"
}
