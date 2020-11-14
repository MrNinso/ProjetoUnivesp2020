package objets

import (
	"encoding/json"
	"github.com/MrNinso/MyGoToolBox/lang/ifs"
)

//Room Column name in database
const RoomColName = "ROOMS"

// A Room entry
type Room struct {
	UId       string
	title     string
	contentMd string
	imageUId  string
}

// Create a Room struct using json with given UId
func NewRoomFromJson(id string, j []byte) (*Room, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(j, &m)

	if err == nil {
		r := RoomFromMap(m)
		r.UId = id
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
		UId:       ifs.IfNil(m["UId"], "").(string),
		title:     ifs.IfNil(m["title"], "").(string),
		contentMd: ifs.IfNil(m["contentMd"], "").(string),
		imageUId:  ifs.IfNil(m["imageUId"], "").(string),
	}
}

// Get Room UID (which is used as file name for rendered Room)
func (r Room) GetUID() string {
	return r.UId
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
	return r.imageUId
}

// convert a Room to Map
func (r Room) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"UId":       r.UId,
		"title":     r.title,
		"contentMd": r.contentMd,
		"imageUId":  r.imageUId,
	}
}

func (r Room) ToJson() ([]byte, error) {
	return json.Marshal(r.ToMap())
}
