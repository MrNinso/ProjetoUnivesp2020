package objets

import (
	"ProjetoUnivesp2020/utils"
	"encoding/json"
)

const (
	IMAGE_COL_NAME   = "IMAGES"
	IMAGE_HEADER_KEY = "IMAGE"
)

type Image struct {
	UId             string
	Name            string
	DockerImageName string
	DockerFile      string
	Created         string
}

func ImageFromJson(j []byte) (*Image, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(j, &m)

	if err != nil {
		return nil, err
	}

	return ImageFromMap(m), nil
}

func ImageFromMap(m map[string]interface{}) *Image {
	return &Image{
		UId:             utils.IfNil(m["UId"], "").(string),
		Name:            utils.IfNil(m["Name"], "").(string),
		DockerImageName: utils.IfNil(m["DockerImageName"], "").(string),
		DockerFile:      utils.IfNil(m["DockerFile"], "").(string),
		Created:         utils.IfNil(m["Created"], "").(string),
	}
}

func (i Image) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"UId":             i.UId,
		"Name":            i.Name,
		"DockerImageName": i.DockerImageName,
		"DockerFile":      i.DockerFile,
		"Created":         i.Created,
	}
}
