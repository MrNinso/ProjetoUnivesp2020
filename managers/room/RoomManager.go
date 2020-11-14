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
	"strings"
)

const (
	// Room json key in a Request Header
	RoomJsonHeaderKey = "ROOM"
)

// Room Manager type
type Manager []objets.Room

// Add Room to end of Array (this room will not render automatic)
func AddRoomToManager(r *Manager, room *objets.Room) error {
	if err := r.RenderRoom(room); err != nil {
		return err
	}

	*r = append(*r, *room)
	return nil
}

// Remove Room from Array (the rendered file will not deleted)
func RemoveRoomFromManager(r *Manager, pos int) {
	(*r)[len(*r)-1], (*r)[pos] = (*r)[pos], (*r)[len(*r)-1]

	*r = (*r)[:len(*r)-1]
}

// Update Room to on Array (this room will not render automatic)
func UpdateRoomFromManager(r *Manager, room *objets.Room) {
	(*r)[r.GetRoomPosByID(room.GetUID())] = *room
}

// Get Room Position using ID
func (r Manager) GetRoomPosByID(id string) int {
	for i := 0; i < len(r); i++ {
		if r[i].GetUID() == id {
			return i
		}
	}

	return -1
}

// Get Room using ID
func (r Manager) GetRoomByID(id string) *objets.Room {
	if pos := r.GetRoomPosByID(id); pos != -1 {
		return &r[pos]
	}

	return nil
}

// Render a Room with given id
func (r Manager) RenderRoomByID(id string) error {
	room := r.GetRoomByID(id)

	return r.RenderRoom(room)
}

// Render given Room
func (r Manager) RenderRoom(room *objets.Room) error {
	if room == nil {
		return errors.New("room é nula")
	}

	return render(room)
}

// Create CSV of Rooms
func (r Manager) ToCSV() string {
	str := ""
	for _, room := range r {
		str += fmt.Sprintf("%s,%s\n", room.GetUID(), room.GetTitle())
	}

	return str
}

func (r Manager) ToJson() (string, error) {
	var sb strings.Builder
	sb.WriteString("[")
	for i := 0; i < len(r); i++ {
		j, err := r[i].ToJson()

		if err != nil {
			return "", err
		}

		sb.Write(j)
		if i != len(r)-1 {
			sb.WriteString(",")
		}
	}

	sb.WriteString("]")

	return sb.String(), nil
}

// Render ALL Rooms in database
func RenderRooms() *Manager {
	if _, err := os.Stat("./public/res/temp"); !os.IsNotExist(err) {
		err = os.RemoveAll("./public/res/temp")
		utils.CheckPanic(&err)
	}

	err := os.MkdirAll("./public/res/temp", 0777)

	utils.CheckPanic(&err)

	rs := make(Manager, 0)

	database.Conn.ForEachRoom(func(id int, r *objets.Room) (moveOn bool) {
		e := render(r)

		if e != nil {
			log.Manager.GetLogLevel("docker").AppendLog(e.Error())
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
