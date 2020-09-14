package objets

import (
	"ProjetoUnivesp2020/utils"
	"encoding/json"
	"time"
)

const IMAGE_COL_NAME = "IMAGES"

type Image struct {
	UId string
	PodmanId string
	Name string
	DockerFile string
	Created time.Time
}

func ImageFromJson(j []byte) *Image {
	m := make(map[string]interface{})
	_ = json.Unmarshal(j, &m)
	return ImageFromMap(m)
}

func ImageFromMap(m map[string]interface{}) *Image {
	return &Image{
		UId:         utils.IfNil(m["UId"], "").(string),
		PodmanId:   utils.IfNil(m["PodmanId"], "").(string),
		Name:       utils.IfNil(m["Name"], "").(string),
		DockerFile: utils.IfNil(m["DockerFile"], "").(string),
		Created:	utils.IfNil(m["Created"], "").(time.Time),
	}
}

func (i Image) ToMap() *map[string]interface{} {
	return &map[string]interface{} {
		"UId": i.UId,
		"Name": i.Name,
		"PodmanId": i.PodmanId,
		"DockerFile": i.DockerFile,
		"Created": i.Created,
	}
}