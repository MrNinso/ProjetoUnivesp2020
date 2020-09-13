package rooms

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"strings"

	"testsGin/utils"
)

type RoomManager []Room

func (r RoomManager) GetRoomPosByID(id string) int {
	for i := 0; i < len(r); i++ {
		if r[i].getID() == id {
			return i
		}
	}

	return -1
}

func (r RoomManager) GetRoomByID(id string) *Room {
	if pos := r.GetRoomPosByID(id); pos != -1 {
		return &r[pos]
	}

	return nil
}

func (r RoomManager) RenderRoomByID(id string) {
	room := r.GetRoomByID(id)

	render(room.id, room.title)
}

func (r RoomManager) RenderRoom(room *Room) {
	render(room.id, room.title)
}

func (r RoomManager) ToCSV() string {
	str := ""
	for _, room := range r {
		str += fmt.Sprintf("%s,%s\n", room.id, room.title)
	}

	return str
}

func Render() *RoomManager {
	files, err := ioutil.ReadDir("public/rooms")

	utils.CheckPanic(&err)

	if _, errTemp := os.Stat("public/temp"); !os.IsNotExist(errTemp) {
		errTemp = os.RemoveAll("public/temp")
		utils.CheckPanic(&errTemp)
	}

	err = os.Mkdir("public/temp", 0777)

	utils.CheckPanic(&err)

	rs := make(RoomManager, len(files))

	for i := 0; i < len(files); i++ {
		rs[i] = NewRoom(
			render(files[i].Name(), uuid.New().String()),
			strings.ReplaceAll(files[i].Name(), ".md", ""),
		)
	}

	return &rs
}

func render(name, id string) string {

	pathIn := "public/rooms/" + name
	pathOut := "public/temp/" + id

	content, err := ioutil.ReadFile(pathIn)

	utils.CheckPanic(&err)

	_, _ = os.Create(pathOut)

	err = ioutil.WriteFile(pathOut, markdown.ToHTML(content, nil, nil), 666)

	utils.CheckPanic(&err)

	return id
}
