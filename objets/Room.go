package objets

import (
	"ProjetoUnivesp2020/utils"
	"encoding/json"
)

const ROOM_COL_NAME = "ROOMS"

type Room struct {
	id        string
	title     string
	contentMd string
	imageUID  string
}

func NewRoom(id, title, contentMd, imageUId string) Room {
	return Room{
		id:        id,
		title:     title,
		contentMd: contentMd,
	}
}

func NewRoomFromJson(id string, j []byte) (*Room, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(j, &m)

	if err == nil {
		r := RoomFromMap(m)
		r.id = id
		return r, nil
	}
	return nil, err
}

func RoomFromJson(j []byte) *Room {
	m := make(map[string]interface{})
	_ = json.Unmarshal(j, &m)
	return RoomFromMap(m)
}

func RoomFromMap(m map[string]interface{}) *Room {
	return &Room{
		id:        utils.IfNil(m["id"], "").(string),
		title:     utils.IfNil(m["title"], "").(string),
		contentMd: utils.IfNil(m["contentMd"], "").(string),
		imageUID:  utils.IfNil(m["imageUID"], "").(string),
	}
}

func (r Room) GetUID() string {
	return r.id
}

func (r Room) GetTitle() string {
	return r.title
}

func (r Room) GetContetMd() string {
	return r.contentMd
}

func (r Room) GetImageUID() string {
	return r.imageUID
}

func (r Room) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"id":        r.id,
		"title":     r.title,
		"contentMd": r.contentMd,
		"imageUID":  r.imageUID,
	}
}
