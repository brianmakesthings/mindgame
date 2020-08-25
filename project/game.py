import os, struct
import random

userIds = {}
rooms = {}
maxLevel = 3

class room:
    def __init__(self, ownerId, roomId):
        self.id = roomId
        self.owner = ownerId
        self.players = [playerData(ownerId)]
        self.lives = 0
        self.level = 0
    def joinRoom(self, userId):
        if len(self.players) >=4:
            return "Room is Full"
        for player in self.players:
            if player.id == userId:
                return "Player already in room"
        self.players.append(playerData(userId))
        return "Success"
    def startRound(self, level):
        self.level = level
        dealtCards = random.sample(range(1,101), k=len(self.players)*level)
        for i in range(len(self.players)):
            self.players[i].giveCards(dealtCards[i*level:(i+1)*level])
    def startGame(self):
        self.lives = len(self.players)
        self.startRound(1)
    def hasPassedLevel(self):
        for player in self.players:
            if len(player.cards) > 0:
                return False
        return True
    def playCard(self, userId, playedCard):
        correctCard = 999
        for player in self.players:
            for card in player.cards:
                if card < correctCard:
                    correctCard = card
                if card == playedCard:
                    if userId != player.id:
                        return "Illegal Move"
                    else:
                        player.cards.remove(playedCard)
        
        if playedCard == correctCard:
            if self.hasPassedLevel():
                if self.level+1>maxLevel:
                    return "gameWin"
                self.startRound(self.level+1)
                return "gameNext"
            else:
                return "correct"
        else:
            for player in self.players:
                for card in player.cards:
                    if card <= playedCard:
                        player.cards.remove(card)
            self.lives = self.lives - 1
            if self.lives <= 0:
                return "gameLose"
            else:
                if self.hasPassedLevel():
                    if self.level+1>maxLevel:
                        return "gameWin"
                    self.startRound(self.level+1)
                    return "gameNext"
                return "incorrect"
        


class playerData:
    def __init__(self, id):
        self.id = id
    
    def giveCards(self, cards):
        self.cards = cards

    def takeCard(self, card):
        self.cards.remove(card)

def generateUserId(username):
    idGenerated = False
    idCandidate = ""
    while (not idGenerated):
        randInt = struct.unpack('I', os.urandom(4))[0] % 1000000
        idCandidate = username + str(randInt)
        if (idCandidate not in userIds):
            userIds[idCandidate] = username
            idGenerated = True
    return idCandidate

def createRoom(ownerId):
    idGenerated = False
    randInt = 0
    while (not idGenerated):
        randInt = struct.unpack('I', os.urandom(4))[0] % 1000000
        if (randInt not in rooms):
            rooms[randInt] = room(ownerId, randInt)
            idGenerated = True
    return randInt

def getParticipants(roomId):
    print(rooms)
    result = ([],rooms[int(roomId)].owner)
    for player in rooms[int(roomId)].players:
        result[0].append(player.id)
    return result


def playCard(roomId, userId, card):
    roomId = int(roomId)
    moveResult = rooms[roomId].playCard(userId, int(card))
    gameState = getGameState(roomId)
    gameState["moveResult"] = moveResult
    if moveResult == "gameWin" or moveResult =="gameLose":
        rooms.pop(roomId)
    return gameState

def removeUser(roomId, userId):
    for player in rooms[roomId].players:
        if player.id == userId:
            rooms[roomId].players.remove(player)
            return "success"
    return "no such player found"

def joinRoom(userId, roomId):
    print(rooms)
    return rooms[int(roomId)].joinRoom(userId)

def startGame(roomId):
    return rooms[int(roomId)].startGame()

def getGameState(roomId):
    room = rooms[roomId]
    for player in room.players:
        player.cards.sort()
    return {"lives" : room.lives, "level" : room.level, "playerState" : room.players}