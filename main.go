package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	connStr := os.Getenv("DATABASE_URL")
	dbPool, _ = pgxpool.New(context.Background(), connStr)

	http.HandleFunc("/search", findDrivers)
	http.ListenAndServe(":8080", nil)
}

func findDrivers(w http.ResponseWriter, r *http.Request) {
	// MG Road, Bengaluru Coordinates
	userLng, userLat := 77.6113, 12.9738

	// The query uses the <-> operator for "Nearest Neighbor" sorting
	query := `
		SELECT id, ST_X(location::geometry) as lng, ST_Y(location::geometry) as lat
		FROM drivers
		WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 5000)
		ORDER BY location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		LIMIT 10;`

	rows, err := dbPool.Query(context.Background(), query, userLng, userLat)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	// Use make() to avoid returning 'null' if list is empty
	list := make([]map[string]interface{}, 0)

	for rows.Next() {
		var id string
		var lng, lat float64
		rows.Scan(&id, &lng, &lat)
		list = append(list, map[string]interface{}{"id": id, "lat": lat, "lng": lng})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
