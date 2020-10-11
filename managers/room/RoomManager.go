package room

import (
	"ProjetoUnivesp2020/managers/database"
	"ProjetoUnivesp2020/managers/log"
	"ProjetoUnivesp2020/objets"
	"ProjetoUnivesp2020/utils"
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	"os"
)

const (
	ROOM_JSON_HEADER_KEY = "ROOM"
)

type RoomManager []objets.Room

func AddRoomToManager(r *RoomManager, room *objets.Room) error {
	if err := r.RenderRoom(room); err != nil {
		return err
	}

	*r = append(*r, *room)
	return nil
}

func RemoveRoomFromManager(r *RoomManager, pos int) {
	(*r)[len(*r)-1], (*r)[pos] = (*r)[pos], (*r)[len(*r)-1]

	*r = (*r)[:len(*r)-1]
}

func UpdateRoomFromManager(r *RoomManager, room *objets.Room) {
	(*r)[r.GetRoomPosByID(room.GetUID())] = *room
}

func (r RoomManager) GetRoomPosByID(id string) int {
	for i := 0; i < len(r); i++ {
		if r[i].GetUID() == id {
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

func (r RoomManager) RenderRoomByID(id string) error {
	room := r.GetRoomByID(id)

	return r.RenderRoom(room)
}

func (r RoomManager) RenderRoom(room *objets.Room) error {
	if room == nil {
		return errors.New("room é nula")
	}

	return render(room)
}

func (r RoomManager) ToCSV() string {
	str := ""
	for _, room := range r {
		str += fmt.Sprintf("%s,%s\n", room.GetUID(), room.GetTitle())
	}

	return str
}

func RenderRooms() *RoomManager {
	if _, err := os.Stat("./public/res/temp"); !os.IsNotExist(err) {
		err = os.RemoveAll("./public/res/temp")
		utils.CheckPanic(&err)
	}

	err := os.MkdirAll("./public/res/temp", 0777)

	utils.CheckPanic(&err)

	rs := make(RoomManager, 0)

	database.Conn.ForEachRoom(func(id int, r *objets.Room) (moveOn bool) {
		e := render(r)

		if e != nil {
			log.LogManager.GetLogLevel("docker").AppendLog(e.Error())
		}

		rs = append(rs, *r)

		return true
	})

	return &rs
}

func render(r *objets.Room) error {
	pathOut := "public/res/temp/" + r.GetUID()

	_, err := os.Create(pathOut)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pathOut, markdown.ToHTML([]byte(r.GetContetMd()), nil, nil), 666)

	return err
}
