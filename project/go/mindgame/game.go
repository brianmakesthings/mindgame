package mindgame

import (
	"fmt"
	"math/big"
	"crypto/rand"
)

type Player struct {
	id string
	cards []int
}

type GameInfo struct {
	lives int
	level int
	players []Player
	playedCards []int
}

var deckSize int64= 100

func selectN (n int) map[int]bool {
	results := map[int]bool{}
	for len(results) <= n {
		randBig, err := rand.Int(rand.Reader, big.NewInt(deckSize))
		rand := int(randBig.Int64())+1
		if err != nil {
			panic("could not get random number")
		}
		if _, exists := results[rand]; !exists {
			results[rand] = true
		}
	}
	return results
}

func hasPassedLevel( players []Player) bool {
	hasPassedLevel := true
	for _, player := range players {
		if len(player.cards) != 0 {
			hasPassedLevel = false
			break;
		}
	}
	return hasPassedLevel
}

func GenerateUserId(request map[string]string, userIds *map[string]bool) map[string]string {
	name := request["string"]
	// will probably fail on very large user base
	randInt, err := rand.Int(rand.Reader, big.NewInt(4294967296))
	if err != nil {
		panic("Couldn't generate random number")
	}
	userId := fmt.Sprintf("%s%v", name, randInt)
	_, exists := (*userIds)[userId];
	for exists {
		randInt, err = rand.Int(rand.Reader, big.NewInt(4294967296))
		if err != nil {
			panic("Couldn't generate random number")
		}
		userId = fmt.Sprintf("%s%v", name, randInt)
		_, exists = (*userIds)[userId];
	}
	(*userIds)[userId] = true
	return map[string]string{"id" : userId}
}

func CreateRoom(request map[string]string, rooms *map[string]*RoomInfo) map[string]string {
	ownerId := request["ownerId"]
	randInt, err := rand.Int(rand.Reader, big.NewInt(4294967296))
	if err != nil {
		panic("Couldn't generate random number")
	}
	id := fmt.Sprintf("%v", randInt)
	_, exists := (*rooms)[id];
	for exists {
		randInt, err = rand.Int(rand.Reader, big.NewInt(4294967296))
		if err != nil {
			panic("Couldn't generate random number")
		}
		id = fmt.Sprintf("%v", randInt)
		_, exists = (*rooms)[id];
	}
	(*rooms)[id] = &RoomInfo{ownerId, []string{ownerId}, "lobby", GameInfo{}}
	return map[string]string{"roomId" : id}
}

func JoinRoom(request map[string]string, rooms *map[string]*RoomInfo) map[string]string {
	roomId := request["roomId"]
	userId := request["userId"]
	roomStat, exists := (*rooms)[roomId]
	if !exists {
		return map[string]string{"error" : "no such room"}
	}
	if len(roomStat.participants) >= 4 {
		return map[string]string{"error" : "room full"}
	}
	for _, participant := range roomStat.participants {
		if participant == userId {
			return map[string]string{"error" : "user already in room"}
		}
	}

	(*rooms)[roomId] = &RoomInfo{roomStat.owner, append(roomStat.participants, userId), roomStat.status, GameInfo{}}
	return map[string]string{"success" : "success"}
}

func GetParticipants(request map[string]string, rooms *map[string]*RoomInfo) map[string]string {
	roomId := request["roomId"]
	roomStat, exists := (*rooms)[roomId]
	if !exists {
		return map[string]string{"error" : "no such room"}
	}
	return map[string]string{ "participants" : fmt.Sprintf("%v", roomStat.participants), "owner" : roomStat.owner}
}

func StartGame(request map[string]string, rooms *map[string]*RoomInfo) map[string]string {
	roomId := request["roomId"]
	roomStat, exists := (*rooms)[roomId]
	if !exists {
		return map[string]string{"error" : "no such room"}
	}
	if len((*rooms)[roomId].participants) < 2 {
		return map[string]string{"error" : "Not Enough Players"}
	} else if len((*rooms)[roomId].participants) > 4 {
		return map[string]string{"error" : "Too Many Players"}
	}
	(*rooms)[roomId] = &RoomInfo{roomStat.owner, roomStat.participants, "inProgress", GameInfo{}}
	(*rooms)[roomId].StartGame()
	fmt.Printf("%v\n", (*rooms)[roomId].gameData)
	return map[string]string{"success" : "success"}
}

func GetGameState(request map[string]string, rooms *map[string]*RoomInfo) map[string]string{
	roomId := request["roomId"]
	_, exists := (*rooms)[roomId]
	if !exists {
		return map[string]string{"error" : "no such room"}
	}
	lives := fmt.Sprintf("%v", (*rooms)[roomId].gameData.lives)
	level := fmt.Sprintf("%v", (*rooms)[roomId].gameData.level)
	playerState := fmt.Sprintf("%v", (*rooms)[roomId].gameData.players)
	playedCards := fmt.Sprintf("%v", (*rooms)[roomId].gameData.playedCards)
	return map[string]string{ "lives" : lives, "level" : level, "playerState" : playerState, "playedCards" : playedCards}
}

func PlayCard(request map[string]string, rooms *map[string]*RoomInfo) map[string]string {
	roomId := request["roomId"]
	card := request["card"]
	userId := request["userId"]
	_, exists := (*rooms)[roomId]
	if !exists {
		return map[string]string{"error" : "no such room"}
	}
	gameResult := (*rooms)[roomId].PlayCard(userId, card)
	fmt.Printf("%v\n", (*rooms)[roomId].gameData)
	println(gameResult)
	gameState := GetGameState(request, rooms)
	gameState["moveResult"] = gameResult
	switch gameResult {
	case "correct":
	case "incorrect":
	case "gameNext":
	case "gameWin":
		delete((*rooms), roomId)
	case "gameLose":
		delete((*rooms), roomId)
	case "illegalMove":
		println("Illegal Move")
	}
	return gameState
}

func RemoveUser(request map[string]string, rooms *map[string]*RoomInfo) map[string]string {
	roomId := request["roomId"]
	userId := request["userId"]
	_, exists := (*rooms)[roomId]
	if !exists {
		return map[string]string{"error" : "no such room"}
	}
	result := (*rooms)[roomId].RemoveUser(userId)
	if result == "delete" {
		delete((*rooms), roomId)
	}
	return map[string]string{"success" : result}
}

func DeleteRoom(request map[string]string, rooms *map[string]*RoomInfo) map[string]string {
	roomId := request["roomId"]
	_, exists := (*rooms)[roomId]
	if !exists {
		return map[string]string{"error" : "no such room"}
	}
	delete((*rooms), roomId)
	return map[string]string{"success" : "deleted room"}
}

func AddSid(request map[string]string, sids *map[string]string) map[string]string {
	userId := request["userId"]
	sid := request["sid"]
	(*sids)[sid] = userId
	return map[string]string{"success" : "addedSid"}
}

func PopSid(request map[string]string, sids *map[string]string) map[string]string {
	sid := request["sid"]
	userId, exists := (*sids)[sid]
	if !exists {
		return map[string]string{"error" : "no such user"}
	}
	delete((*sids), sid)
	return map[string]string{"userId" : userId}
}