package objets

import (
	"encoding/json"
	"github.com/MrNinso/MyGoToolBox/lang/ifs"
)

//Room Column name in database
const RoomColName = "ROOMS"

// A Room entry
type Room struct {
	id        string
	title     string
	contentMd string
	imageUID  string
}

// Create a Room struct using json with given id
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

// Create a Room struct using json
func RoomFromJson(j []byte) *Room {
	m := make(map[string]interface{})
	_ = json.Unmarshal(j, &m)
	return RoomFromMap(m)
}

// Create a Room from map
func RoomFromMap(m map[string]interface{}) *Room {
	return &Room{
		id:        ifs.IfNil(m["id"], "").(string),
		title:     ifs.IfNil(m["title"], "").(string),
		contentMd: ifs.IfNil(m["contentMd"], "").(string),
		imageUID:  ifs.IfNil(m["imageUID"], "").(string),
	}
}

// Get Room UID (which is used as file name for rendered Room)
func (r Room) GetUID() string {
	return r.id
}

// Get Room Title
func (r Room) GetTitle() string {
	return r.title
}

// Get not rendered content
func (r Room) GetContetMd() string {
	return r.contentMd
}

// Get Docker Image UID (this is not the docker image name)
func (r Room) GetImageUID() string {
	return r.imageUID
}

// convert a Room to Map
func (r Room) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"id":        r.id,
		"title":     r.title,
		"contentMd": r.contentMd,
		"imageUID":  r.imageUID,
	}
}
