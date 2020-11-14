package objets

import (
	"encoding/json"
	. "github.com/MrNinso/MyGoToolBox/lang/ifs"
)

const (
	IMAGE_COL_NAME   = "IMAGES"
	IMAGE_HEADER_KEY = "IMAGE"
)

type Image struct {
	UId             string
	Name            string
	DockerImageName string
	Dockerfile      string
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
		UId:             IfNil(m["UId"], "").(string),
		Name:            IfNil(m["Name"], "").(string),
		DockerImageName: IfNil(m["DockerImageName"], "").(string),
		Dockerfile:      IfNil(m["Dockerfile"], "").(string),
		Created:         IfNil(m["Created"], "").(string),
	}
}

func (i Image) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"UId":             i.UId,
		"Name":            i.Name,
		"DockerImageName": i.DockerImageName,
		"Dockerfile":      i.Dockerfile,
		"Created":         i.Created,
	}
}
