package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/CoryEvans2324/SystemsDesignAppAPI/database"
	"github.com/CoryEvans2324/SystemsDesignAppAPI/models"
	"gorm.io/gorm/clause"
)

func Index(w http.ResponseWriter, r *http.Request) {

}

func UploadTracks(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracks := models.LoadFromFile(file)
	log.Println(len(tracks))
	tx := database.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&tracks)
	if tx.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func GetTracks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	if limit > 1000 {
		limit = 1000
	}

	var tracks []struct {
		models.Track
		Distance float64 `json:"distance"`
	}

	// result := database.DB.Limit(limit).Find(&tracks)
	result := database.DB.Raw(`
SELECT
	calculate_distance(?, ?, latlon.ll[1], latlon.ll[0], 'K') as distance,
	id, description, status, shape_length, object_type_description, geometry
FROM (
	SELECT geometry::POLYGON::POINT as ll, id, description, status, shape_length, object_type_description, geometry
	FROM tracks
) as latlon
ORDER BY distance
LIMIT ?;
	`, lat, lon, limit).Scan(&tracks)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(tracks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
