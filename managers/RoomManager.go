package managers

import (
	"ProjetoUnivesp2020/objets"
	"ProjetoUnivesp2020/utils"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"strings"
)

type RoomManager []objets.Room

func (r RoomManager) GetRoomPosByID(id string) int {
	for i := 0; i < len(r); i++ {
		if r[i].GetID() == id {
			return i
		}
	}

	return -1
}

func (r RoomManager) GetRoomByID(id string) *objets.Room {
	if pos := r.GetRoomPosByID(id); pos != -1 {
		return &r[pos]
	}

	return nil
}

func (r RoomManager) RenderRoomByID(id string) {
	room := r.GetRoomByID(id)

	render(room.GetID(), room.GetTitle())
}

func (r RoomManager) RenderRoom(room *objets.Room) {
	render(room.GetID(), room.GetTitle())
}

func (r RoomManager) ToCSV() string {
	str := ""
	for _, room := range r {
		str += fmt.Sprintf("%s,%s\n", room.GetID(), room.GetTitle())
	}

	return str
}

func RenderRooms() *RoomManager {
	files, err := ioutil.ReadDir("./public/res/rooms")

	utils.CheckPanic(&err)

	if _, errTemp := os.Stat("./public/res/temp"); !os.IsNotExist(errTemp) {
		errTemp = os.RemoveAll("./public/res/temp")
		utils.CheckPanic(&errTemp)
	}

	err = os.Mkdir("./public/res/temp", 0777)

	utils.CheckPanic(&err)

	rs := make(RoomManager, len(files))

	for i := 0; i < len(files); i++ {
		rs[i] = objets.NewRoom(
			render(files[i].Name(), uuid.New().String()),
			strings.ReplaceAll(files[i].Name(), ".md", ""),
		)
	}

	return &rs
}

func render(name, id string) string {
	pathIn := "public/res/rooms/" + name
	pathOut := "public/res/temp/" + id

	content, err := ioutil.ReadFile(pathIn)

	utils.CheckPanic(&err)

	_, _ = os.Create(pathOut)

	err = ioutil.WriteFile(pathOut, markdown.ToHTML(content, nil, nil), 666)

	utils.CheckPanic(&err)

	return id
}
