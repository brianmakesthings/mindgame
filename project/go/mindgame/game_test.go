package mindgame

import (
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
