package models

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Track struct {
	gorm.Model

	ID                    uint      `json:"-" gorm:"primaryKey,unique"`
	GlobalID              uuid.UUID `json:"global_id"`
	Description           string    `json:"description"`
	Status                string    `json:"status"`
	ObjectTypeDescription string    `json:"object_type_description"`
	ShapeLength           float64   `json:"shape_length"`
	Geometry              Path      `json:"geometry"`
}

func LoadFromFile(r io.Reader) []Track {
	data, _ := ioutil.ReadAll(r)

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	features := result["features"].([]interface{})
	arrSize := len(features)
	var tracks = make([]Track, arrSize)

	for i := 0; i < arrSize; i++ {
		trackData := features[i].(map[string]interface{})
		properties := trackData["properties"].(map[string]interface{})

		globalID, _ := uuid.Parse(properties["GlobalID"].(string))
		objectID := uint(properties["OBJECTID"].(float64))
		description := properties["DESCRIPTION"].(string)
		status := properties["STATUS"].(string)
		objectTypeDescription := properties["OBJECT_TYPE_DESCRIPTION"].(string)
		shapeLength := properties["SHAPE_Length"].(float64)

		path := Path{}
		geometry := trackData["geometry"].(map[string]interface{})
		coordinates := geometry["coordinates"].([]interface{})[0].([]interface{})
		for _, coord := range coordinates {
			points := coord.([]interface{})
			path.Points = append(path.Points, Point{
				X: points[0].(float64),
				Y: points[1].(float64),
			})
		}

		track := Track{
			ID:                    objectID,
			GlobalID:              globalID,
			Description:           description,
			Status:                status,
			ObjectTypeDescription: objectTypeDescription,
			ShapeLength:           shapeLength,
			Geometry:              path,
		}
		tracks[i] = track
	}
	return tracks
}
