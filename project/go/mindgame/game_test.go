package mindgame

import (
	"fmt"
	"strings"
	"testing"
)

func beforeEach(t *testing.T) {

}

func TestGenerateUserId(t *testing.T) {
	request := make(map[string]string)
	userIds := make(map[string]bool)
	request["string"] = "abc"
	abcUserName1 := GenerateUserId(request, &userIds)
	if !strings.HasPrefix(abcUserName1["id"], "abc") {
		t.Errorf("Generate User Id returned %s; it should start with %s", abcUserName1["id"], "abc")
	}
	if !userIds[abcUserName1["id"]] {
		t.Errorf("Generate User Id failed to update list of user ids")
	}

	abcUserName2 := GenerateUserId(request, &userIds)
	if !strings.HasPrefix(abcUserName2["id"], "abc") {
		t.Errorf("Generate User Id returned %s; it should start with %s", abcUserName1["id"], "abc")
	}
	if abcUserName1["id"] == abcUserName2["id"] {
		t.Errorf("A collision of userIds occured. String 1: %s, String 2: %s", abcUserName1["id"], abcUserName2["id"])
	}
}

func TestCreateRoom(t *testing.T) {
	request := make(map[string]string)
	rooms := make(map[string]*(RoomInfo))
	ownerId := "abc23874687"
	request["ownerId"] = ownerId
	result := CreateRoom(request, &rooms)
	room := rooms[result["roomId"]]
	if room == nil {
		t.Errorf("Failed to create a room")
	}
	if room.owner != ownerId {
		t.Errorf("Owner is %s; expected %s", room.owner, ownerId)
	}
}

func TestJoinRoom(t *testing.T) {
	createRoomRequest := make(map[string]string)
	rooms := make(map[string]*RoomInfo)
	ownerId := "abc123"
	createRoomRequest["ownerId"] = ownerId
	result := CreateRoom(createRoomRequest, &rooms)

	joinRoomRequest := make(map[string]string)
	joinRoomRequest["roomId"] = result["roomId"]
	joinRoomRequest["userId"] = ownerId
	result = JoinRoom(joinRoomRequest, &rooms)
	_, hasError := result["error"]
	if !hasError {
		t.Errorf("Failed to reject duplicate join")
	}

	okJoin := [3]string{"two123", "three123", "four123"}
	for i, participant := range okJoin {
		joinRoomRequest["userId"] = participant
		joinResult := JoinRoom(joinRoomRequest, &rooms)
		_, ok := joinResult["success"]
		if !ok {
			t.Errorf("Failed to add participant %d to room", i+2)
		}
	}

	fullUser := "five123"
	joinRoomRequest["userId"] = fullUser
	fullResult := JoinRoom(joinRoomRequest, &rooms)
	_, hasError = fullResult["error"]
	if !hasError {
		t.Errorf("Failed to reject adding participant to full room")
	}
}

func TestGetParticipants(t *testing.T) {
	createRoomRequest := make(map[string]string)
	rooms := make(map[string]*RoomInfo)
	ownerId := "abc123"
	createRoomRequest["ownerId"] = ownerId
	result := CreateRoom(createRoomRequest, &rooms)

	joinRoomRequest := make(map[string]string)
	joinRoomRequest["roomId"] = result["roomId"]
	okJoin := [3]string{"two123", "three123", "four123"}
	for _, participant := range okJoin {
		joinRoomRequest["userId"] = participant
		JoinRoom(joinRoomRequest, &rooms)
	}

	getParticipantsRequest := make(map[string]string)
	getParticipantsRequest["roomId"] = result["roomId"]
	participantsResult := GetParticipants(getParticipantsRequest, &rooms)
	listOwner := [1]string{"abc123"}
	actualParticipantList := append(listOwner[:], okJoin[:]...)
	stringifiedParticipants := fmt.Sprintf("%v", actualParticipantList)
	if participantsResult["participants"] != stringifiedParticipants {
		t.Errorf("Get Participants failed, expected: %v, got %v", stringifiedParticipants, participantsResult["participants"])
	}
	if participantsResult["owner"] != ownerId {
		t.Errorf("Get Participants: Owner failed, expected %s, got %s", ownerId, participantsResult["owner"])
	}
}
